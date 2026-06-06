# Examples Catalog — Architecture

## 0. Solution (TL;DR)

Add an `examples/` package with a registry of named scenarios. Each example is a function that populates a `*Game` with entities, positions, etc. The `example` subcommand selects which scenario to load. When no subcommand is given, the `"default"` example runs. No new dependencies, no structural changes to the engine.

```
go run . example                # list available examples
go run . example combat         # launches a combat demo
go run .                        # launches the default world (via registry)
```

## 1. Context & scope

### Current landscape

The game has an ECS engine, a camera system, pathfinding, resource gathering, unit spawning, patrol orders, and an Input abstraction layer. But currently, the "world" is hardcoded in `main()`:

```go
// main.go — hardcoded world setup
ironMineBuilder := NewEntityBuilder()...
game.Append(ironMineBuilder.WithPosition(...).Build())
townCenterBuilder := NewEntityBuilder()...
townCenter := game.Append(townCenterBuilder.WithPosition(...).Build())
// ...
```

Every feature addition bloats `main()` further. There is no way to:

- Demonstrate a specific feature in isolation
- Quickly iterate on a scenario without rebuilding everything
- Show the game to someone without them seeing a crowded test world

### What is being built

An **examples catalog** — a lightweight framework under `examples/` where each scenario is a self-contained Go file that builds a world. The CLI entry point (`main.go`) dispatches to the chosen example, or runs `"default"` when no subcommand is given.

Future examples might demonstrate:

- Combat between two groups of units
- Resource gathering loop (mine → carry → deposit)
- Patrol routes along a border
- Mass spawning from multiple town centers
- A "tech demo" with every entity type on screen

## 2. Goals & non-goals

### Goals

- **CLI dispatch:** `go run . example combat` loads the combat scenario; `go run .` → `"default"` example from the registry
- **Self-contained examples:** Each example file owns its entity setup, positioning, and initial camera position — including `"default"`
- **Feature independence:** Examples live outside `main.go` and don't require changes to core engine files
- **Discoverability:** `go run . example` (no name) prints all registered examples with descriptions
- **Minimal wiring:** Zero boilerplate beyond adding a file and registering it

### Non-goals

- Not a scripting system — examples are compiled Go, not JSON/YAML data files
- Not a map editor — the goal is code-driven scenario setup, not a GUI tool
- Not a save/load system — examples are static worlds, not persisted state
- Not a modding API — only Go code, no plugin boundary

## 3. Design

### Architecture overview

```mermaid
flowchart LR
    A[main.go] -->|example combat| B{examples.Registry}
    B -->|lookup "combat"| C[combat.Build]
    C -->|Append entities| D[Game]
    D -->|ebiten.RunGame| E[Rendered world]
```

### Registry

A single map that maps example name → builder function:

```
examples/registry.go
  type Builder func(*Game, BuilderDeps)
  type BuilderDeps struct { UnitBuilder EntityBuilder; ... }
  var Registry = map[string]struct{
    Name        string
    Description string
    Build       Builder
  }
```

The `BuilderDeps` struct carries shared resources (font, pre-built entity builders) so each example doesn't recreate them. The example only calls `game.Append(...)` and optionally sets `game.Camera`.

### Example structure (combat)

```go
// examples/combat.go
package examples

import "age_of_empires/..."

func init() {
    Registry["combat"] = Entry{
        Name:        "Combat",
        Description: "Two groups of units fight each other",
        Build:       buildCombat,
    }
}

func buildCombat(g *Game, deps BuilderDeps) {
    // Red team
    redBuilder := deps.UnitBuilder.WithSolid(NewFilledCircleImage(100, color.RGBA{0xff,0,0,0xff}))
    for x := 0; x < 5; x++ {
        g.Append(redBuilder.WithPosition(physics.Point{X: 400 + x*120, Y: 800}).Build())
    }
    // Blue team
    blueBuilder := deps.UnitBuilder.WithSolid(NewFilledCircleImage(100, color.RGBA{0,0,0xff,0xff}))
    for x := 0; x < 5; x++ {
        g.Append(blueBuilder.WithPosition(physics.Point{X: 400 + x*120, Y: 1200}).Build())
    }
    // Camera position
    g.Camera = Camera{X: 800, Y: 1000, Zoom: 1.0}
}
```

### CLI dispatch

In `main.go`, before creating the Game:

1. Look at `os.Args` for the subcommand pattern: `example` or `example <name>`
2. If `go run . example` (no name), print the registry and exit
3. If `go run . example <name>`, look up the name in `examples.Registry`:
   - If found, call the builder with the game pointer before `ebiten.RunGame`
   - If not found, print error + available names and exit
4. If no `example` subcommand (`go run .`), dispatch to `"default"` via the same registry

This avoids mixing concerns: flags like `--profile` are truly cross-cutting, while `example` is a mode selector — naturally a subcommand.

### Package structure

