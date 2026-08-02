# ADR-0001: Rendering host for the playable demo

- **Status:** Accepted
- **Date:** 2026-08-02
- **Context:** [DEMO-0001](../demo/DEMO-0001-new-game-to-whelk.md), Unit 6
- **Supersedes:** nothing

## Context

The repository had **zero third-party dependencies** and every quality gate ran
on a bare container with no system libraries. DEMO-0001 requires a window,
keyboard input, and eventually audio output, none of which the standard library
provides.

The decision is therefore not only "which library" but "what does adding the
first dependency in this project's history cost, and can it be undone".

## Decision

Adopt **Ebitengine v2** (`github.com/hajimehoshi/ebiten/v2` v2.9.9,
Apache-2.0), under two constraints that matter more than the choice itself:

1. **Confined to one leaf package.** `internal/engine/ebitenhost` is the only
   package permitted to import it. `audit.CheckImportBoundaries` enforces this
   as a *direct-import* rule and fails the build otherwise. Replacing the host
   means rewriting one package.

2. **Behind a `gui` build tag.** The default build — and therefore
   `go build ./...`, `go vet ./...`, `go test ./...`, and every CI gate except
   one dedicated job — compiles no third-party code at all.

The headless host (`internal/engine/headless`, standard library only) is a
**co-equal permanent host**, not a stepping stone. It is the authoritative one
for tests.

## Why the build tag is the design, not a fallback

Measured on 2026-08-02 against ebiten v2.9.9, building a minimal
`ebiten.RunGame` program:

| Target | `CGO_ENABLED=0` |
|---|---|
| `darwin/amd64` | **fails** — cgo required |
| `linux/amd64` | **fails** — cgo required (`undefined: glfw.Window`) |
| `windows/amd64` | builds |
| `js/wasm` | builds |

Ebitengine's `purego` work removes cgo on some platforms, but **not** on the
two this project actually develops and tests on. Without the tag, the standard
gates would acquire a system OpenGL/X11 dependency, and CI could no longer run
the demo's test surface on a bare container. The tag keeps that cost inside one
opt-in job.

Verified after adoption: `go list -deps ./cmd/ff6demo` on the default build
links **no** third-party package.

## Alternatives considered

| Option | License | Distribution cost | Verdict |
|---|---|---|---|
| **Ebitengine v2** | Apache-2.0 | cgo on macOS/Linux; cgo-free on Windows and js/wasm; no system package to install on macOS | **Accepted** |
| raylib-go (`gen2brain/raylib-go`) | zlib | cgo mandatory on every platform, vendoring a large C library into every build | Rejected. A C toolchain for every contributor and CI job, buying no capability this demo needs |
| go-sdl2 (`veandco/go-sdl2`) | BSD-3 | cgo **plus** a system SDL2 install (brew/apt) | Rejected. Worst distribution story: it breaks clone-and-build. Its finer control buys nothing over one integer-scaled blit |
| Pure stdlib, headless only | n/a | none | Rejected *as the only host* — no window, input, or audio path. **Adopted as the second host**, which is the part that matters |

License compatibility: Apache-2.0 is compatible with this repository's MIT
licence. Ebitengine is actively maintained and widely used for 2D games.

## Consequences

**Accepted:**

- The first `go.sum` in the project's history. The version is pinned;
  `go mod verify` belongs in the pre-commit check set.
- One extra CI job installing `libgl1-mesa-dev xorg-dev libasound2-dev` and
  building `-tags gui`. Without it the GUI path would leave the gates entirely
  and could rot unnoticed.
- A `//go:build gui` code path that the default `go vet` and `go test` do not
  cover. The dedicated job compiles it; behaviour is exercised through the
  shared `engine.Machine`, which the headless tests cover fully.

**Gained:**

- `GOOS=js/wasm` builds today, so a browser demo later needs no code change.
- The dependency is provably optional: the headless host and every test landed
  and passed *before* ebiten entered `go.mod` (commit `5a8ffd8`, one commit
  before this one). That ordering was deliberate.

## The abstraction's acceptance criterion

Ebitengine's `Game` interface is `Update`/`Draw`/`Layout`. The adapter is a
thin translation to `engine.Machine`, holding no game state.

**The moment domain code wants an ebiten type, this abstraction has failed.**
That is why the boundary is a mechanical check rather than a comment: the rule
was verified to fire by temporarily importing ebiten from
`internal/game/scenes`, which produced

```text
import-boundaries: internal/game/scenes -> github.com/hajimehoshi/ebiten/v2 —
only internal/engine/ebitenhost may touch the rendering host …
```

## Host responsibilities, and what the host must never be

The host owns integer scaling (largest whole-number fit, never fractional —
a fractional scale resamples the framebuffer and destroys the pixel grid),
nearest-neighbour filtering, letterboxing, device-pixel layout on HiDPI
displays, and the keyboard map.

The host is **not** authoritative for tests. Ebitengine may call `Update` more
or fewer times than `Draw` and pauses both on focus loss unless told otherwise.
The simulation stays exact because `Tick` is 1:1 with `Update`, but any
frame-for-frame comparison uses `internal/engine/headless`.

## Revisiting

Reconsider if: Ebitengine's licence or maintenance changes; the demo needs
rendering the host cannot express; or the cgo requirement blocks a target that
matters. Reversal cost is one package plus the `go.mod` entry.
