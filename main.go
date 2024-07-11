package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"

	"log"
	"os"
	"time"

	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	GAME_WAITING = iota
	GAME_RUNNING
)

type Game struct {
	level        *Level
	startTime    time.Time
	currentTime  time.Time
	needsDraw    bool
	keys         []ebiten.Key
	wormPicture  *ebiten.Image
	candyPicture *ebiten.Image
	restart      bool
	gameState    int
	writeScore   bool
}

func NewGame() *Game {
	g := &Game{}
	g.Initialize()
	return g
}

func (g *Game) Initialize() {

	reader, err := os.Open("mato.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	wormImg, _, _ := image.Decode(reader)

	g.wormPicture = ebiten.NewImageFromImage(wormImg)

	candyReader, err := os.Open("candy.png")
	if err != nil {
		log.Fatal(err)
	}
	defer candyReader.Close()

	candyImg, _, _ := image.Decode(candyReader)

	g.candyPicture = ebiten.NewImageFromImage(candyImg)

	level := NewLevel()
	level.AddWorm("Mato1")
	level.NewCandy()
	g.level = level

	g.gameState = GAME_WAITING

	g.WaitForStart()
}

func (g *Game) WaitForStart() {
	g.gameState = GAME_WAITING
}

func (g *Game) StartGame() {
	g.level.score = 0
	g.gameState = GAME_RUNNING
	g.startTime = time.Now()
	g.currentTime = g.startTime
}

func (g *Game) Update() error {
	t := time.Now()
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	interestingKeys := [5]ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight, ebiten.KeyR}
	foundKey := ebiten.KeyMax
	// find the last pressed valid key, so mashing the buttons on keyboard will move the worm
	if (len(g.keys)) >= 1 {
		for i := len(g.keys) - 1; i >= 0 && foundKey == ebiten.KeyMax; i-- {
			for j := range interestingKeys {
				if g.keys[i] == interestingKeys[j] {
					foundKey = g.keys[i]
					break
				}
			}
		}
		if !g.HasGameEnded() {
			if foundKey == ebiten.KeyArrowUp {
				g.level.NewOrientation(180)
			} else if foundKey == ebiten.KeyArrowDown {
				g.level.NewOrientation(0)
			} else if foundKey == ebiten.KeyArrowLeft {
				g.level.NewOrientation(270)
			} else if foundKey == ebiten.KeyArrowRight {
				g.level.NewOrientation(90)
			} else if foundKey == ebiten.KeyR {
				g.restart = true
			}
		} else {
			if foundKey == ebiten.KeyR {
				g.restart = true
			}
		}

	}
	if g.gameState == GAME_RUNNING && g.HasGameEnded() {
		g.gameState = GAME_WAITING
		g.writeScore = true
	} else if g.gameState == GAME_RUNNING && t.Sub(g.currentTime) > 200000000 {
		g.currentTime = t
		if g.restart && foundKey != ebiten.KeyR { //Wait until key is released
			g.restart = false
			g.level.Restart()
		} else {
			g.level.MoveWorm()
		}
		g.needsDraw = true
	} else if g.gameState == GAME_WAITING && g.restart {
		// Wait until key is released
		if foundKey != ebiten.KeyR {
			g.restart = false
			g.level.Restart()
			g.StartGame()
		}

	}
	return nil
}

func (g *Game) HasGameEnded() bool {
	if g.level.worm == nil {
		return true
	} else {
		return false
	}
}

func (g *Game) DrawText(screen *ebiten.Image, str string, x float64, y float64) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 0x80})

	fontBaseSize := 20.0
	scale := 1.0
	text.Draw(screen, str, &text.GoTextFace{
		Source: s,
		Size:   fontBaseSize * float64(scale),
	}, op)

	op.GeoM.Reset()
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.gameState == GAME_RUNNING {
		if g.needsDraw {
			wormPositions := g.level.GetWormPositions()

			for i := 0; i < len(wormPositions); i++ {
				if i >= 0 {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(wormPositions[i].x*100), float64(wormPositions[i].y*100))
					screen.DrawImage(g.wormPicture, op)
				}
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.level.candy.x*100), float64(g.level.candy.y*100))
			screen.DrawImage(g.candyPicture, op)
		}
	} else if g.gameState == GAME_WAITING {
		if g.writeScore {
			scoreStr := fmt.Sprintf("Score: %d", g.level.score)
			g.DrawText(screen, scoreStr, 410, 400)
		}
		g.DrawText(screen, "Press r to begin", 350, 480)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Unnamed worm game")

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
