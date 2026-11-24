package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Shape       *Rect
	WorldWidth  int
	WorldHeight int
}

func NewPlayer(WorldWidth int, WorldHeight int, Textures []Texture) *Player {
	Player := new(Player)
	Player.WorldWidth = WorldWidth
	Player.WorldHeight = WorldHeight
	Player.Shape = NewRect(32, 32, GetTexture("Player.png", Textures).Data)
	x := rand.Intn(WorldWidth)
	if x > 32 {
		x -= 32
	}
	Player.Shape.X = float64(x)
	Player.Shape.Y = float64(WorldHeight - 64)
	return Player
}

func (this *Player) MoveLeft() {
	if this.Shape.X > 0 {
		this.Shape.X -= 10
	}
}

func (this *Player) MoveUp() {
	if this.Shape.Y > 0 {
		this.Shape.Y -= 10
	}
}

func (this *Player) MoveRight() {
	if this.Shape.X < float64(this.WorldWidth-32) {
		this.Shape.X += 10
	}
}

func (this *Player) MoveDown() {
	if this.Shape.Y < float64(this.WorldHeight-32) {
		this.Shape.Y += 10
	}
}

func (this *Player) Draw(screen *ebiten.Image) {
	this.Shape.Draw(screen)
}
