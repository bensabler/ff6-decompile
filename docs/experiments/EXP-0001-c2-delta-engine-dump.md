# EXP-0001: Verify Session 003 disassembly claims by ROM dump

- **Status:** running (2026-07-29, overnight session)
- **Question:** Do the ROM bytes at `ROMCPU:$C212F0`–`$C2141F` contain the
  instructions `internal/game/battle/battle.go` claims — dispatch
  `JSR ($131F,X)` near `$C212FF`, HP-delta routine entry `$C21323` with
  heal-clamp store `$C21338` and damage store `$C21347`, MP routine entry
  `$C21350`, death handler `$C21390` with zero store `$C21396`, delta
  fetch `$C213A7` — and does the claimed pointer table at `ROMCPU:$C2131F`
  contain `$1323`/`$1350`-family entries?
- **Starting state:** Fresh Mesen boot to the title screen with
  `mesen/bridge.lua` loaded. ROM identity per
  [ROM_IDENTITY.md](../research/ROM_IDENTITY.md) (SHA-256 `0f51b4fc…d5e2`).
  ROM contents are static; no gameplay state is involved, so the dump is
  deterministic regardless of when it is taken.
- **Controlled variables:** none needed beyond ROM identity; CPU-bus reads
  of ROM banks do not depend on emulation state.
- **Observation method:** bridge command `read cpu C212F0 304` (decimal
  length) → hex bytes to `mesen/out/resp.txt`; archived to
  `mesen/out/rom_C212F0_304.hex`. Hand-disassembly of the bytes, verified
  arithmetically (branch targets, operand addresses).
- **Expected outcomes:**
  - *Supports claims:* bytes decode to the claimed instructions at the
    claimed addresses; store opcodes are 3-byte `STA abs,X`-family at
    `$C21338`/`$C21347`/`$C21396` (matching the +3 post-instruction PC
    interpretation of the Session 003 write captures).
  - *Refutes claims:* bytes decode to different instructions, or the
    claimed addresses fall mid-instruction.
- **Falsifying outcome:** any claimed instruction that the dump contradicts
  downgrades that specific battle.go claim to Refuted and opens a
  contradiction record; the +3-PC interpretation is falsified if the bytes
  at `$C21338/$C21347/$C21396` are not 3-byte stores.
- **Raw evidence paths:** `mesen/out/rom_C212F0_304.hex` (+ SHA-256 in the
  result section), new `mesen/out/events.log` session header. Session 003
  originals archived at `mesen/out/session003/` (events.log SHA-256
  `bcfc7f4c…a99d03` verified at archive time; resp.txt
  `4405af70…5281cb`).
- **Result:** **Every battle.go claim verified byte-exact** (dump
  `mesen/out/rom_C212F0_304.hex`, SHA-256 `2800f34b…d5d56a`; annotated
  disassembly promoted to
  [02_DISCOVERED_FUNCTIONS.md](../sessions/02_DISCOVERED_FUNCTIONS.md)):
  dispatch `JSR ($131F,X)` at `ROMCPU:$C21300` (pushed return `$1302`
  matches the Session 003 stack captures; the Session 003 record's
  `$C212FF` inference was an off-by-one, corrected); table `$C2131F` =
  `[$1323, $1350]`; HP routine `$C21323` with stores `$C21338`/`$C21347`;
  MP routine `$C21350` over `$3C08`/`$3C30` with `$3C95`-bit-0 death call
  and `LDA #$0080 / JMP $464C` tail; death handler `$C21390` (clears
  `$3A89`, store `$C21396`, `$3EE4,Y` `BIT #$0002` suppression,
  `JMP $0E32` continuation); delta fetch `$C213A7` over `$33E4`/`$33D0`
  with `$FFFF` sentinels. The +3 post-instruction-PC interpretation of
  write-callback captures is confirmed (all three logged PCs are store
  address + 3). New unknowns discovered: `WRAM:+$11A2` bit 7 selects the
  MP path at `$C212F9`; post-dispatch tail `$C21303–$C2131B`
  (`$327C,Y`/`$3018,Y`/`TRB $3419`); fetch gates `$3A3C`,
  `$3A81`&`$3A82`. Bonus: Mesen's log window recorded the ROM header
  parse — Type HiROM, FastROM, Map Mode `$31` — upgrading the
  ROM-identity mapping row.
- **Status:** completed (2026-07-29)
- **Confidence:** Engine code Confirmed (byte-exact dump, deterministic
  read, branch offsets verified arithmetically). Semantic meanings
  (MP vs HP arrays, "death") remain at their prior levels pending live
  observation — code roles are now Confirmed, labels are not.
- **Next action:** promote records (done with this experiment); queue
  live MP observation and `$11A2`/tail investigations as open questions.
