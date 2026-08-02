# Compression format recovery

## Required inputs

- exact target or queued question;
- ROM identity;
- a preserved capture holding the **decompressed output**;
- emulator identity when runtime evidence is used;
- current checkpoint;
- evidence and stopping standards.

## Procedure

1. **Establish that the data is compressed at all.** Run `ff6lab state origin`
   over the captured region first. EXP-0050 found 47-52% of a field scene's
   VRAM present in the ROM verbatim, after readiness X1 had asserted for ten
   units that compression gated it. A verbatim span needs a slice, not a
   decoder.
2. Freeze the decompressed reference and hash it.
3. Locate the compressed source: DMA trace, staging-buffer ancestry, or a
   differential search over unowned ROM spans (`ff6lab rom gaps`).
4. Read the decompressor, statically or at runtime. Prefer reading the routine
   over guessing the format; EXP-0051 spent four encodings guessing and found
   nothing.
5. State the algorithm with a falsifier that is **byte-exact reproduction** of
   the captured output from the ROM span.
6. Implement a Go decoder with table-driven tests, a fuzz target for malformed
   input, and the golden vector pairing source to output.
7. Extend `internal/extract` and land the asset with manifest provenance.
8. Record what the format does **not** cover, and which regions remain
   unexplained.

## Required outputs

- experiment and discovery records, with the falsifier stated;
- decompressed reference hash and ROM span;
- Go decoder, tests, fuzz target, golden vector;
- extractor and manifest entry;
- confidence and alternatives;
- exact next action;
- checkpoint when the unit stops.

## Stopping rule

A negative result is a result. If no algorithm reproduces the capture
byte-exactly, record that, list the surviving alternatives, and **do not
implement on a hypothesis**.
