package main

import (
	"container/heap"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"slices"
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

const END_DAY = 26
const OPEN_TIME = 1

type Node struct {
	position     string
	released     int
	currentFlow  int
	openedValves map[string]bool
	day          int
	priority     int // The priority of the item in the queue.
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
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
	// all valves are needed for distance calculations
	valves := makeValves(valvesArr)
	fmt.Println("# Valves:", len(valves))

	// distances: Floyd-Warshall algorithm
	dist := floydWarshall(valves)

	// now we can keep only relevant valves
	for name, v := range valves {
		if v.Flow == 0 {
			delete(valves, name)
		}
	}
	fmt.Println("# Valves (pruned):", len(valves))

	maxFlow := 0
	for _, v := range valves {
		maxFlow += v.Flow
	}

	// A* search
	startingNode := &Node{
		position:     "AA",
		openedValves: make(map[string]bool, 0),
	}
	fmt.Println("Generating best solution for elf")
	bestSolutionNode, openedCount := findBest(startingNode, valves, dist, maxFlow)
	fmt.Println("T:", bestSolutionNode.day, "A:", bestSolutionNode.released)
	fmt.Println("# Opened:", openedCount)

	fmt.Println("Generating best complementary solution")
	startingNode = &Node{
		position:     "AA",
		openedValves: make(map[string]bool, 0),
	}
	for k, v := range bestSolutionNode.openedValves {
		startingNode.openedValves[k] = v
	}
	remainingSolutionNode, openedCount := findBest(startingNode, valves, dist, maxFlow)
	fmt.Println("T:", remainingSolutionNode.day, "A:", remainingSolutionNode.released)
	fmt.Println("# Opened:", openedCount)

	// use remainingSolutionNode.released as heuristic (prune below) for new "exhaustive search"
	fmt.Println("Generating all feasible solutions (ones better than the complementary one)")
	startingNode = &Node{
		position:     "AA",
		openedValves: make(map[string]bool, 0),
	}
	allFeasible, openedCount := findAll(startingNode, valves, dist, maxFlow, remainingSolutionNode.released)
	fmt.Println("# Found:", len(allFeasible))
	fmt.Println("# Opened:", openedCount)

	bestTotalReleased := 0
	bestTotalSolution := &Node{}
	bestTotalSolutionComplementary := &Node{}
	for _, feasible := range allFeasible {
		startingNode = &Node{
			position:     "AA",
			openedValves: make(map[string]bool, 0),
		}
		for k, v := range feasible.openedValves {
			startingNode.openedValves[k] = v
		}
		bestComplementary, _ := findBest(startingNode, valves, dist, maxFlow)

		totalReleased := feasible.released + bestComplementary.released
		if totalReleased > bestTotalReleased {
			bestTotalReleased = totalReleased
			bestTotalSolution = feasible
			bestTotalSolutionComplementary = bestComplementary
		}
	}
	fmt.Println("Best total solution:", bestTotalReleased)
	fmt.Println(bestTotalSolution, bestTotalSolutionComplementary)
}

func findAll(startingNode *Node, valves map[string]Valve, dist map[string]map[string]int, maxFlow int, threshold int) ([]*Node, int) {
	q := make([]*Node, 0)
	q = append(q, startingNode)
	solutions := make([]*Node, 0)
	openedCount := 0

	for len(q) > 0 {
		openedCount++
		currentNode := q[0]
		q = q[1:]

		nbs := generateNeighbors(currentNode, valves, dist, maxFlow)
		for _, node := range nbs {
			// if actual release is above threshold, it is feasible solution
			// -> save it
			if node.released > threshold {
				solutions = append(solutions, node)
			}
			// if priority is above threshold, it is possible this will be feasible solution
			// -> continue processing
			if node.priority > threshold {
				q = append(q, node)
			}
		}
	}
	return solutions, openedCount
}

func findBest(startingNode *Node, valves map[string]Valve, dist map[string]map[string]int, maxFlow int) (*Node, int) {
	q := make(PriorityQueue, 0)
	heap.Init(&q)
	heap.Push(&q, startingNode)

	bestSolution := 0
	bestSolutionNode := &Node{}
	openedCount := 0

	for len(q) > 0 {
		openedCount++
		currentNode := q[0]
		q = q[1:]

		nbs := generateNeighbors(currentNode, valves, dist, maxFlow)
		for _, node := range nbs {
			if node.released > bestSolution {
				bestSolution = node.released
				bestSolutionNode = node
			}
			if node.priority >= bestSolution {
				heap.Push(&q, node)
			}
		}
	}
	return bestSolutionNode, openedCount
}

func generateNeighbors(currentNode *Node, valves map[string]Valve, dist map[string]map[string]int, maxFlow int) []*Node {
	nbs := make([]*Node, 0)
	// Initial action: Travel to a node and open it.
	// Since we're allowing travel to any node we can assume
	// we always want to open it. If we don't want to open it
	// we must want to open another node. But in that case we
	// can travel directly to that node.
	//
	// Generate neighbors for all possible actions
	// 1) move to an unopened valve and open it
	// 2) wait till the end
	//
	// No other actions are needed.

	// 1)
	for valveName := range valves {
		_, alreadyOpened := currentNode.openedValves[valveName]
		if alreadyOpened {
			continue
		}
		travelTime := dist[currentNode.position][valveName]
		totalTime := travelTime + OPEN_TIME
		endTime := currentNode.day + totalTime
		// Make sure we can travel to the node in the given time.
		if endTime < END_DAY {
			releasedAtEndTime := currentNode.released + currentNode.currentFlow*totalTime
			// Optimistic estimate of total flow
			// - actually released at end time
			heuristic := releasedAtEndTime
			// - if everything was opened after that
			heuristic += maxFlow * (END_DAY - endTime)
			n := &Node{
				position:    valveName,
				currentFlow: currentNode.currentFlow + valves[valveName].Flow,
				released:    releasedAtEndTime,
				day:         endTime,
				priority:    heuristic,
			}
			n.openedValves = make(map[string]bool)
			for k, v := range currentNode.openedValves {
				n.openedValves[k] = v
			}
			n.openedValves[valveName] = true
			nbs = append(nbs, n)
		}
	}
	// 2)
	daysTillEnd := END_DAY - currentNode.day
	if daysTillEnd < 0 {
		panic("invalid value")
	}
	if daysTillEnd > 0 {
		totalReleased := currentNode.released + currentNode.currentFlow*daysTillEnd
		n := &Node{
			position:    currentNode.position,
			currentFlow: currentNode.currentFlow,
			released:    totalReleased,
			day:         END_DAY,
			priority:    totalReleased,
		}
		n.openedValves = make(map[string]bool)
		for k, v := range currentNode.openedValves {
			n.openedValves[k] = v
		}
		nbs = append(nbs, n)
	}
	return nbs
}

// convert to map for faster/easier access
func makeValves(valvesArr []InputValve) map[string]Valve {
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
	return valves
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
