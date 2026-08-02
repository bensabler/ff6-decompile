Use the static-runtime-correlator, 65816-analyst, mesen-operator, experiment-designer, function-recovery, and documentation-curator skills.

Target: $ARGUMENTS

Follow `docs/static-analysis/GHIDRA_MESEN_WORKFLOW.md`.

Create one bounded static/runtime correlation record. Start from existing evidence or a specific ROMCPU address. Treat Ghidra disassembly, function boundaries, labels, references, and decompiler output as hypotheses until confirmed by Mesen evidence.

Required deliverables:

- one correlation record based on `.claude/templates/STATIC_CORRELATION.md`;
- explicit CPU-address, file-offset, and Ghidra-address mapping;
- processor-state assumptions at each disassembly entry;
- Mesen trace or breakpoint evidence;
- competing interpretations and a falsification test;
- links to any affected experiment, function, variable, or discovery record;
- exact next action.

Do not broaden into neighboring systems. Register incidental observations through `/census-observations`.
