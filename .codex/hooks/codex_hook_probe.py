#!/usr/bin/env python3
"""Sanitize one Codex hook payload and append a diagnostic JSONL record."""

import fcntl
import hashlib
import json
import os
import stat
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path


MAX_INPUT_BYTES = 1024 * 1024
MAX_DEPTH = 8
MAX_ARRAY_ITEMS = 64
MAX_OBJECT_FIELDS = 256

SAFE_COMMANDS = {
    "printf 'codex-hook-probe-success\\n'": 0,
    "sh -c 'printf \"codex-hook-probe-failure\\n\" >&2; exit 7'": 7,
}

ALLOWED_HOOK_EVENTS = {
    "SessionStart",
    "SubagentStart",
    "SubagentStop",
    "PreToolUse",
    "PostToolUse",
}

PRESENCE_FIELDS = (
    "session_id",
    "turn_id",
    "cwd",
    "permission_mode",
    "agent_id",
    "agent_type",
    "tool_name",
    "tool_use_id",
)

CORRELATION_FIELDS = {
    "session_id",
    "turn_id",
    "agent_id",
    "agent_type",
    "tool_name",
    "tool_use_id",
}

FORBIDDEN_INPUT_FIELDS = {
    "workflow_id",
    "run_id",
    "contract_hash",
    "source_kind",
    "trust_basis",
    "provenance",
    "observations",
    "reconciliation",
    "verdict",
    "completion_verdict",
}

EXIT_STATUS_KEYS = {"exitcode", "exitstatus"}

REPOSITORY_ROOT = Path(__file__).resolve().parents[2]
OUTPUT_PATH = REPOSITORY_ROOT / "local_artifacts" / "codex-hook-probe" / "events.jsonl"


def json_type(value):
    if value is None:
        return "null"
    if isinstance(value, bool):
        return "boolean"
    if isinstance(value, (int, float)):
        return "number"
    if isinstance(value, str):
        return "string"
    if isinstance(value, list):
        return "array"
    if isinstance(value, dict):
        return "object"
    raise TypeError("value is not a JSON type")


def digest_text(value):
    encoded = value.encode("utf-8")
    return {
        "length_bytes": len(encoded),
        "sha256": hashlib.sha256(encoded).hexdigest(),
    }


def text_length(value):
    return {"length_bytes": len(value.encode("utf-8"))}


def reject_forbidden_fields(value):
    if isinstance(value, dict):
        for key, nested in value.items():
            if key in FORBIDDEN_INPUT_FIELDS:
                raise ValueError("production workflow field is not accepted")
            reject_forbidden_fields(nested)
    elif isinstance(value, list):
        for nested in value:
            reject_forbidden_fields(nested)


def summarize_shape(value, path=(), depth=0):
    kind = json_type(value)
    summary = {"json_type": kind}
    if depth >= MAX_DEPTH:
        summary["shape_truncated"] = True
        return summary

    if kind == "object":
        keys = sorted(value)
        summary["field_count"] = len(keys)
        summary["fields"] = [
            {
                "name": key,
                "shape": summarize_shape(value[key], path + (key,), depth + 1),
            }
            for key in keys[:MAX_OBJECT_FIELDS]
        ]
        if len(keys) > MAX_OBJECT_FIELDS:
            summary["truncated_field_count"] = len(keys) - MAX_OBJECT_FIELDS
        return summary

    if kind == "array":
        summary["length"] = len(value)
        summary["items"] = [
            summarize_shape(nested, path + ("[]",), depth + 1)
            for nested in value[:MAX_ARRAY_ITEMS]
        ]
        if len(value) > MAX_ARRAY_ITEMS:
            summary["truncated_item_count"] = len(value) - MAX_ARRAY_ITEMS
        return summary

    if kind == "string":
        summary.update(text_length(value))
        if path == ("tool_input", "command"):
            summary.update(digest_text(value))
            if value in SAFE_COMMANDS:
                summary["safe_command"] = value
        return summary

    return summary


def field_presence(payload, name):
    if name not in payload:
        return {"present": False}
    value = payload[name]
    descriptor = {"present": True, "json_type": json_type(value)}
    if isinstance(value, str):
        descriptor["length_bytes"] = len(value.encode("utf-8"))
        if name in CORRELATION_FIELDS:
            descriptor["sha256"] = hashlib.sha256(value.encode("utf-8")).hexdigest()
    return descriptor


def safe_command_expected_status(payload):
    tool_input = payload.get("tool_input")
    if not isinstance(tool_input, dict):
        return None
    command = tool_input.get("command")
    if not isinstance(command, str):
        return None
    return SAFE_COMMANDS.get(command)


def find_exit_status_fields(value, path=("tool_response",), depth=0):
    if depth >= MAX_DEPTH:
        return []
    found = []
    if isinstance(value, dict):
        for key in sorted(value):
            nested = value[key]
            nested_path = path + (key,)
            normalized = key.lower().replace("_", "").replace("-", "")
            if normalized in EXIT_STATUS_KEYS:
                descriptor = {
                    "path": ".".join(nested_path),
                    "json_type": json_type(nested),
                }
                if isinstance(nested, int) and not isinstance(nested, bool):
                    descriptor["value"] = nested
                elif isinstance(nested, str):
                    descriptor.update(text_length(nested))
                found.append(descriptor)
            found.extend(find_exit_status_fields(nested, nested_path, depth + 1))
    elif isinstance(value, list):
        for index, nested in enumerate(value[:MAX_ARRAY_ITEMS]):
            found.extend(
                find_exit_status_fields(nested, path + ("[{}]".format(index),), depth + 1)
            )
    return found


