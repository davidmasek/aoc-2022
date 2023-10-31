package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type Sensor struct {
	X, Y             int
	BeaconX, BeaconY int // beacon
}

// Load sensors from a file
func load(filename string) []Sensor {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var sensors []Sensor
	err = json.Unmarshal(data, &sensors)
	if err != nil {
		panic(err)
	}

	return sensors
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *Sensor) MhDistance() int {
	return Abs(s.X-s.BeaconX) + Abs(s.Y-s.BeaconY)
}

func (s *Sensor) MhDistanceTo(target_y int) int {
	return Abs(s.Y - target_y)
}

func (s *Sensor) MhDistanceToXY(x, y int) int {
	return Abs(s.Y-y) + Abs(s.X-x)
}

func Dist(a, b Sensor) int {
	return Abs(a.X-b.X) + Abs(a.Y-b.Y)
}

func Sign(x int) int {
	if x == 0 {
		return 0
	}
	if x < 0 {
		return -1
	}
	return 1
}

type Point struct {
	X, Y int
}

func ManhattanPoints(p Point, d int) map[Point]bool {
	points := make(map[Point]bool)
	for o := 0; o < d; o++ {
		inv := d - o
		points[Point{p.X + o, p.Y + inv}] = true
		points[Point{p.X + inv, p.Y - o}] = true
		points[Point{p.X - o, p.Y - inv}] = true
		points[Point{p.X - inv, p.Y + o}] = true
	}
	return points
}

func solve(sensors []Sensor, target_y, limit int) {
	occupied := make(map[int]bool)
	for _, s := range sensors {
		d := s.MhDistance()
		d_y := s.MhDistanceTo(target_y)
		if d >= d_y {
			dy := d - d_y
			occupied[s.X] = true
			for i := 1; i <= dy; i++ {
				occupied[s.X-i] = true
				occupied[s.X+i] = true
			}

		}
	}

	for _, s := range sensors {
		if s.BeaconY == target_y {
			occupied[s.BeaconX] = false
			// fmt.Println(s)
		}
	}

	keys := make([]int, 0, len(occupied))
	for k, v := range occupied {
		if v {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	fmt.Println("A:", len(keys))

	// fmt.Println(keys)
	var xBeacon, yBeacon int
	// bar := progressbar.Default(int64(len(sensors) * len(sensors)))
	// outer:
	possible := make(map[Point]bool)
	for i, first := range sensors {
		for j, second := range sensors {
			if i == j || j < i {
				continue
			}
			d := Dist(first, second)
			coverageFirst := first.MhDistance()
			coverageSecond := second.MhDistance()
			if d == coverageFirst+coverageSecond+2 {
				fromFirst := ManhattanPoints(Point{first.X, first.Y}, coverageFirst+1)
				fromSecond := ManhattanPoints(Point{second.X, second.Y}, coverageSecond+1)
				// common := make(map[Point]bool)
				for p := range fromFirst {
					if fromSecond[p] {
						// common[p] = true
						// possible = append(possible, p)
						possible[p] = true
					}
				}
				// if first.X == 8 && first.Y == 7 && second.X == 20 && second.Y == 14 {
				// 	fmt.Println("possible", common)
				// }
			}
		}
		// bar.Add(1)
	}

	for p, _ := range possible {
		ok := true
		for _, s := range sensors {
			if s.MhDistanceToXY(p.X, p.Y) <= s.MhDistance() {
				ok = false
				break
			}
		}
		if ok {
			fmt.Println("ok", p)
			xBeacon = p.X
			yBeacon = p.Y
		}
	}

	fmt.Println("search space size", len(possible))

	fmt.Println(xBeacon, yBeacon)
	fmt.Println("B:", xBeacon*4000000+yBeacon)
}

func main() {
	fmt.Println("Example")
	fmt.Println("------")
	sensors := load("15.ex.json")
	solve(sensors, 10, 20)
	fmt.Println("Test")
	fmt.Println("------")
	// A: 4961647
	// B: 12274327017867
	sensors = load("15.in.json")
	solve(sensors, 2000000, 4000000)
}
