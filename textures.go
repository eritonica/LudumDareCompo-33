package main

import (
	"image/color"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Texture struct {
	Name string
	Data *ebiten.Image
}

func NewTextures(dirname string) []Texture {
	dir, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	Textures := make([]Texture, 0)
	for _, Entry := range dir {
		if Entry.IsDir() {
			continue
		}
		var Tex Texture
		Tex.Name = Entry.Name()
		Tex.Data, _, err = ebitenutil.NewImageFromFile(filepath.Join(dirname, Entry.Name()))
		if err != nil {
			log.Fatal(err)
		}
		Textures = append(Textures, Tex)
	}
	return Textures
}

func GetTexture(TextureName string, Textures []Texture) Texture {
	for Texture := range Textures {
		if Textures[Texture].Name == TextureName {
			return Textures[Texture]
		}
	}
	return Textures[0]
}

// Rect is a textured, tinted rectangle (replacement for sf.RectangleShape).
type Rect struct {
	X, Y          float64
	Width, Height float64
	Color         color.RGBA
	Texture       *ebiten.Image
}

func NewRect(width, height float64, tex *ebiten.Image) *Rect {
	return &Rect{Width: width, Height: height, Color: color.RGBA{255, 255, 255, 255}, Texture: tex}
}

func (this *Rect) Intersects(other *Rect) bool {
	return this.X < other.X+other.Width && other.X < this.X+this.Width &&
		this.Y < other.Y+other.Height && other.Y < this.Y+this.Height
}

func (this *Rect) Draw(screen *ebiten.Image) {
	if this.Texture == nil {
		return
	}
	b := this.Texture.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(this.Width/float64(b.Dx()), this.Height/float64(b.Dy()))
	op.GeoM.Translate(this.X, this.Y)
	op.ColorScale.ScaleWithColor(this.Color)
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(this.Texture, op)
}
