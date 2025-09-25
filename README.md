## Debugging

- Restart after kernel update until `glxinfo | grep "OpenGL version"` works

## TODO

- time must come from engine, not now
- see the current action (debug?)
- spawner: show queue
- keyboard config
- bigger map
- move the map view
- ecs array: on delete, replace by last and update id->index map
- collision based move w smaller accuracy
- perf: look for nil comparison -> more data oriented
- perf: keep a short array of selected and use it instead of `Value.IsSelected`
- perf: component columns and systems

## Bugs

- if B is < 1 cell away, it will go and merge with A
- stops moving if someone settles at destination

## Dependencies

- `sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`
