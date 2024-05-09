package main

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	g.startTime = time.Now()
	g.currentTime = g.startTime
}

func (g *Game) Update() error {
	t := time.Now()
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if (len(g.keys)) == 1 {
		if g.keys[0] == ebiten.KeyArrowUp {
			g.level.NewOrientation(180)
		} else if g.keys[0] == ebiten.KeyArrowDown {
			g.level.NewOrientation(0)
		} else if g.keys[0] == ebiten.KeyArrowLeft {
			g.level.NewOrientation(270)
		} else if g.keys[0] == ebiten.KeyArrowRight {
			g.level.NewOrientation(90)
		} else if g.keys[0] == ebiten.KeyR {
			g.restart = true
		}
	}
	if t.Sub(g.currentTime) > 200000000 {
		g.currentTime = t
		if g.restart {
			g.restart = false
			g.level.Restart()
		} else {
			g.level.MoveWorms()
		}
		g.needsDraw = true
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.needsDraw {
		wormPositions := g.level.GetWormPositions()

		//op.GeoM.Reset()

		for i := 0; i < len(wormPositions); i++ {
			//op.GeoM.Translate(float64(0), float64(0))
			if i >= 0 {
				op := &ebiten.DrawImageOptions{}
				//op.GeoM.Reset()
				op.GeoM.Translate(float64(wormPositions[i].x*100), float64(wormPositions[i].y*100))

				//op.GeoM.Translate(float64(100), float64(100))
				screen.DrawImage(g.wormPicture, op)
			}
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.level.candy.x*100), float64(g.level.candy.y*100))
		screen.DrawImage(g.candyPicture, op)
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
	DoesItWork()

}
