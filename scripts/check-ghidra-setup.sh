#!/usr/bin/env bash
set -euo pipefail

ROOT="${FF6_RE_WORKSTATION:-$HOME/Learning/GitHub/FF6-Reverse-Engineering}"
GHIDRA_DIR="${GHIDRA_INSTALL_DIR:-$ROOT/tools/ghidra/ghidra_12.1.2_PUBLIC}"
WORKSPACE_DIR="${GHIDRA_WORKSPACE_DIR:-$ROOT/workspaces/ghidra}"
ROM_PATH="${FF6_ROM_PATH:-$ROOT/private/roms/Final Fantasy III (USA).sfc}"

fail=0
check_path() {
  local label="$1" path="$2"
  if [[ -e "$path" ]]; then
    printf 'OK   %-22s %s\n' "$label" "$path"
  else
    printf 'MISS %-22s %s\n' "$label" "$path"
    fail=1
  fi
}

check_path "Ghidra install" "$GHIDRA_DIR"
check_path "Ghidra launcher" "$GHIDRA_DIR/ghidraRun"
check_path "Ghidra workspace" "$WORKSPACE_DIR"
check_path "Local ROM" "$ROM_PATH"

if command -v java >/dev/null 2>&1; then
  printf 'OK   %-22s %s\n' "Java" "$(java -version 2>&1 | head -n 1)"
else
  printf 'MISS %-22s %s\n' "Java" "not found"
  fail=1
fi

if [[ -f "$ROM_PATH" ]]; then
  printf 'INFO %-22s %s\n' "ROM SHA-256" "$(shasum -a 256 "$ROM_PATH" | awk '{print $1}')"
fi

printf '\nManual Ghidra checks still required:\n'
printf '  - SNES Loader version matches Ghidra version\n'
printf '  - import format is SNES ROM Loader\n'
printf '  - language is 65816:LE:24:snes\n'
printf '  - memory map contains bank_c0_hirom\n'
printf '  - ROMCPU:$C09B5C and ROMCPU:$C0B593 resolve\n'

exit "$fail"
