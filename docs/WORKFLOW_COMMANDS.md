# Workflow Command Reference

## `/bootstrap-v4`

**Use:** Once after installing or upgrading to Version 4.

**Performs:**

- inventories the repository;
- preserves earlier research;
- checks directory structure;
- records Go, Git, Claude Code, Mesen, and OS versions;
- computes the ROM hash locally;
- verifies Git ignore boundaries;
- initializes dashboards and indexes;
- creates a migration report and checkpoint.

**Does not:** Begin a new reverse-engineering experiment.

---

## `/resume-session`

**Use:** At the beginning of every fresh Claude session.

**Performs:**

- reads the latest checkpoint;
- reads current focus, queue, blockers, and open contradictions;
- checks the working tree;
- identifies interrupted work;
- summarizes the exact next action.

**Rule:** Interrupted work takes precedence over new work.

---

## `/orchestrate`

**Use:** Let Claude choose the next bounded research unit.

**Performs:**

- scores queued tasks by information gain, dependencies, cost, and reproducibility;
- selects one task;
- writes an experiment plan;
- delegates to specialists;
- executes until a stopping condition;
- updates all durable records.

**Do not use:** When you already know the exact target. Use the specific command instead.

---

## `/checkpoint`

**Use:** Before stopping, changing models, compacting context, installing files, or interrupting Mesen.

**Performs:**

- records active question and state;
- lists evidence and changed files;
- captures exact next action;
- records test and Git status;
- updates `dashboards/CURRENT_FOCUS.md`.

---

## `/investigate-function <CPU address or function ID>`

Recovers boundaries, callers, callees, inputs, outputs, side effects, control flow, and candidate purpose. Requires trace evidence and at least one falsification experiment.

Example:

```text
/investigate-function C10DF3
```

---

## `/investigate-variable <address or variable ID>`

Recovers reads, writes, owner, lifetime, width, stride, reset behavior, and relationships.

Example:

```text
/investigate-variable WRAM+$2EB5
```

---

## `/reconstruct-struct <record or address range>`

Builds an offset table from repeated access patterns and experiments. Unknown fields remain unknown.

---

## `/trace-caller <function>`

Traces direct and indirect caller paths, call conditions, and upstream semantics.

---

## `/trace-dma <scene, asset, or transfer>`

Records DMA channel, mode, source, destination, length, trigger routine, buffers, and source provenance.

---

## `/capture-graphics <target>`

Captures deterministic runtime graphics evidence: frame, PPU state, VRAM, CGRAM, OAM, layers, viewers, and hashes.

A capture is evidence—not yet a reconstructed source asset.

---

## `/reconstruct-sprite <target>`

Recovers tiles, palettes, OAM/layout, offsets, flips, priorities, frame identity, and animation relationships. Produces a local preview and metadata manifest.

---

## `/reconstruct-tileset <target>`

Recovers tile data, palettes, tilemap entries, layer dimensions, flips, priorities, loading paths, and compression.

---

## `/reconstruct-background <target>`

Handles complete multi-layer scenes, including scroll, mode, tilemap dimensions, color math, windowing, and runtime validation.

---

## `/recover-palette <target>`

Traces CGRAM state to loading tables or transformation routines and documents color format and palette-bank usage.

---

## `/validate-graphics <asset ID>`

Compares reconstructed output to deterministic runtime evidence. Reports exact pixel, palette, layout, and layer differences.

---

## `/investigate-audio <cue>`

Creates a controlled audio trigger and separates:

- CPU/APU messaging;
- SPC700 driver behavior;
- sequence data;
- BRR samples;
- DSP registers;
- timing and voice allocation.

---

## `/recover-brr <sample or address>`

Recovers BRR block boundaries, directory entries, loop/end flags, filters, range, pitch relationship, and verified PCM decoding.

---

## `/recover-sequence <cue>`

Recovers music or SFX command stream, control flow, timing, instruments, loops, and driver interpretation.

---

## `/trace-spc-command <game action>`

Traces the main CPU command through ports to SPC driver dispatch and resulting sequence/sample activity.

---

## `/validate-audio <asset ID>`

Prefers command, event, sample, and DSP-state comparison. Waveform comparison is secondary and must document latency and rendering assumptions.

---

## `/validate-hypothesis <hypothesis ID>`

Designs the smallest discriminating experiment. A successful result raises confidence only by the documented criteria.

---

## `/resolve-contradiction <topic>`

Freezes implementation, gathers conflicting claims and source evidence, runs a discriminating experiment, and records the resolution.

---

## `/implement-discovery <discovery ID>`

Converts only sufficiently proven findings into Go:

- code;
- tests;
- fuzz targets where appropriate;
- docs;
- manifests;
- architecture updates.

---

## `/run-quality-gates`

Runs formatting, tests, vet, schema checks, repository-boundary audit, broken-reference audit, and evidence consistency checks.

---

## `/audit-project`

Performs a deeper read-only-first review of:

- unsupported claims;
- stale dashboards;
- inconsistent addresses;
- missing provenance;
- package architecture;
- test gaps;
- legal-boundary violations.

---

## `/weekly-review`

Recalculates status, merges duplicate questions, flags stale hypotheses, audits documentation drift, updates milestones, and reprioritizes the queue.

---

## `/prepare-release <version>`

Runs release checks and produces a release report. It must fail if ROM-derived bytes are tracked or quality gates fail.

---

## `/bootstrap-ghidra`

Verifies the local Ghidra/SNES Loader environment, ROM identity, HiROM mapping, known address resolution, external workspace location, and current 65816 state limitations. It does not begin research or rename functions.

---

## `/correlate-static-runtime <ROMCPU address or question>`

Creates one bounded correlation between Ghidra static leads and Mesen runtime evidence. Requires explicit processor-state assumptions, a falsification test, and a correlation record.

Example:

```text
/correlate-static-runtime ROMCPU:$C0B593
```

---

## `/export-ghidra-symbols`

Exports reviewed labels/bookmarks to local text or JSON for audit. It never commits the Ghidra database or promotes names automatically.
