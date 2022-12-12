package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

var infinity int = math.MaxInt

type node struct {
	x, y             int
	height, distance int
	visited          bool
}

type heightMap struct {
	start, end *node
	elevations [][]*node
	queue      []*node
}

func (h *heightMap) reset() {
	h.queue = make([]*node, 0, len(h.elevations)*len(h.elevations[0]))
	for _, row := range h.elevations {
		for _, node := range row {
			node.distance = infinity
			node.visited = false
		}
	}
}

func (h *heightMap) findClosestInQueue() *node {
	closest := h.queue[0]
	closestIndex := 0
	for i, n := range h.queue {
		if n.distance < closest.distance {
			closest = n
			closestIndex = i
		}
	}
	h.queue = append(h.queue[:closestIndex], h.queue[closestIndex+1:]...)

	return closest
}

func (h heightMap) findUnvisitedNeighbours(
	middleNode *node) (neighbours []*node) {

	x, y := middleNode.x, middleNode.y
	if x > 0 {
		n := h.elevations[y][x-1]
		if !n.visited {
			neighbours = append(neighbours, n)
		}
	}
	if x < len(h.elevations[0])-1 {
		n := h.elevations[y][x+1]
		if !n.visited {
			neighbours = append(neighbours, n)
		}
	}
	if y > 0 {
		n := h.elevations[y-1][x]
		if !n.visited {
			neighbours = append(neighbours, n)
		}
	}
	if y < len(h.elevations)-1 {
		n := h.elevations[y+1][x]
		if !n.visited {
			neighbours = append(neighbours, n)
		}
	}
	return
}

func (h heightMap) dijkstra(findEnd, findClosestGround bool) int {
	h.reset()
	for _, row := range h.elevations {
		for _, node := range row {
			if findEnd && node == h.start ||
				findClosestGround && node == h.end {
				node.distance = 0
			}
			h.queue = append(h.queue, node)
		}
	}

	for len(h.queue) > 0 {
		node := h.findClosestInQueue()
		if node.distance == infinity {
			break
		}
		if findEnd && node == h.end || findClosestGround && node.height == 0 {
			return node.distance
		}

		for _, neighbour := range h.findUnvisitedNeighbours(node) {
			if findEnd && (neighbour.height-node.height) > 1 ||
				findClosestGround && (node.height-neighbour.height) > 1 {
				continue
			}

			newDist := node.distance + 1
			if newDist < neighbour.distance {
				neighbour.distance = newDist
			}
		}

		node.visited = true
	}

	return infinity
}

func (h heightMap) findShortestPathLength() int {
	return h.dijkstra(true, false)
}

func (h heightMap) findClosestStartingDistance() int {
	return h.dijkstra(false, true)
}

func parseHeightMap(scanner *bufio.Scanner) *heightMap {
	scanner.Split(bufio.ScanLines)
	heights := heightMap{
		elevations: make([][]*node, 0),
	}

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]*node, 0, len(line))

		for j, e := range line {
			n := &node{
				x:        j,
				y:        i,
				distance: infinity,
				visited:  false,
			}

			if e == 'S' {
				heights.start = n
				e = 'a'
			} else if e == 'E' {
				heights.end = n
				e = 'z'
			}

			intHeight := int(e - 97)
			n.height = intHeight

			row = append(row, n)
		}

		heights.elevations = append(heights.elevations, row)
		i++
	}

	return &heights
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	heights := parseHeightMap(scanner)
	shortestPath := heights.findShortestPathLength()
	fmt.Println("The shortest path from S to E is:", shortestPath)

	closestDistance := heights.findClosestStartingDistance()
	fmt.Println("The closest distance to E at height a is:", closestDistance)
}
