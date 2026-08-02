Use the static-runtime-correlator, 65816-analyst, context-manager, and documentation-curator skills.

Bootstrap the local Ghidra integration without beginning a new experiment.

Required checks:

1. read `docs/research/ROM_IDENTITY.md` and verify the local ROM SHA-256;
2. verify Ghidra and the SNES extension versions;
3. verify the ROM was imported with the SNES ROM Loader, not Raw Binary;
4. verify the language is `65816:LE:24:snes`;
5. verify the memory map contains canonical HiROM blocks such as `bank_c0_hirom`;
6. verify `ROMCPU:$C09B5C` resolves in Ghidra;
7. verify `ROMCPU:$C0B593` resolves in Ghidra;
8. identify the external Ghidra project location and ensure it is outside Git;
9. record unresolved 65816 processor-state assumptions;
10. write or update `local_artifacts/static-analysis/ghidra-environment.md`.

Do not rename functions, apply broad disassembly, or promote Ghidra output into canonical records during bootstrap.
