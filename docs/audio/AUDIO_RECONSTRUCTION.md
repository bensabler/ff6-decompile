# Audio Reconstruction

The SNES audio subsystem is a separate computer: SPC700 CPU, ARAM, timers, ports, and DSP.

## Reconstruction layers

1. Main CPU command/upload protocol
2. SPC700 driver
3. Sequence interpreter
4. Music and sound-effect data
5. Instrument/sample directory
6. BRR sample streams
7. DSP voice and echo behavior
8. Timing and final rendering

A captured WAV proves only that a cue sounded under one state. It does not prove sequence or sample provenance.
