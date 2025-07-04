## Debugging

- Restart after kernel update until `glxinfo | grep "OpenGL version"` works

## TODO

- atm, stops moving if someone settles at destination
- an entity may spawn other entities with a cost
- spawn order
- bigger map
- move the map view
- bug: select A, gather, move to B, (it stops), gather (no storage)

## Dependencies

- `sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel mesa-libGL-devel libXi-devel libXxf86vm-devel`
