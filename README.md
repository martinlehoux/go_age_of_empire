Run `go run . example` to browse and launch curated gameplay scenarios (combat, gathering, etc.)  

## Debugging

- Restart after kernel update until `glxinfo | grep "OpenGL version"` works

## Testing

- **Integration tests** Use `UpdateSimulation(Input{...})` to drive ticks with simulated inputs — no ebiten runtime needed.
- **Performance benchmarks**

## TODO

- [ ] see the current action (debug?)
- [ ] spawner: show queue
- [ ] keyboard config
- [ ] display selected entity name in bottom banner
- [ ] ecs array: on delete, replace by last and update id->index map
- [ ] collision based move w smaller accuracy
- [ ] perf: look for nil comparison -> more data oriented
- [ ] perf: keep a short array of selected and use it instead of `Value.IsSelected`
- [ ] perf: component columns and systems
- [ ] ecs: bitfield for fast checking
- [ ] ecs: rename to systems
- [ ] perf: measure frame budget
- [ ] getMoveMap could be a maintained index
- [ ] entityAt could be a maintained index

## Bugs

- [ ] spawner spawns 2 units instead of 1

## Dependencies

- `sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`