```
go_age_of_empire/
├── examples/
│   ├── registry.go       # Registry map, Builder type, BuilderDeps
│   ├── combat.go         # Combat scenario
│   ├── gathering.go      # Resource gathering scenario
│   └── patrol.go         # Patrol route scenario
├── main.go               # CLI dispatch + default fallback
├── camera.go
├── ...
```

### Flow diagram

```mermaid
flowchart TD
    Start([os.Args]) --> Subcmd{os.Args[1] == "example"?}
    Subcmd -->|No| Lookup["Look up 'default' in examples.Registry"]
    Subcmd -->|Yes| HasName{os.Args[2] exists?}
    HasName -->|No| List["Print registry and exit"]
    HasName -->|Yes| Lookup["Look up name in examples.Registry"]
    Lookup --> Found{Found?}
    Found -->|No| Error["Print error + available names, exit"]
    Found -->|Yes| Build["Call example.Build(game, deps)"]
    Build --> Run["ebiten.RunGame(game)"]
    Default --> Run
```

### BuilderDeps vs global state

Shared resources (font, entity templates) are passed through `BuilderDeps` instead of being global variables. This keeps examples testable and independent.

```go
type BuilderDeps struct {
    FaceSource   *text.GoTextFaceSource
    UnitBuilder  EntityBuilder
    // future: BuildingBuilder, TreeBuilder, etc.
}
```

`main.go` constructs the `BuilderDeps` once and passes it to every example.

## 4. Alternatives

### A — Data-driven (JSON/YAML) world definitions

A JSON file per example describing entities, positions, and properties. A generic loader iterates over the data and calls `game.Append()`.

| Pro | Con |
|---|---|
| Non-programmers can author scenarios | Can't express logic (patrol routes, timers, AI behavior) |
| Loadable at runtime without recompilation | Type safety requires schema validation |
| Easy to share and version | Entity composition is awkward in flat JSON |

**Verdict:** Rejected for now. The examples catalog is for developers demonstrating features, and the feature logic (e.g., "spawn 5 units then patrol them") is too procedural for a data format. Data files could be added later for map data.

### B — Standalone binaries per example

Each example gets its own `main.go` in a subdirectory: `examples/combat/main.go`, `examples/gathering/main.go`.

| Pro | Con |
|---|---|
| Truly independent, no dispatch code needed | Massive duplication of engine setup (font, icon, window, Ebiten boilerplate) |
| Each can have its own `go.mod` | Can't share entity builders or utility code without a library package |
| Clean separation | `go run ./examples/combat` is more typing than `go run . example combat` |

**Verdict:** Rejected. The duplication cost outweighs the isolation benefit. A single binary with CLI dispatch is more ergonomic.

### C — Plugin system (Go plugins / shared objects)

Each example is compiled as a `.so` plugin loaded at runtime via `plugin.Open`.

| Pro | Con |
|---|---|
| Examples can be added without recompiling the main binary | Linux-only, fragile ABI, not supported on all platforms |
| | Debugging is painful |
| | Overkill for first-party examples shipped in the same repo |

**Verdict:** Rejected. Go's plugin system has too many caveats and zero advantage for in-repo examples.

## 5. Open questions

1. **Should the default world (current `main()` code) also be an example?**
   - ✅ **Decision:** Yes — move it to `examples/default.go` and register it as `"default"`.
   - When `go run .` (no subcommand), dispatch to `"default"` via the same registry.
   - Single code path: every mode of execution goes through `examples.Registry`.
   - `main.go` has zero world-building logic.

2. **Should examples be able to add new ECS components?**
   - ✅ **Decision:** No. New components belong in `ecs/component.go`.
   - Examples only wire existing components into the world.

3. **How should the camera initial position be set?**
   - ✅ **Decision:** Add `Camera.SetCenter(x, y)` for clarity. Examples call `deps.SetCamera(game, x, y)` or set `game.Camera.X/Y` directly.

4. **Should `example` without an argument or with an unknown name print the list?**
   - ✅ **Decision:** No argument → list. Unknown name → error + list.

5. **How are build resources (entity builders with specific colors/sizes) shared?**
   - ✅ **Decision:** `BuilderDeps` holds a base `UnitBuilder`. Examples clone it and call `.WithSolid(...)` to customize color/size.

## 6. Parties involved

- **Author:** Martin (implementing the examples)
- **Consumers:** Anyone who wants to quickly see a feature in action (demo purposes, debugging, onboarding)
- **Future contributors:** Adding a new example should be a single-file change with no engine modifications

## 7. Cross concerns

| Concern | Mitigation |
|---|---|
| **CLI hygiene** | Use `os.Args` scanning (no external flag dependency needed); subcommand avoids conflating modes with flags. Print usage on unknown subcommand |
| **Example isolation** | Each example is a `func(*Game, BuilderDeps)`. No mutable global state is shared between examples |
| **Testability** | Examples can be tested by calling `example.Build(game, deps)` and asserting on `game.Mask`, `game.Position`, etc. |
| **Observability** | `go run . example` lists examples; unknown names print error + list. The active example name is logged at startup |
