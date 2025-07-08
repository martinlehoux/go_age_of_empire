## Debugging

- Restart after kernel update until `glxinfo | grep "OpenGL version"` works

## TODO

- atm, stops moving if someone settles at destination
- spawn may register a target, on spawn unit will do as a right click on target (move, gather, ...)
- time must come from engine, not now
- selection priority: selecting units + building -> select only units
- see the current action (debug?)
- keyboard config
- spawn order
- bigger map
- move the map view

## Bugs
- select A, gather, move to B, (it stops), gather (no storage)
- if B is < 1 cell away, it will go and merge with A

## Dependencies

- `sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`
