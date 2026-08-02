#!/usr/bin/env python3
"""Synthetic sanitizer tests. These fixtures are never live hook evidence."""

import hashlib
import io
import json
import multiprocessing
import stat
import tempfile
import unittest
from datetime import datetime, timezone
from pathlib import Path

import codex_hook_probe as probe


SUCCESS_COMMAND = "printf 'codex-hook-probe-success\\n'"
FAILURE_COMMAND = "sh -c 'printf \"codex-hook-probe-failure\\n\" >&2; exit 7'"


def field_shape(shape, name):
    for field in shape["fields"]:
        if field["name"] == name:
            return field["shape"]
    raise AssertionError("missing shape field {!r}".format(name))


def append_worker(path, index):
    probe.append_record(
        Path(path),
        {
            "collector_schema_version": "1.0",
            "collector_record_id": "capture-{:02d}".format(index),
        },
    )


class SanitizerTests(unittest.TestCase):
    def make_payload(self, command=SUCCESS_COMMAND, exit_code=0):
        return {
            "session_id": "session-private-123",
            "turn_id": "turn-private-456",
            "transcript_path": "/private/conversation/rollout.jsonl",
            "cwd": "/private/repository/path",
            "permission_mode": "default",
            "hook_event_name": "PostToolUse",
            "tool_name": "Bash",
            "tool_use_id": "tool-private-789",
            "tool_input": {"command": command},
            "tool_response": {
                "output": "complete command output must not survive",
                "exit_code": exit_code,
                "metadata": {"credential": "nested credential must not survive"},
                "wall_time_seconds": 0.25,
            },
            "prompt": "user prompt must not survive",
            "last_assistant_message": "assistant response must not survive",
            "environment": {"TOKEN": "credential must not survive"},
        }

    def sanitize(self, payload):
        return probe.sanitize_payload(
            payload,
            now=datetime(2026, 8, 2, 20, 0, tzinfo=timezone.utc),
            record_id="capture-test",
        )

    def test_safe_command_and_response_are_sanitized(self):
        record = self.sanitize(self.make_payload())
        serialized = json.dumps(record, sort_keys=True)

        self.assertEqual(record["hook_event_name"], "PostToolUse")
        self.assertEqual(record["collector_timestamp"], "2026-08-02T20:00:00Z")
        for secret in (
            "session-private-123",
            "turn-private-456",
            "/private/conversation/rollout.jsonl",
            "/private/repository/path",
            "tool-private-789",
            "complete command output must not survive",
            "nested credential must not survive",
            "user prompt must not survive",
            "assistant response must not survive",
            "credential must not survive",
        ):
            self.assertNotIn(secret, serialized)

        input_shape = record["tool_input_shape"]["shape"]
        command_shape = field_shape(input_shape, "command")
        self.assertEqual(command_shape["safe_command"], SUCCESS_COMMAND)
        self.assertEqual(command_shape["json_type"], "string")

        response_shape = record["tool_response_shape"]["shape"]
        output_shape = field_shape(response_shape, "output")
        self.assertEqual(output_shape["json_type"], "string")
        self.assertEqual(
            output_shape["length_bytes"],
            len(b"complete command output must not survive"),
        )
        self.assertNotIn("sha256", output_shape)
        metadata_shape = field_shape(response_shape, "metadata")
        credential_shape = field_shape(metadata_shape, "credential")
        self.assertEqual(
            credential_shape["length_bytes"],
            len(b"nested credential must not survive"),
        )
        self.assertNotIn("sha256", credential_shape)

        status = record["explicit_exit_status"]
        self.assertTrue(status["exists"])
        self.assertEqual(status["expected_for_safe_command"], 0)
        self.assertEqual(status["fields"][0]["path"], "tool_response.exit_code")
        self.assertEqual(status["fields"][0]["value"], 0)
        self.assertTrue(status["fields"][0]["matches_expected"])

    def test_non_allowlisted_command_is_digest_only(self):
        command = "echo sensitive-unapproved-command"
        record = self.sanitize(self.make_payload(command=command))
        serialized = json.dumps(record, sort_keys=True)
        self.assertNotIn(command, serialized)

        command_shape = field_shape(record["tool_input_shape"]["shape"], "command")
        self.assertNotIn("safe_command", command_shape)
        self.assertEqual(command_shape["length_bytes"], len(command.encode("utf-8")))
        self.assertEqual(
            command_shape["sha256"], hashlib.sha256(command.encode("utf-8")).hexdigest()
        )

    def test_failure_command_preserves_only_safe_command_and_exit_status(self):
        record = self.sanitize(self.make_payload(command=FAILURE_COMMAND, exit_code=7))
        status = record["explicit_exit_status"]
        self.assertEqual(status["expected_for_safe_command"], 7)
        self.assertEqual(status["fields"][0]["value"], 7)
        self.assertTrue(status["fields"][0]["matches_expected"])

    def test_top_level_explicit_exit_status_is_detected(self):
        payload = self.make_payload()
        del payload["tool_response"]["exit_code"]
        payload["exit_status"] = 0
        record = self.sanitize(payload)
        status = record["explicit_exit_status"]
        self.assertTrue(status["exists"])
        self.assertEqual(status["fields"], [{
            "path": "exit_status",
            "json_type": "number",
            "value": 0,
            "matches_expected": True,
        }])

    def test_presence_is_independent_and_identifiers_are_digested(self):
        record = self.sanitize(
            {
                "hook_event_name": "SubagentStart",
                "session_id": "session-1",
                "turn_id": "turn-1",
                "cwd": "/private/repository/path",
                "agent_id": "agent-1",
                "agent_type": "specialist",
                "permission_mode": "default",
            }
        )
        presence = record["field_presence"]
        self.assertTrue(presence["session_id"]["present"])
        self.assertIn("sha256", presence["session_id"])
        self.assertTrue(presence["agent_type"]["present"])
        self.assertFalse(presence["tool_name"]["present"])
        self.assertTrue(presence["cwd"]["present"])
        self.assertNotIn("sha256", presence["cwd"])

    def test_forbidden_production_fields_are_not_accepted(self):
        for forbidden in sorted(probe.FORBIDDEN_INPUT_FIELDS):
            payload = {"hook_event_name": "SessionStart", forbidden: "value"}
            with self.subTest(field=forbidden):
                with self.assertRaises(ValueError):
                    self.sanitize(payload)

    def test_non_object_malformed_and_unapproved_events_fail_closed(self):
        with self.assertRaises(ValueError):
            probe.read_hook_object(io.BytesIO(b"[]"))
        with self.assertRaises(ValueError):
            probe.read_hook_object(io.BytesIO(b"{} {}"))
        with self.assertRaises(ValueError):
            self.sanitize({"hook_event_name": "UserPromptSubmit"})

    def test_non_bash_tool_payload_is_rejected(self):
        payload = self.make_payload()
        payload["tool_name"] = "mcp__private__tool"
        payload["tool_input"] = {"arbitrary_payload": "must not be inspected"}
        with self.assertRaises(ValueError):
            self.sanitize(payload)


