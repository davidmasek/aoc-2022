package main

import (
	"container/heap"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"slices"
	"strings"
)

type InputValve struct {
	Name      string   `json:"name"`
	Flow      int      `json:"flow"`
	Neighbors []string `json:"neighbors"`
}

type Valve struct {
	Name             string
	Flow             int
	NeighborDistance map[string]int
}

const N_DAYS = 30

type Node struct {
	position     string
	released     int
	currentFlow  int
	openedValves map[string]bool
	day          int
	actions      []string
	priority     int // The priority of the item in the queue.
	index        int // The index of the item in the heap.
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Compute distances between all pairs using Floyd-Warshall algorithm
func floydWarshall(valves map[string]Valve) map[string]map[string]int {
	dist := make(map[string]map[string]int)

	// Initialize distances to infinity
	for name := range valves {
		dist[name] = make(map[string]int)
		for other := range valves {
			dist[name][other] = math.MaxInt64
		}
	}

	// Distance to self is 0
	for name := range valves {
		dist[name][name] = 0
	}

	// Fill in the distances for the edges that exist
	for name, v := range valves {
		for nbName, nbDistance := range v.NeighborDistance {
			dist[name][nbName] = nbDistance
		}
	}

	// Floyd-Warshall algorithm
	for k := range valves {
		for i := range valves {
			for j := range valves {
				if dist[i][k] != math.MaxInt64 && dist[k][j] != math.MaxInt64 && dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	return dist
}

func solve(valvesArr []InputValve) {
	// convert to map for faster/easier access
	valves := make(map[string]Valve)
	for _, v := range valvesArr {
		valves[v.Name] = Valve{
			Name:             v.Name,
			Flow:             v.Flow,
			NeighborDistance: make(map[string]int),
		}
		for _, nb := range v.Neighbors {
			valves[v.Name].NeighborDistance[nb] = 1
		}
	}
	valvesArr = nil
	fmt.Println("# Valves (initial):", len(valves))
	// prune nodes with 0 flow
	// for _, v := range valves {
	// 	// we'll keep AA as en entrypoint and remove/skip it later
	// 	if v.Name == "AA" {
	// 		continue
	// 	}
	// 	if v.Flow == 0 {
	// 		// for my neighbor (target)
	// 		for targetName := range v.NeighborDistance {
	// 			target, ok := valves[targetName]
	// 			if !ok {
	// 				panic("target not found")
	// 			}
	// 			// add my neighbors
	// 			for nbName, nbDistance := range v.NeighborDistance {
	// 				// don't add target to itself
	// 				if nbName == targetName {
	// 					continue
	// 				}
	// 				// +1 for traveling through this node
	// 				target.NeighborDistance[nbName] = nbDistance + 1
	// 			}
	// 			// remove me from target neighbors
	// 			delete(target.NeighborDistance, v.Name)
	// 			valves[targetName] = target
	// 		}
	// 		// remove me from relevant nodes
	// 		delete(valves, v.Name)
	// 	}
	// }
	fmt.Println("# Valves:", len(valves))

	// Floyd-Warshall algorithm
	dist := floydWarshall(valves)

	const END_DAY = N_DAYS
	const OPEN_TIME = 1

	maxFlow := 0
	for _, v := range valves {
		maxFlow += v.Flow
	}

	// BFS
	// q := make([]Node, 0)
	// A* search
	q := make(PriorityQueue, 0)
	heap.Init(&q)
	// // note that this is not perfect since same valves opened in different order
	// // will be considered different
	// // we can probably just ignore seen for now, since it's not possible to return
	// // seen := make(map[Node]bool)
	for valveName := range valves {
		// Never move to start - this is never needed.
		// Assumess start flow is 0.
		if valveName == "AA" {
			continue
		}
		// Traveling to a node with 0 flow is pointless.
		if valves[valveName].Flow == 0 {
			continue
		}
		// Initial action: Travel to a node and open it.
		// Since we're allowing travel to any node we can assume
		// we always want to open it. If we don't want to open it
		// we must want to open another node. But in that case we
		// can travel directly to that node.
		//
		// Make sure we can travel to the node in the given time.
		travelTime := dist["AA"][valveName]
		totalTime := travelTime + OPEN_TIME
		if totalTime < END_DAY {
			// q = append(q, Node{
			// Heuristic == value that will definitely flow + what would happen if we opened all the valves
			// Should always be optimistic
			heuristic := maxFlow * (END_DAY - totalTime)
			heap.Push(&q, &Node{
				position:    valveName,
				currentFlow: valves[valveName].Flow,
				released:    0,
				day:         totalTime,
				openedValves: map[string]bool{
					valveName: true,
				},
				actions:  []string{fmt.Sprintf("[%d] open %s", totalTime, valveName)},
				priority: heuristic,
			})
		}
	}

	bestSolution := 0
	bestSolutionNode := &Node{}
	openedCount := 0

	for len(q) > 0 {
		openedCount++
		currentNode := q[0]
		q = q[1:]
		// Generate neighbors for all possible actions
		// 1) move to an unopened valve and open it
		// - if possible
		// 2) wait till the end
		// - otherwise
		//
		// No other actions are needed. See explanation for initial action.
		movePossible := false
		// 1)
		for valveName := range valves {
			isStart := valveName == "AA"
			isSamePlace := valveName == currentNode.position
			_, alreadyOpened := currentNode.openedValves[valveName]
			hasZeroFlow := valves[valveName].Flow == 0
			if isStart || isSamePlace || alreadyOpened || hasZeroFlow {
				continue
			}
			travelTime := dist[currentNode.position][valveName]
			totalTime := travelTime + OPEN_TIME
			endTime := currentNode.day + totalTime
			actions := make([]string, len(currentNode.actions))
			copy(actions, currentNode.actions)
			actions = append(actions, fmt.Sprintf("[%d] open %s", endTime, valveName))
			if endTime < END_DAY {
				// See above for heuristic explanation
				heuristic := currentNode.released + currentNode.currentFlow*totalTime
				heuristic += maxFlow * (END_DAY - endTime)
				n := &Node{
					position:    valveName,
					currentFlow: currentNode.currentFlow + valves[valveName].Flow,
					released:    currentNode.released + currentNode.currentFlow*totalTime,
					day:         endTime,
					actions:     actions,
					priority:    heuristic,
				}
				n.openedValves = make(map[string]bool)
				for k, v := range currentNode.openedValves {
					n.openedValves[k] = v
				}
				n.openedValves[valveName] = true
				if heuristic >= bestSolution {
					// q = append(q, n)
					heap.Push(&q, n)
				}
				movePossible = true
			}
		}
		// 2)
		if !movePossible {
			daysTillEnd := END_DAY - currentNode.day
			if daysTillEnd < 0 {
				panic("invalid value")
			}
			totalReleased := currentNode.released + currentNode.currentFlow*daysTillEnd
			if totalReleased > bestSolution {
				actions := make([]string, len(currentNode.actions))
				copy(actions, currentNode.actions)
				actions = append(actions, fmt.Sprintf("[%d] wait %d", currentNode.day+daysTillEnd, daysTillEnd))
				currentNode.actions = actions
				bestSolution = totalReleased
				if bestSolution == 1651 || bestSolution == 1701 {
					fmt.Println("Found best solution after", openedCount, "opened nodes")
				}
				bestSolutionNode = currentNode
			}
		}
	}
	fmt.Println("T:", END_DAY, "A:", bestSolution)
	fmt.Println(strings.Join(bestSolutionNode.actions, ", "))
	fmt.Println(bestSolutionNode.priority)
	fmt.Println("# Opened:", openedCount)
}

func main() {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	idx := slices.IndexFunc[[]fs.DirEntry](entries, func(e fs.DirEntry) bool { return e.Name() == "ex.json" })
	if idx == -1 {
		log.Fatal("Could not find example file ex.json")
	}
	fmt.Println("------")
	fmt.Println("Example")
	fmt.Println("------")
	valves := Load[InputValve]("ex.json")
	solve(valves)
	fmt.Println("------")
	fmt.Println("Test")
	fmt.Println("------")
	idx = slices.IndexFunc[[]fs.DirEntry](entries, func(e fs.DirEntry) bool { return e.Name() == "in.json" })
	if idx == -1 {
		log.Fatal("Could not find input file in.json")
	}
	valves = Load[InputValve]("in.json")
	solve(valves)
}
