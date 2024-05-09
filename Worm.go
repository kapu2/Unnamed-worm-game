package main

type LinkedList struct {
	value    Coordinate
	next     *LinkedList
	previous *LinkedList
}

type Coordinate struct {
	x int
	y int
}
type Worm struct {
	wormName    string
	orientation int
	//length      int
	head   *LinkedList
	tail   *LinkedList
	toGrow int
}

func NewWorm(wormName string, space Coordinate, orientation int) *Worm {
	worm := &Worm{}
	worm.wormName = wormName
	worm.orientation = orientation
	worm.toGrow = 10
	worm.head = &LinkedList{space, nil, nil}
	worm.tail = worm.head
	return worm
}
func (worm *Worm) Turn(orientation int) {
	worm.orientation = orientation
}

func (worm *Worm) SelfCollision() bool {
	for next := worm.head.next; next != nil; next = next.next {
		if next.value == worm.head.value {
			return true
		}
	}
	return false
}

func (worm *Worm) Move() bool {
	newCoord := worm.head.value
	switch worm.orientation {
	case 0:
		newCoord.y++
	case 90:
		newCoord.x++
	case 180:
		newCoord.y--
	case 270:
		newCoord.x--
	}
	newHead := LinkedList{newCoord, worm.head, nil}
	worm.head.previous = &newHead
	worm.head = &newHead

	// growing done, let's destroy tail
	if worm.toGrow > 0 {
		worm.toGrow--
	} else {
		newTail := worm.tail.previous
		newTail.next = nil
		worm.tail = newTail
	}
	return worm.SelfCollision()
}
