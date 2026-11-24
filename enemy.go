package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	HasChild        bool
	HasFather       bool
	HasMother       bool
	FatherHasWeapon bool
	Shape           *Rect
}

func NewEnemy(WorldWidth int, Textures []Texture) *Enemy {
	Enemy := new(Enemy)
	Enemy.HasChild = true
	if rand.Intn(100) < 80 {
		Enemy.HasMother = true
	} else {
		Enemy.HasMother = false
	}
	if Enemy.HasMother == true {
		if rand.Intn(1) == 1 {
			Enemy.HasFather = true
		} else {
			Enemy.HasFather = false
		}
	} else {
		Enemy.HasFather = true
	}
	if rand.Intn(1000) == 1 {
		Enemy.FatherHasWeapon = true
	} else {
		Enemy.FatherHasWeapon = false
	}
	if Enemy.HasMother == true && Enemy.HasFather == true && Enemy.HasChild == true {
		Enemy.Shape = NewRect(96, 32, GetTexture("ParentsAndChild.png", Textures).Data)
		Enemy.Shape.Color = color.RGBA{128, 128, 128, 255}
	} else if Enemy.HasMother == true && Enemy.HasFather == false && Enemy.HasChild == true {
		Enemy.Shape = NewRect(64, 32, GetTexture("MotherAndChild.png", Textures).Data)
		Enemy.Shape.Color = color.RGBA{128, 0, 128, 255}
	} else {
		Enemy.Shape = NewRect(64, 32, GetTexture("FatherAndChild.png", Textures).Data)
		if Enemy.FatherHasWeapon == true {
			Enemy.Shape.Color = color.RGBA{0, 128, 255, 255}
		} else {
			Enemy.Shape.Color = color.RGBA{0, 128, 128, 255}
		}
	}
	x := rand.Intn(WorldWidth)
	if x > 96 {
		x -= 96
	}
	Enemy.Shape.X = float64(x)
	Enemy.Shape.Y = 0
	return Enemy
}

func (this *Enemy) GetHasFather() bool {
	return this.HasFather
}

func (this *Enemy) GetHasMother() bool {
	return this.HasMother
}

func (this *Enemy) GetFatherHasWeapon() bool {
	return this.FatherHasWeapon
}

func (this *Enemy) MoveDown() {
	this.Shape.Y += 5
}

func (this *Enemy) Collision(Player *Player) bool {
	return Player.Shape.Intersects(this.Shape)
}

func (this *Enemy) Destroy(WorldHeight int) bool {
	return this.Shape.Y > float64(WorldHeight)
}

func (this *Enemy) Draw(screen *ebiten.Image) {
	this.Shape.Draw(screen)
}
