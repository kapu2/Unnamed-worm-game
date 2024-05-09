package main

import (
	"fmt"
	"math/rand"
)

type Level struct {
	name    string
	xcoords int
	ycoords int
	worm    *Worm
	candy   *Coordinate
}

func DoesItWork() {
	fmt.Print("toimiiko ollenkaan")
}

func NewLevel() *Level {
	theLevel := &Level{}
	theLevel.xcoords = 10
	theLevel.ycoords = 10

	theLevel.name = "theworld"
	return theLevel
}

func AbsoluteValue(f float64) float64 {
	if f < 0.0 {
		return -f
	} else {
		return f
	}
}

func (level *Level) AddWorm(wormName string) {
	space := level.GetFreeSpaceForWorm()
	howRight := float64(space.x) / float64(level.xcoords)
	howUp := float64(space.y) / float64(level.ycoords)

	orientation := 0
	if AbsoluteValue(howRight) < AbsoluteValue(howUp) {
		if howUp > 0.5 {
			orientation = 180
		} else {
			orientation = 0
		}
	} else {
		if howRight > 0.5 {
			orientation = 270
		} else {
			orientation = 90
		}
	}
	worm := NewWorm(wormName, space, orientation)
	level.worm = worm
}

func (level *Level) GetFreeSpaceForWorm() Coordinate {
	x := rand.Int() % level.xcoords
	if x == 0 {
		x++
	}

	y := rand.Int() % level.ycoords
	if y == 0 {
		y++
	}

	coord := Coordinate{x, y}
	return coord
}

func (level *Level) NewCandy() {
	levelXY := make([][]bool, level.ycoords)
	for i := range levelXY {
		levelXY[i] = make([]bool, level.xcoords)
	}

	wormOccupiedPositions := level.GetWormPositions()

	for _, Coord := range wormOccupiedPositions {
		levelXY[Coord.y][Coord.x] = true
	}
	randIndex := rand.Int() % ((level.xcoords * level.ycoords) - len(wormOccupiedPositions))
	for i := 0; i < randIndex; i++ {
		if levelXY[i/level.xcoords][i%level.xcoords] {
			randIndex++
		}
	}
	level.candy = &Coordinate{x: randIndex % level.xcoords, y: randIndex / level.xcoords}
}

func (level *Level) GetWormPositions() []Coordinate {
	occupiedWormPositions := []Coordinate{}
	if level.worm == nil {
		return occupiedWormPositions
	}
	for wormPiece := level.worm.head; wormPiece != nil; {
		occupiedWormPositions = append(occupiedWormPositions, wormPiece.value)
		wormPiece = wormPiece.next
	}
	return occupiedWormPositions
}

func (level *Level) WormWallCollision() bool {
	if level.worm == nil {
		return false
	}
	if (level.worm.orientation == 0 && level.worm.head.value.y == level.ycoords-1) ||
		(level.worm.orientation == 90 && level.worm.head.value.x == level.xcoords-1) ||
		(level.worm.orientation == 180 && level.worm.head.value.y == 0) ||
		(level.worm.orientation == 270 && level.worm.head.value.x == 0) {
		return true
	}
	return false
}

func (level *Level) MoveWorms() {
	if level.WormWallCollision() {
		level.worm = nil
	} else if level.worm.Move() {
		// kill the worm!
		level.worm = nil
	} else if level.worm.head.value == *level.candy {
		level.worm.toGrow++
		level.NewCandy()
	}
}

func (level *Level) NewOrientation(orientation int) {
	if Abs(orientation-level.worm.orientation) != 180 {
		level.worm.orientation = orientation
	}
}

func (level *Level) Restart() {
	wormName := "MatoMatala"
	if level.worm != nil {
		wormName = level.worm.wormName
	}
	level.worm = nil
	level.AddWorm(wormName)
}
