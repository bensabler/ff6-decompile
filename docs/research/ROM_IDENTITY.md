# ROM Identity

Record filename locally, byte size, copier-header status, region/revision, SHA-256, and mapping. Do not commit the ROM or full byte ranges.

## Recorded identity (2026-07-29, V4 migration)

| Field | Value | Confidence |
|---|---|---|
| Local filename | `Final Fantasy III (USA).sfc` (stored outside the repository, on the local machine) | Confirmed |
| Byte size | 3,145,728 bytes (`0x300000`; 24 Mbit) | Confirmed (measured) |
| Copier header | None (size ≡ 0 mod `0x8000`) | Confirmed (measured) |
| SHA-256 | `0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2` | Confirmed (computed 2026-07-29) |
| Region/label | Final Fantasy III (USA), SNES | Confirmed (filename/title screen) |
| Revision | Not yet determined (internal ROM header not inspected) | Unknown |
| Mapping | HiROM | Strong hypothesis — bank `$C1`/`$C2` code addresses and the 3 MiB size are consistent with HiROM; not yet verified against `ROMFILE:` offsets |

Session 002 (and the undocumented Session 003 evidence in `mesen/out/`) used
this file per the session record and launch command. Session 001 predates
identity recording; it is presumed to have used the same local copy, but that
is an inference, not an observation. If the local ROM is ever replaced,
recompute the hash and record a new identity block; do not overwrite this one.
