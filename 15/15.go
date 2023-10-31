package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Point struct {
	X, Y int
}
type Sensor struct {
	Position Point
	Beacon   Point
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

func (s *Sensor) Coverage() int {
	return Abs(s.Position.X-s.Beacon.X) + Abs(s.Position.Y-s.Beacon.Y)
}

func (s *Sensor) MhDistanceTo(target_y int) int {
	return Abs(s.Position.Y - target_y)
}

func (s *Sensor) MhDistanceToXY(x, y int) int {
	return Abs(s.Position.Y-y) + Abs(s.Position.X-x)
}

func Dist(a, b Sensor) int {
	return Abs(a.Position.X-b.Position.X) + Abs(a.Position.Y-b.Position.Y)
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
	// Part A
	// find which x coordinates are occupied at target_y
	// based on the sensors' coverage
	occupied := make(map[int]bool)
	for _, s := range sensors {
		d := s.Coverage()
		d_y := s.MhDistanceTo(target_y)
		if d >= d_y {
			dy := d - d_y
			occupied[s.Position.X] = true
			for i := 1; i <= dy; i++ {
				occupied[s.Position.X-i] = true
				occupied[s.Position.X+i] = true
			}

		}
	}
	// remove the ones that are occupied by beacons
	for _, s := range sensors {
		if s.Beacon.Y == target_y {
			occupied[s.Beacon.X] = false
		}
	}

	// count unique occupied x coordinates
	keys := make([]int, 0, len(occupied))
	for k, v := range occupied {
		if v {
			keys = append(keys, k)
		}
	}
	fmt.Println("A:", len(keys))

	// Part B
	// The beacon has to be between the areas covered by two sensors.
	// Since there is only one solution (assumed in task) the space
	// between the areas can be only "one point wide" (this gives the equality E_1 below).
	// This does not cover the edge case, where the beacon is on the edge of the area,
	// but we'll just ignore that.
	var xBeacon, yBeacon int
	// we want to check each point only once
	possible := make(map[Point]bool)
	for i, first := range sensors {
		for j, second := range sensors {
			if i == j || j < i {
				continue
			}
			d := Dist(first, second)
			coverageFirst := first.Coverage()
			coverageSecond := second.Coverage()
			// E_1: the distance has to be as close as possible but without overlapping
			if d == coverageFirst+coverageSecond+2 {
				// get all points just behind the covered area
				// we are interested in the interesection of these points from both sensors
				fromFirst := ManhattanPoints(Point{first.Position.X, first.Position.Y}, coverageFirst+1)
				fromSecond := ManhattanPoints(Point{second.Position.X, second.Position.Y}, coverageSecond+1)
				for p := range fromFirst {
					if fromSecond[p] {
						possible[p] = true
					}
				}
			}
		}
	}

	fmt.Println("search space size", len(possible))
	// check which of the possible points is not covered by any sensor
	for p, _ := range possible {
		ok := true
		for _, s := range sensors {
			if s.MhDistanceToXY(p.X, p.Y) <= s.Coverage() {
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
