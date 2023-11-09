package main

import (
	"fmt"
)

type Point struct {
	x, y int
}

// Possible rocks:
// ####

// .#.
// ###
// .#.

// ..#
// ..#
// ###

// #
// #
// #
// #

// ##
// ##
type Rock struct {
	position Point
	blocks   []Point
	w, h     int
}

func RockMinus(position Point) *Rock {
	return &Rock{
		position: position,
		blocks: []Point{
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
		},
		w: 4,
		h: 1,
	}
}

func RockPlus(position Point) *Rock {
	return &Rock{
		position,
		[]Point{
			{1, 0},
			{0, -1},
			{1, -1},
			{2, -1},
			{1, -2},
		},
		3,
		3,
	}
}

func RockL(position Point) *Rock {
	return &Rock{
		position,
		[]Point{
			{2, 0},
			{2, -1},
			{0, -2},
			{1, -2},
			{2, -2},
		},
		3,
		3,
	}
}

func RockTall(position Point) *Rock {
	return &Rock{
		position: position,
		blocks: []Point{
			{0, 0},
			{0, -1},
			{0, -2},
			{0, -3},
		},
		w: 1,
		h: 4,
	}
}

func RockSquare(position Point) *Rock {
	return &Rock{
		position: position,
		blocks: []Point{
			{0, 0},
			{1, 0},
			{0, -1},
			{1, -1},
		},
		w: 2,
		h: 2,
	}
}

