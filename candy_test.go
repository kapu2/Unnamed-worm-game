package main

import (
	"testing"
)

func TestNewCandy(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := NewLevel()
		l.AddWorm("testimato")
		l.NewCandy()

		fail := false
		failPos := Coordinate{}
		positions := l.GetWormPositions()
		for _, pos := range positions {
			if pos == *l.candy {
				fail = true
				failPos = pos
			}
		}
		if fail {
			t.Errorf("Candy in same position as worm. Position x:%d y:%d", failPos.x, failPos.y)
		}
	}
}
