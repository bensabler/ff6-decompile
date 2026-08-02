Read `CLAUDE.md`, `.claude/README.md`, `docs/FRESH_SESSION_SEQUENCE.md`, `docs/checkpoints/LATEST.md`, `dashboards/CURRENT_FOCUS.md`, and all files under `docs/static-analysis/`.

Begin with:

```text
/resume-session
/bootstrap-ghidra
```

The local workstation now has Ghidra 12.1.2 with a matching SNES Loader build. The ROM is imported through the SNES ROM Loader as HiROM using `65816:LE:24:snes`. Known ROMCPU addresses including `$C09B5C` and `$C0B593` resolve directly.

A candidate function was manually created at `ROMCPU:$C0B593`. Do not accept its boundary, disassembly context, pseudocode, or semantic meaning as confirmed. Audit the repository's prior claims about this address, verify bytes, recover the 65816 entry state from Mesen, and run a bounded static/runtime correlation before renaming it.

Use:

```text
/correlate-static-runtime ROMCPU:$C0B593
```

Preserve the existing experiment queue unless the latest checkpoint identifies interrupted work. Static analysis may run in parallel as bounded reconnaissance, but runtime verification and one active experiment remain the authority.
