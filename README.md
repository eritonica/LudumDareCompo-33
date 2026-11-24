# LudumDareCompo-33
This is repository for my Ludum Dare 33 competition source code.

## Running

Requires [Go](https://go.dev/dl/) 1.25 or newer. Run from the repository root (the game loads assets from `res/`):

```sh
go run .
```

Or build a binary and keep it next to the `res/` folder:

```sh
go build -o childservice.exe .
./childservice.exe
```

On Linux you also need the usual Ebitengine system packages (e.g. `libgl1-mesa-dev xorg-dev libasound2-dev` on Debian/Ubuntu).

**Controls:** arrow keys / WASD to move, Enter to start or restart, Esc to quit.

The original 2015 version used `bitbucket.org/krepa098/gosfml2` (SFML bindings), which is no longer available; the game has since been ported to [Ebitengine](https://ebitengine.org/) with the same gameplay.
