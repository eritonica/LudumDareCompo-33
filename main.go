package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WorldWidth  = 1920
	WorldHeight = 1080
)

func main() {
	// Resources are loaded relative to the working directory; if res/ isn't
	// there, fall back to the directory containing the executable.
	if _, err := os.Stat("res"); err != nil {
		if exe, err := os.Executable(); err == nil {
			_ = os.Chdir(filepath.Dir(exe))
		}
	}

	Game := NewGame(WorldWidth, WorldHeight)

	// Fit the 1920x1080 world into the monitor; the window is resizable.
	scale := 1.0
	if mw, mh := ebiten.Monitor().Size(); mw > 0 && mh > 0 {
		scale = min(1.0, 0.9*float64(mw)/WorldWidth, 0.9*float64(mh)/WorldHeight)
	}
	ebiten.SetWindowSize(int(WorldWidth*scale), int(WorldHeight*scale))
	ebiten.SetWindowTitle("Child Service")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)

	if err := ebiten.RunGame(Game); err != nil {
		log.Fatal(err)
	}
}
