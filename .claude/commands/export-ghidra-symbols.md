Use the static-runtime-correlator and documentation-curator skills.

Export reviewed Ghidra labels and bookmarks into a text or JSON artifact under `local_artifacts/static-analysis/exports/`.

Rules:

- never export the ROM or Ghidra database into Git;
- include Ghidra version, extension version, ROM SHA-256, and export timestamp;
- distinguish user labels from auto-generated labels;
- exclude speculative semantic names unless they cite a correlation, function, experiment, or discovery record;
- do not overwrite canonical indexes automatically;
- produce a review report listing additions, conflicts, and unsupported names.