class StorageTests(unittest.TestCase):
    def test_run_returns_success_without_writing_invalid_input(self):
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "probe" / "events.jsonl"
            self.assertEqual(probe.run(io.BytesIO(b"not-json"), output), 0)
            self.assertFalse(output.exists())

    def test_private_modes_and_intact_append(self):
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "probe" / "events.jsonl"
            probe.append_record(output, {"collector_record_id": "capture-1"})
            probe.append_record(output, {"collector_record_id": "capture-2"})

            self.assertEqual(stat.S_IMODE(output.parent.stat().st_mode), 0o700)
            self.assertEqual(stat.S_IMODE(output.stat().st_mode), 0o600)
            records = [json.loads(line) for line in output.read_text().splitlines()]
            self.assertEqual(
                [record["collector_record_id"] for record in records],
                ["capture-1", "capture-2"],
            )

    def test_concurrent_process_appends_are_intact_without_order_claim(self):
        with tempfile.TemporaryDirectory() as temporary:
            output = Path(temporary) / "probe" / "events.jsonl"
            context = multiprocessing.get_context("fork")
            processes = [
                context.Process(target=append_worker, args=(str(output), index))
                for index in range(12)
            ]
            for process in processes:
                process.start()
            for process in processes:
                process.join(10)
                self.assertEqual(process.exitcode, 0)

            records = [json.loads(line) for line in output.read_text().splitlines()]
            self.assertEqual(len(records), len(processes))
            self.assertEqual(
                {record["collector_record_id"] for record in records},
                {"capture-{:02d}".format(index) for index in range(12)},
            )

    def test_symlink_output_is_rejected(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            output_dir = root / "probe"
            output_dir.mkdir()
            target = root / "target"
            target.write_text("unchanged")
            output = output_dir / "events.jsonl"
            output.symlink_to(target)
            with self.assertRaises(OSError):
                probe.append_record(output, {"collector_record_id": "capture-1"})
            self.assertEqual(target.read_text(), "unchanged")


if __name__ == "__main__":
    unittest.main()
