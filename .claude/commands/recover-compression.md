Use the graphics-archaeologist, dma-tracer, experiment-designer, and go-implementer skills. Follow `.claude/playbooks/RECOVER_COMPRESSION.md`.

Target: $ARGUMENTS

**First establish that the data is compressed at all.** EXP-0050 found 47-52% of a field scene's VRAM present in the ROM verbatim after readiness X1 had asserted for ten units that compression gated it. Run `ff6lab state origin` before assuming a format exists.

Record the decompressed reference and its hash, the ROM source span, the algorithm with a byte-exact falsifier, the Go decoder, its fuzz target, and the golden vector pairing source to output.