def find_top_level_exit_status_fields(payload):
    found = []
    for key in sorted(payload):
        normalized = key.lower().replace("_", "").replace("-", "")
        if normalized not in EXIT_STATUS_KEYS:
            continue
        value = payload[key]
        descriptor = {"path": key, "json_type": json_type(value)}
        if isinstance(value, int) and not isinstance(value, bool):
            descriptor["value"] = value
        elif isinstance(value, str):
            descriptor.update(text_length(value))
        found.append(descriptor)
    return found


def sanitize_payload(payload, now=None, record_id=None):
    if not isinstance(payload, dict):
        raise ValueError("hook input must be one JSON object")
    reject_forbidden_fields(payload)

    hook_event_name = payload.get("hook_event_name")
    if hook_event_name not in ALLOWED_HOOK_EVENTS:
        raise ValueError("hook event is outside the diagnostic allowlist")
    if hook_event_name in {"PreToolUse", "PostToolUse"} and payload.get("tool_name") != "Bash":
        raise ValueError("tool event is outside the Bash-only diagnostic boundary")

    if now is None:
        now = datetime.now(timezone.utc)
    if now.tzinfo is None:
        raise ValueError("collector timestamp must be timezone-aware")
    if record_id is None:
        record_id = "capture-" + uuid.uuid4().hex

    record = {
        "collector_schema_version": "1.0",
        "collector_record_id": record_id,
        "collector_timestamp": now.astimezone(timezone.utc).isoformat().replace("+00:00", "Z"),
        "hook_event_name": hook_event_name,
        "top_level_fields": [
            {"name": key, "json_type": json_type(payload[key])}
            for key in sorted(payload)
        ],
        "field_presence": {
            name: field_presence(payload, name) for name in PRESENCE_FIELDS
        },
        "tool_input_shape": {"present": "tool_input" in payload},
        "tool_response_shape": {"present": "tool_response" in payload},
    }

    if "tool_input" in payload:
        record["tool_input_shape"]["shape"] = summarize_shape(
            payload["tool_input"], ("tool_input",)
        )
    if "tool_response" in payload:
        record["tool_response_shape"]["shape"] = summarize_shape(
            payload["tool_response"], ("tool_response",)
        )

    expected_status = safe_command_expected_status(payload)
    exit_fields = find_top_level_exit_status_fields(payload)
    exit_fields.extend(find_exit_status_fields(payload.get("tool_response")))
    explicit_status = {"exists": bool(exit_fields), "fields": exit_fields}
    if expected_status is not None:
        explicit_status["expected_for_safe_command"] = expected_status
        for field in exit_fields:
            if "value" in field:
                field["matches_expected"] = field["value"] == expected_status
    record["explicit_exit_status"] = explicit_status
    return record


def read_hook_object(stream):
    raw = stream.read(MAX_INPUT_BYTES + 1)
    if len(raw) > MAX_INPUT_BYTES:
        raise ValueError("hook input exceeds diagnostic size limit")
    try:
        value = json.loads(raw.decode("utf-8"))
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        raise ValueError("hook input is not one UTF-8 JSON value") from exc
    if not isinstance(value, dict):
        raise ValueError("hook input must be one JSON object")
    return value


def ensure_private_directory(path):
    path.mkdir(mode=0o700, parents=True, exist_ok=True)
    info = path.lstat()
    if stat.S_ISLNK(info.st_mode) or not stat.S_ISDIR(info.st_mode):
        raise OSError("probe output directory is not a real directory")
    path.chmod(0o700)


def append_record(path, record):
    ensure_private_directory(path.parent)
    flags = os.O_WRONLY | os.O_APPEND | os.O_CREAT
    if hasattr(os, "O_NOFOLLOW"):
        flags |= os.O_NOFOLLOW
    descriptor = os.open(str(path), flags, 0o600)
    try:
        info = os.fstat(descriptor)
        if not stat.S_ISREG(info.st_mode):
            raise OSError("probe output is not a regular file")
        os.fchmod(descriptor, 0o600)
        line = json.dumps(record, sort_keys=True, separators=(",", ":")).encode("utf-8") + b"\n"
        fcntl.flock(descriptor, fcntl.LOCK_EX)
        try:
            view = memoryview(line)
            while view:
                written = os.write(descriptor, view)
                if written <= 0:
                    raise OSError("short write to probe output")
                view = view[written:]
            os.fsync(descriptor)
        finally:
            fcntl.flock(descriptor, fcntl.LOCK_UN)
    finally:
        os.close(descriptor)


def run(stream, output_path=OUTPUT_PATH, now=None, record_id=None):
    """Run observationally: collection failures never steer or block Codex."""
    try:
        payload = read_hook_object(stream)
        record = sanitize_payload(payload, now=now, record_id=record_id)
        append_record(Path(output_path), record)
    except Exception:
        return 0
    return 0


def main():
    os.umask(0o077)
    return run(sys.stdin.buffer)


if __name__ == "__main__":
    sys.exit(main())
