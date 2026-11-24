package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const SampleRate = 44100

// Label is a piece of positioned text (replacement for sf.Text).
type Label struct {
	Str   string
	Face  *text.GoTextFace
	Color color.RGBA
	X, Y  float64
}

func (this *Label) Width() float64 {
	w, _ := text.Measure(this.Str, this.Face, 0)
	return w
}

func (this *Label) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(this.X, this.Y)
	op.ColorScale.ScaleWithColor(this.Color)
	text.Draw(screen, this.Str, this.Face, op)
}

type Game struct {
	EnemyList            []*Enemy
	EnemyGenerationSpeed int64
	LastEnemyGenerated   int64
	Player               *Player
	Width                uint
	Height               uint
	ChildsCollected      int
	ChildsCollectedText  *Label
	WelcomeText          *Label
	PlayText             *Label
	StartTime            int64
	EndTime              int64
	TimeText             *Label
	EndingText           *Label
	EndingPointsText     *Label
	RestartText          *Label
	CreditsText          *Label
	State                int
	Font                 *text.GoTextFaceSource
	AudioContext         *audio.Context
	Music                *audio.Player
	Textures             []Texture
	ChildStolen          *audio.Player
}

func NewGame(width uint, height uint) *Game {
	Game := new(Game)
	Game.Width = width
	Game.Height = height

	fontData, err := os.ReadFile("res/fonts/UbuntuMono-R.ttf")
	if err != nil {
		log.Fatal(err)
	}
	Game.Font, err = text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}

	Game.AudioContext = audio.NewContext(SampleRate)
	Game.Music = Game.loadMusic("res/audio/Child_Service_-_Going_to_steal_your_Kids.ogg")
	Game.ChildStolen = Game.loadSound("res/audio/childstolen.wav")

	Game.State = 0
	Game.StartStateInit()
	Game.Textures = NewTextures("res/images/")
	return Game
}

func (this *Game) loadMusic(path string) *audio.Player {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := vorbis.DecodeWithSampleRate(SampleRate, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	player, err := this.AudioContext.NewPlayer(audio.NewInfiniteLoop(stream, stream.Length()))
	if err != nil {
		log.Fatal(err)
	}
	return player
}

func (this *Game) loadSound(path string) *audio.Player {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	stream, err := wav.DecodeWithSampleRate(SampleRate, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	player, err := this.AudioContext.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
	}
	return player
}

// restart plays a sound from the beginning, like sf.Music.Play on a playing stream.
func restart(p *audio.Player) {
	p.Pause()
	if err := p.SetPosition(0); err != nil {
		log.Println(err)
	}
	p.Play()
}

func (this *Game) face(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: this.Font, Size: size}
}

func (this *Game) centered(str string, size float64, c color.RGBA, y float64) *Label {
	l := &Label{Str: str, Face: this.face(size), Color: c, Y: y}
	l.X = (float64(this.Width) - l.Width()) / 2
	return l
}

var black = color.RGBA{0, 0, 0, 255}

func (this *Game) GameStateInit() {
	this.Player = NewPlayer(int(this.Width), int(this.Height), this.Textures)
	this.EnemyList = make([]*Enemy, 0)
	this.EnemyList = append(this.EnemyList, NewEnemy(int(this.Width), this.Textures))
	this.LastEnemyGenerated = time.Now().UnixNano()
	this.EnemyGenerationSpeed = 500000000
	this.ChildsCollected = 0
	this.ChildsCollectedText = &Label{Str: fmt.Sprintf("Childs Collected: %v", this.ChildsCollected), Face: this.face(24), Color: black, X: 20, Y: 20}
	this.StartTime = time.Now().Unix()
	this.TimeText = &Label{Str: fmt.Sprintf("Collection time: %v", time.Now().Unix()-this.StartTime), Face: this.face(24), Color: black, Y: 20}
	this.TimeText.X = float64(this.Width) - 20 - this.TimeText.Width()
	restart(this.Music)
}

func (this *Game) StartStateInit() {
	this.WelcomeText = this.centered("You work now in Child Service and you must collect every child who has no mother!", 36, black, 250)
	this.PlayText = this.centered("Press Enter to play!", 42, black, 350)
}