func draw(rocks []*Rock) {
	h := 0
	for _, r := range rocks {
		h = max(r.position.y, h)
	}
	h += 1
	w := 7

	occupied := make([][]bool, h)

	for i := 0; i < h; i++ {
		occupied[i] = make([]bool, w)
	}

	for _, r := range rocks {
		for _, p := range r.blocks {
			occupied[r.position.y+p.y][r.position.x+p.x] = true
		}
	}

	for y := h - 1; y >= 0; y-- {
		fmt.Print("|")
		for x := 0; x < w; x++ {
			isRock := occupied[y][x]
			if isRock {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("+-------+")
}

func (rock *Rock) isBlockedAt(shift Point, pile map[Point]bool, chamberWidth int) bool {
	if rock.position.x+shift.x < 0 {
		// fmt.Println("blocked by left wall")
		return true
	}
	if rock.position.x+shift.x+rock.w > chamberWidth {
		// fmt.Println("blocked by right wall")
		return true
	}
	if rock.position.y-rock.h+1+shift.y < 0 {
		// fmt.Println("blocked by ground")
		return true
	}

	for _, p := range rock.blocks {
		shiftedX := rock.position.x + p.x + shift.x
		shiftedY := rock.position.y + p.y + shift.y
		if _, blocked := pile[Point{shiftedX, shiftedY}]; blocked {
			// fmt.Println("blocked by other rocks", shiftedX, shiftedY, pile)
			return true
		}
	}
	return false
}

func solve(winds string) {
	const CHAMBER_WIDTH = 7
	// start with left edge this from left edge of chamber
	const START_DX = 2
	// start with bottom this above highest rock
	const START_DY = 3
	rocks := make([]*Rock, 0)
	rock := RockMinus(Point{START_DX, START_DY})
	rocks = append(rocks, rock)
	pile := make(map[Point]bool)
	// draw(rocks)
	fallenCounter := 0
	for round := 0; ; round++ {
		windIndex := round % len(winds)
		wind := winds[windIndex]
		// fmt.Println("WIND:", string(wind), "Round:", round)
		var shift Point
		if wind == '>' {
			shift.x = 1
		} else if wind == '<' {
			shift.x = -1
		} else {
			panic("unexpected input")
		}
		if !rock.isBlockedAt(shift, pile, CHAMBER_WIDTH) {
			// fmt.Println("can move")
			rock.position.x += shift.x
		}
		// draw(rocks)
		blockedBelow := rock.isBlockedAt(Point{0, -1}, pile, CHAMBER_WIDTH)
		if blockedBelow {
			// fmt.Println("--- NEW ROCK ---")
			for _, p := range rock.blocks {
				pile[Point{rock.position.x + p.x, rock.position.y + p.y}] = true
			}
			startY := 0
			for _, rock := range rocks {
				startY = max(startY, rock.position.y)
			}
			startY += START_DY
			fallenCounter++
			if fallenCounter == 2022 {
				// if fallenCounter == 10 {
				fmt.Println("A:", startY-START_DY+1)
				break
			}
			switch fallenCounter % 5 {
			case 0:
				rock = RockMinus(Point{START_DX, startY})
			case 1:
				rock = RockPlus(Point{START_DX, startY})
			case 2:
				rock = RockL(Point{START_DX, startY})
			case 3:
				rock = RockTall(Point{START_DX, startY})
			case 4:
				rock = RockSquare(Point{START_DX, startY})
			}
			rocks = append(rocks, rock)
			rock.position.y += rock.h
			// draw(rocks)
		} else {
			// fmt.Println("DOWN")
			rock.position.y--
			// draw(rocks)
		}
	}
}

func main() {
	fmt.Println("------")
	fmt.Println("Example")
	fmt.Println("------")
	solve(">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>")
	fmt.Println("------")
	fmt.Println("Test")
	fmt.Println("------")
	solve(">>>><<<>><<><<><<><>>>><>><><<<<>>>><<<>><<<>><<<>>><<<>><<<<>>><<<<>>><<<>><><<<<>>>><<<<>>>><<>>><>>>><<>>><<<<>>><><<<<><>>>><>>>><>>>><><>><>>><<>>>><>>>><<<>>>><<>>><<<<>>><><<<<>>><>>><<<>>>><>><><<<>>><<<>>>><<>><<><<<>>>><><<>>><<><<<><<<><><><<<<>>>><<<<>>><<<<>><<>>><><<<<>><<<<>>><<>>>><<><>><<<<>>>><<<>>>><<<>>>><><<<>>>><<<<><<><<>>>><<><<<<>><<>>><<<><<>><><<<<>><<<><<<>><<<>><<<<><>>><<<<>>><<<>>><<<><<<>>>><<>>><<<>>>><<<>>><><<<<>><<<>>><<<>><<<><<<>><<>>>><<>><<>>>><<><<<>>><<<<>>>><<<>><<>>>><>><>><>><<<<>><><>><<<>><><>>>><<><<<<>><>>>><<<<>>><<>><<<>><<<>>><<<<><<<><<<<>>>><<<<>>>><<<><<<<>><<<>>>><>>>><<<>>>><<<<><><<<<>><>>>><<<>><><<<<><<>>><>>>><>>><<<<><<><<<><>><<<<>>><>>><>>>><<<><<<><>><<<>>>><<<<>>><<><<>>>><><>>><><<>>>><<<>>>><<<>><><<<<><<<>>><<>>><<<>><<<>><<>>>><<>>><<<><<<>><<<<><>><><<<>>>><<<<>>>><<<<>><<<<>>>><<<<>><<>>>><<<<><<>>><<>>><<<<>>><<<>>><>>><<<>>>><<<<><<>>>><<<<>>>><<<<>>>><>><<>>>><><<<>>><>>><<<>><<<<>>><>>>><<>>>><<<>>>><<<<>><<<>>>><>>>><<<<>>>><<<<>><<>>><<><>>><<<<><<<<><<><<<<>>><>>>><<<<>>><><>><<<<><<<<>><<>>>><>>><<>>><<>>>><<>><><>>>><>><<<<>>><<>>>><<<<>>><<<>>>><<<>>><<<<><<><<<>><><<<<>>><<<><<<>><<<<>>><>>><>>><<<<>>>><<<>>>><<>>><<<<>>><<<<><<<<>>>><<>><<<<>>><<<<>>>><>>>><<<<>><<<<>><<>>><><<><<><<<<>><<>>>><<>><<<>>>><>>>><>>>><<<><<<>>>><<><<>>><<<>><<<<><<<<>><<<<>><<<<>>><<>><>><<>><<<>>>><<<>><>>>><<<>>>><<>><<>>>><>><<<><<<>>>><<<><<>>><<<>>>><<><<<>>>><<<<><<<>>>><<<>>><<>><<><<>>>><<><<<>><<<<>><<>><>>>><>><<<<>>><>>>><<>><<>>><<<<>>>><<<<>><>><<<>>>><><>>><<><<<<><>><<<>>><<<>><><><<<>><<<>>>><<<<>>><>>>><<<><>><<<>><<>>><<<><<<<>><<<>>><><<>><>>><<<>>>><>>><<>><<<>>>><<<>>><><<<<>>><>>>><<<<>>><<<<>>><<<><><><<<>>><><<>>>><<<>>>><>>>><><<<>>>><<<>>>><>>><<>>>><<<>>>><<>><<<<>>><<<>>>><<<><>><<<<>>><>><<<><<<>>>><>>><>>><<>><<<>><<<<><<<>><>>><>><<<>>><<<>><<<<>>><<<><<>>><<<<>>><<<><>>>><<><<>><>>><<>><<>>><<<><<<><<<><<<>><<>>>><>>>><>><<<<>>><<<<><<<<>><>>>><<<>>>><<<>><<<>>><<<<>>><><>><<<<><<><>>>><<<>><><<><<>>>><<>><><<>><<<>><<<<>>>><><<<<>><<<><<>><><<>><>>><>>><<<<><<>><<>><<><<>>>><><<<><<<>>>><<<>>>><<><<<<><>><<<<>>><<<<>><>>>><>><>>><<<>>><<<>>><<<<>><>><<<<>>><<<>>>><>>>><<<><<<>><<>>><<>>><<>><>><<<<>>>><>>><>><<>>>><<<<>><>><<<<>>>><<>><<<<>><>>><<><<<>><<<>>><<<<>>><<><<>>>><<<>>><<<<>>><<<<>><><<>><<>><<<<>>>><<<><<>>>><<<><><>>><><<<>>><><<<<><<<<><<<>><<<>>>><>>><<<<>><<<<><<<>>><<><<<>>>><<<>>><<<<>><<<<>>><<<<><>>>><<<<>><<<><<<<>>>><>><<<<>><<<<>><<<><<<>><<<>>>><>>><<<>>><<<<>><<>>><<<>>>><>>>><<>>><<<><<>>><<<<>>>><<>>>><<<>><<<<>>>><><<<>>><<<>>>><<<>>><<<><<><<<<><><<<<>>>><<<>>>><><><<<<>><<>>>><<<<>>>><>>>><<<<>>>><><<><<>><<><<>><<<<>>><<>>><<<<>>>><<<<>>>><<<><>><<<<>>>><<<><>>>><>>><<><><><<<>>><<>>>><<>>><<>><<<<>>>><<<>>>><<>>><<><>>>><<<>><>>>><<>><>>><>>><<>><><<>>>><<>>>><>><<<<>><<<<>>>><>>>><>>>><>>><>>>><<<>><<>>>><<<>>>><><<<<><<>>><<<>>>><<>>><<><><<<<>>><<<<>>>><<<<>>><<>><<<>><<>><<<<>><<<>>>><<<>>>><<<<>>>><<>>>><<<<><<<>><<<>>><<>>>><>>>><<<<>><<<>><<>><<<>><<>>><<<<><<<<>>><<>><<<<><<<<>>>><><<<<>>><<<><>>><<>><<<<>>><<<<>>>><>><<<<>>><<<<>>>><<<>>>><>>><<><><>>><<<<>>><<<<>>><<<<>><<<<>>><<<<>>>><<<<><<<>>>><>>>><<<<>>><<<>>>><>><<<<><<>>>><>>><<<<><<<<>>>><><<<><<<>>>><<<><>>>><<<<><<<><>>>><<<>><<<>><<<><<<<>><<<<>><<><<>>>><<>>>><<<<>>><>><<<<><<<><<>>>><<<>><>>>><<<>>>><<>><>>><><<><<<>><<>><<<>><<<<><<<<>>><>><<<>><><<<<><<>>><>><<<<><<<>>>><<<><<<<><<>>>><><<<>><<>><<<><<<<>>>><<>><<>><<><<>><>><<<>><<>><<<>><>>><<<<><<<>>>><<<<>><>><<><><>>><>>>><<<<>>>><<<>>>><<<>>><<<<><<<>>>><<><<<<><<<><><<>>><<>>>><<>><<<<>>><<<>>>><<<><<<>>>><<<>>><<<>><<<>><<<<>><<<><>>>><>>>><<<<>>>><<>>><<><>><<<<>>>><<<>>>><>><>>><<><<<><<<><><<<>>><<<<><<<<>>>><<<>>>><<<<>>>><<<><<<<><<<<>><<<>>>><<>>>><>>><<>><<<<>>>><<>><>><<>><<<>><<<<>>><<<<><<<<><<>><<<>><<><<<<>><<><<<>><<<<><<<>>><<<>>><>>>><<<><<<<>><>>><<<>><<>>><<<<>><><<>>><<><<>>><<>>><<<<>><<<<>>>><<<>>>><>>>><<>>>><<<>>>><>>><<<<>><>>>><<<<>>>><>><<<>><>>><<<<><>><><<<>><><<>><<<<>>><<<<><<>>>><<>>><<<<>>>><<<<><<<<><<<>><>>>><<<>>>><<>>>><<>><><>>><<<>>><<>>>><><<<<><<>>><<>>>><<<>>><<<<>><>>><<<<>><>>><><>>><<>>><><>>>><<<>>>><<<>><<<<>>>><<<<>>>><<>><<<<><<<<>>>><<<<>>><<<>>><<<<><<<<>><<<<>>>><<<>>><<>><<><>><<>>><>>><<>>><<<<>>><<<>>>><<<><><<>><<<><<<>>><><>><<<<><<>><<<<>><<<<><<<<>><>>>><<<>>>><<<<><<<><<<<>>><<>>>><>>><<<><<>>>><>>>><>>><<<<>>>><<<>><>>><<><<<<>><<>><<>>><<<<><<>>>><<<>><<<<><<>>><>>><<>><><<<<>><<<><><<<>>>><<<>><>><<>><<<>>><<>>><<<>>>><<<>>><><<>><>>>><>>><>><<<<><><<<>>><>>>><><<<<><<<><>>><<<>><<<<>>>><<<<>>><<>><<<<>><<>><<>>><<<>>>><<>>>><<><<<>>><<<<>><<<>><<<>>><<<<>>>><><<<<>><<>><<<<>><<><<<<>>>><<>>><<<><>>>><<>>>><<>><<<><<>><<<>>>><<<><<><<<<>><<<<>>>><<<>>>><<<<>>>><<<>><><<<>>>><<<<>>><<>>><<<>><><<<<>>><>>>><<<<>>>><<<<>><<<><><<<<><<<>>>><<<<>>>><><<<>><>>>><<<<>>>><<<<><<<><<<<>>><>>>><>>>><<<<>><<>><<<>>>><<<<>>>><<<<><<<<>>>><<>>><<>>><<<>>>><<>>>><<<>>>><<>><<>>>><<<<>>>><<<<><<<><>><<>><<<<><>>>><>><<<>>>><<<><<>>>><<><><><>>>><<<<>>><<><<<>><<<<>><<<>><<<>>><<<<>>>><<<>>><<<<>><<<>>>><<<>><<<>>>><<<<>>>><<>>>><<<<>><<<>>><<>>><<<<>>>><<<<>><<<<>><<<<>>>><<<><<<>><>>><>>>><<>><>>>><<>><<<>>>><>>>><>><>>><<<<>>><<<<>>><<<<>>><<<>>>><<<>>>><<<<>>><<<>>><<>>>><>>><<>>><>><><<<>>><><<<>>><<<>><<>>><<<<>>><>>><<><<>>><<<<>>>><<>>><<<>><<>><<<>>><<<<>>>><<<>>><<<<>>><<<<>>><<<><<>>>><<<<>>><<<>>>><<>>><<>><>>>><<<>>>><<<<>>>><>>>><>>>><<><<><<<>><<<>><<><<<<>>><<<<>><>>>><<>>><<>>>><<<>><<<><<<>>>><<<>>><<<<>>><<<>>><<<<><<><<<<>><<<<><<><<<>><<<<>><><<<<>>>><<<<>><<>><<>>>><><<>><>>><<>><<<><><<>><<<<>><>><<<<>>><<><<>><>>>><>>>><<<>>><<<><<>><<><<<>>>><<<><<>>>><>>>><<<>>>><<>><<<>><>><><<<>>>><<>><<<<>><<<<>>><<<<>>>><<<<>><<<<>>><>>><>>><<><>>><<>><<<<>>>><><<<>><<<<>>>><>><<<>><<><><<><<<<>><<>>>><>>>><<>>><<>><<><><<>>><>><<<<>>><<<><<<<><><<>>>><<>>>><<<<>>><<<>>><<<>>>><<<>>><<<<>>>><<>>>><<>><>><<<>>><<<<><<<>>><<>>><<<<>><><>>>><<<>><<>>><<<><<><<<<>><<<<>>>><><<<>><<<<><<><<><<<>>>><>>><>><<<>><<<>>><><<<><<<<><>>>><><<>><<<>><<>>><<<>>>><<<<><>>>><<<>>><><<<><<<>><<>><<>><<<<><<<<>>>><><<<>>><<>><<<>>>><<>>><<<<><<<<>><<><<<>><<<<>>><<<<><<>>>><<<<>><<<>>>><><<>><>><<<><<>>>><<>>>><<<>>>><<<<>>>><<>>><><<<<>>>><<>>><><<<>><>>><<<<>>>><<<<><<<>>><<>>>><><<>><<>>><<<>>>><>><><<<<>><<<>>>><<>><<<>><<><<<<>><><<><<<>>>><<><<<<>>><<><<<>>>><>><<<>>>><<>>><><<<<>><<<<>>>><<>>>><><>>><<<><<>>>><<>>><<<>>><<<><>>>><<>><<<>><><<<>>>><>><>>><<<<>><<>>>><><<<>>>><<>>>><<<<><>>>><<<>>>><<<>><<<>>><<<<><<>>>><<<<>>><<>>><><<<<>>><>>>><<><<<>>>><<><<><<<<><<<><<>>>><>>><><<<<>>>><<<<>>>><><<<<>>><<>>>><<<<>>>><<<>>>><<>>><<<<>><<>>>><<><<><<>><<<<>>><<><<<<><>>>><<>><<>>>><<><>><<<>><<>>>><<<>>><<<<><<<<>>>><><<<<><<<><<>>><<<><<<<>>>><>>>><<<><<>><>>>><<<<><>>><<<>>><<><<<<>>>><<<<>><<<<>>>><<<<>><<<<><>><<<>>><>><<<<>>>><<<<><<<>>><<><<<<>>>><<>>>><<<<>>>><<>>>><<<<>>>><<<>>>><><>><<<>>>><<<<>>><<<<>>>><<<>>><<<<>>>><<><<>><><<>>>><>>>><<<<>>>><<>>>><>>><<<>><>>>><<<<>>><>>>><<<<>>>><<><>><<><<<<>>><<<<>><<<<>>>><>>>><<<>>><<>>>><<<<>>>><<>>><<<<><<>>>><>>><>><<<<>><<<<>>>><<<<>>><<><<>>><<><><<<>>>><<<>>><<<>><<<<><<<<>>><<>><<<<>>>><<<<>>>><<<><<<><><<>>><<<<><<<>><<<>>><<<<><><<>><<<>>>><<<<><<>>>><<<<><<>>><<>>>><<<<><>><<<<><<<<><<>>>><<><<<>><>><<<>><>>><<<>><>>>><<<>><<>>>><<><<>>><<>><>>><<<<>>><>><<<><<><><><<>>><<<><>><<>><<>><<>>><<<<>>>><<<>>><<<>>>><<<<>>><>><<>>><>>>><<<>><<<>><<<<>>><<<<><<>>>><<<<><>>>><<<<>>><<<>>>><<<>>>><<<<>><<>>>><>>><><<<<>>><<<<>><>>><<<<><<<<>>><><<<>>>><<<<>><<<>>>><<<<>>><<>><<>><<<<>><<>><<<<>><<>>>><>>>><<<<>>>><<<<>>><><>>><<>>><<<<>>>><>>>><>><<<<>>>><<><<<><<>>>><<<><<<><<>><<<<>>>><<<><<>><<<>><<>><<>>><<><<<>>><<<<>>><<<>>><>>>><<<<><<>>><>>><<<>>><>>>><<<<>>>><<<>>><<<<>><<<><<<<><<<>>><><<<<>>><<<>>>><<>>>><<>><<<>>>><<<<>>><<<><<<>><<<<>>><<<>>>><<<>>><<<>><<<>>>><<>><<<<>><<><><<>>><>>>><<<<>>>><<>>>><<<>>><<<>><<<<>>><<<<><<<<>>><<<<>>><><<>>>><<<<>>>><><<<<>><<>><<><<><<<>>><<<<><<>>>><<><<><>><<<<>>>><>><>><<>>>><<<>><<<>>>><<<>><>><<>><<<>>><<<<><<<>>><<<>>><<>>>><>><<<<>><<<>><><<<<>>>><>>>><>>>><<<>>><>>>><<<>>>><<<>>>><>><<<<>><><>>><><<>>>><<<>>>><<>>><<<>><<<<>>><<>>><>>><<<<>>>><><<<<>>><<>>><><<<>><<<<><<<><>>>><<>><<>><<><<<<><<<><<>>><<<>>>><<>><<<<>>><<<<>><><<<<>>>><<<<>>><><<<>><<>>>><>><>><<<>>>><<>><>>><<>>><<><<>>>><>>><><<>>><<>><<>>>><<<><<<<><>>><<<<>>>><<<>>><<<<><<<>>><<<<><<>><<<<><<<<>>>><>>><>><>>><>><>><><<<>>><><>>>><<<>>>><>><<>>>><<<>><<<>>><<>>>><>>>><>>><<<<>>><<>>><<<>><<<<>><<<>>>><<<>><<<<>><>>><<>>><>>><<<>>>><<>><<<<><<<><>>><><<<>>><<>><><<<>>>><>>>><<<<>>>><<<><<<><>>><<<>>>><<<<><<<<><<>>><<>><<<>><>>><<>>><<>>>><<<<><<<<>><<<>><<>>>><><<<<><<<<>>><<<>>>><<<<>><<>>>><>>>><>><<<>>><<<>><<<><<<>>>><><<<<>><<>>>><>>><<<>>>><<><<<<>><<<>>>><<<>>><<<>><<<>>><><<<<>><<<<><<<<>>>><<<>>>><>>><<><>>><<<>><>><>>><<>>>><<<>>><<>>>><>>>><<<><<<>>><<<<><<<>><<<>><<<><<>>><<<>><<<<>>><<<<>><<>><<<<>>>><<>><<<>><>>><<<<>><>>>><<<<><<<<>><>>>><<<<><<><<<><<>><><<<>>>><<<<>>><<><>>>><<>><<<<>><<<<>><<<><<>>><<>>>><<<<>>>><<><<<<>>>><<><<<><<<><<<<>>><<<>>><<<>><>><><<<<>>>><>><<>>>><<><><<>>><<<<>><<<><<<>>>><<<><<<<>><<<<>>>><<<>>><<<>>><<>>><><<<>><<<<>>><<>>><<<<>>><<>>><>>>><<<>>><<<<>>><<<>><<<<>><<<>><<><<<<>>>><<>><<<>>>><<<>>><<<<>>><<<<>>><<<<><<<<>>>><<<<>>><<><<<><<<<>>><<<<>>>><<<>>><<<<>>>><<><<><<<<>><<<<>>><>>>><>>><<>>>><>><<>>>><<<>><<<>>>><><<<>>>><>>><<<<>><><<<>>><<>>><<<<>><>><<>>>><<<<>><<>><>>><<<<>>>><<<>><>><<<<>>>><<<>><<>>><<<>><<<>><<><>>>><><<<>>><<<<>>><<<<><<<<>>><>>><<<<><<<>>><<>><<<<>>>><<<<>>><<<<><><<<<>>><<>><>>><<<>><<><<<<>>><<><<<><>><<<>><<>>><<<>><<<<><<<><<<>>>><<<<>><>><<<<>><<><<<>><<<>>>><<<<>>><>>><<>>><>>>><<<<>>><<<<>>><<>><<>>>><>><>>>><<>><<<>><<<>>>><<<<>>><<>>>><><>><><<<><<>>><<>>><<<<><>>><<<<>><<<<>>><>><<<<>><<<<>>><>>>><><<<<><>>><>>><<<>>><>>><<<<><>><<<><<<<>><<>><><<>>>><<>><<<<><>><<>>>><<>><<>>><<>>>><<<<><>>>><<<>>><<><<>><<<>>><<<<>>><<>><><<<<>>>><<<<><>>>><<<<><<<<><<<>><<<>>>><<<>>>><<<>>>><<<>><<<<>>>><<<<><<<<>>>><<>><<>><<<>><<<<>><<<>><")
}