func (this *Game) endStateInit(message string, restartColor color.RGBA) {
	this.EndTime = time.Now().Unix()
	this.EndingText = this.centered(message, 36, black, 250)
	this.EndingPointsText = this.centered(fmt.Sprintf("You collected %v children in %v seconds.", this.ChildsCollected, this.EndTime-this.StartTime), 36, black, 350)
	this.RestartText = this.centered("Press Enter to restart the game!", 42, restartColor, 450)
	this.CreditsText = this.centered("Created for Ludum Dare 33 Compo by Toni Korpela!", 42, black, 550)
}

func (this *Game) KilledStateInit() {
	this.endStateInit("You were killed by a Father who had a weapon.", black)
}

func (this *Game) MonsterStateInit() {
	this.endStateInit("You were such a monster to try and steal child from a mother.", color.RGBA{255, 0, 0, 255})
}

func (this *Game) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	enter := inpututil.IsKeyJustReleased(ebiten.KeyEnter) || inpututil.IsKeyJustReleased(ebiten.KeyNumpadEnter)

	if this.State == 0 {
		// Opening state
		if enter {
			this.State = 1
			this.GameStateInit()
		}
	} else if this.State == 1 {
		// Game state
		t := time.Now().UnixNano()
		this.TimeText.Str = fmt.Sprintf("Collection time: %v", time.Now().Unix()-this.StartTime)
		if t > this.LastEnemyGenerated+this.EnemyGenerationSpeed {
			this.LastEnemyGenerated = t
			this.EnemyList = append(this.EnemyList, NewEnemy(int(this.Width), this.Textures))
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			this.Player.MoveLeft()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			this.Player.MoveUp()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			this.Player.MoveRight()
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			this.Player.MoveDown()
		}
		for i := 0; i < len(this.EnemyList); i++ {
			if this.EnemyList[i].Collision(this.Player) == true {
				if this.EnemyList[i].GetHasFather() == true && this.EnemyList[i].GetHasMother() == false {
					if this.EnemyList[i].GetFatherHasWeapon() == false {
						restart(this.ChildStolen)
						this.ChildsCollected += 1
						this.ChildsCollectedText.Str = fmt.Sprintf("Childs Collected: %v", this.ChildsCollected)
						this.EnemyList = append(this.EnemyList[:i], this.EnemyList[i+1:]...)
						i--
					} else if this.EnemyList[i].GetFatherHasWeapon() == true {
						// End the game, because the Child Service worker was killed.
						this.KilledStateInit()
						this.State = 2
					}
				} else {
					// End the game and call the Child Service worker a Monster.
					this.MonsterStateInit()
					this.State = 3
				}
			} else {
				this.EnemyList[i].MoveDown()
				if this.EnemyList[i].Destroy(int(this.Height)) == true {
					this.EnemyList = append(this.EnemyList[:i], this.EnemyList[i+1:]...)
					i--
				}
			}
		}
		if this.ChildsCollected > 20 && this.ChildsCollected < 40 {
			this.EnemyGenerationSpeed = 500000000
		} else if this.ChildsCollected > 40 && this.ChildsCollected < 60 {
			this.EnemyGenerationSpeed = 400000000
		} else if this.ChildsCollected > 60 && this.ChildsCollected < 80 {
			this.EnemyGenerationSpeed = 300000000
		} else if this.ChildsCollected > 80 && this.ChildsCollected < 100 {
			this.EnemyGenerationSpeed = 200000000
		} else if this.ChildsCollected > 100 && this.ChildsCollected < 150 {
			this.EnemyGenerationSpeed = 100000000
		} else if this.ChildsCollected > 150 && this.ChildsCollected < 300 {
			this.EnemyGenerationSpeed = 50000000
		} else if this.ChildsCollected > 300 {
			this.EnemyGenerationSpeed = 25000000
		}
	} else if this.State == 2 || this.State == 3 {
		// Ending states (killed / called a monster)
		this.Music.Pause()
		if enter {
			this.State = 0
		}
	}
	return nil
}

func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{44, 176, 55, 255})
	if this.State == 0 {
		this.WelcomeText.Draw(screen)
		this.PlayText.Draw(screen)
	} else if this.State == 1 {
		for _, Enemy := range this.EnemyList {
			Enemy.Draw(screen)
		}
		this.Player.Draw(screen)
		this.ChildsCollectedText.Draw(screen)
		this.TimeText.Draw(screen)
	} else if this.State == 2 || this.State == 3 {
		this.EndingText.Draw(screen)
		this.EndingPointsText.Draw(screen)
		this.RestartText.Draw(screen)
		this.CreditsText.Draw(screen)
	}
}

// Layout keeps the original 1920x1080 world and scales it to the window.
func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(this.Width), int(this.Height)
}
