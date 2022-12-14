package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	grainStartX = 500
	grainStartY = 0
)

type coord struct {
	x, y int
}

func parseCoord(s string) coord {
	split := strings.Split(s, ",")
	x, _ := strconv.Atoi(split[0])
	y, _ := strconv.Atoi(split[1])
	return coord{x: x, y: y}
}

type caveEntity interface{}

type rock struct {
	x, y int
}

type grain struct {
	x, y      int
	isResting bool
}

func (g *grain) rest() {
	g.isResting = true
}

type sandCave struct {
	entities   map[coord]caveEntity
	overflowed bool
	full       bool
	hasFloor   bool
	floorLevel int
}

func (c sandCave) String() string {
	var minCoord coord
	var maxCoord coord
	for key := range c.entities {
		minCoord, maxCoord = key, key
		break
	}

	if c.hasFloor {
		maxCoord.y++
	}

	for coord := range c.entities {
		if coord.x < minCoord.x {
			minCoord.x = coord.x
		}
		if coord.y < minCoord.y {
			minCoord.y = coord.y
		}
		if coord.x > maxCoord.x {
			maxCoord.x = coord.x
		}
		if coord.y > maxCoord.y {
			maxCoord.y = coord.y
		}
	}

	grid := ""
	for i := minCoord.y; i <= maxCoord.y; i++ {
		line := ""
		for j := minCoord.x; j <= maxCoord.x; j++ {
			v, ok := c.entities[coord{x: j, y: i}]
			if !ok {
				line += "."
			}
			switch v.(type) {
			case rock:
				line += "#"
			case grain:
				line += "o"
			}
		}
		grid += line + "\n"
	}

	if c.hasFloor {
		line := ""
		for i := minCoord.x; i <= maxCoord.x; i++ {
			line += "#"
		}
		grid += line + "\n"
	}

	return grid
}

func (c sandCave) getRestingSandAmount() (amount int) {
	for _, e := range c.entities {
		if g, ok := e.(grain); ok && g.isResting {
			amount++
		}
	}
	return
}

func (c *sandCave) insertFloor() {
	c.hasFloor = true
}

func (c *sandCave) insertWall(start, end coord) {
	if start.x > end.x {
		start.x, end.x = end.x, start.x
	}
	if start.y > end.y {
		start.y, end.y = end.y, start.y
	}

	for x := start.x; x <= end.x; x++ {
		for y := start.y; y <= end.y; y++ {
			if y > (c.floorLevel - 2) {
				c.floorLevel = y + 2
			}

			pos := coord{x: x, y: y}
			c.entities[pos] = rock{x: x, y: y}
		}
	}
}

func (c *sandCave) pourSand() {
	for {
		g := grain{x: grainStartX, y: grainStartY}
		for !g.isResting {
			c.moveGrain(&g)
		}
		if c.overflowed || c.full {
			break
		}

		c.entities[coord{x: g.x, y: g.y}] = g
	}
}

func (c sandCave) grainCanMoveDown(g grain) bool {
	_, ok := c.entities[coord{x: g.x, y: g.y + 1}]
	return !ok
}

func (c sandCave) grainCanMoveLeft(g grain) bool {
	_, ok := c.entities[coord{x: g.x - 1, y: g.y + 1}]
	return !ok
}

func (c sandCave) grainCanMoveRight(g grain) bool {
	_, ok := c.entities[coord{x: g.x + 1, y: g.y + 1}]
	return !ok
}

func (c sandCave) grainInAbyss(g grain) bool {
	for pos := range c.entities {
		if pos.y > g.y {
			return false
		}
	}
	return true
}

func (c sandCave) grainIsBlocked(g grain) bool {
	_, ok := c.entities[coord{x: g.x, y: g.y}]
	return ok
}

func (c *sandCave) moveGrain(g *grain) {
	if c.grainIsBlocked(*g) {
		g.rest()
		c.full = true
		return
	}
	if !c.hasFloor && c.grainInAbyss(*g) {
		g.rest()
		c.overflowed = true
		return
	}

	if c.hasFloor && g.y == c.floorLevel-1 {
		g.rest()
		return
	}

	if c.grainCanMoveDown(*g) {
		g.y++
	} else if c.grainCanMoveLeft(*g) {
		g.x--
		g.y++
	} else if c.grainCanMoveRight(*g) {
		g.x++
		g.y++
	} else {
		g.rest()
	}
}

func (c *sandCave) clearSand() {
	for coord, entity := range c.entities {
		if _, ok := entity.(grain); ok {
			delete(c.entities, coord)
		}
	}
	c.overflowed = false
	c.full = false
}

func newCave() *sandCave {
	return &sandCave{
		entities: make(map[coord]caveEntity),
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	cave := newCave()
	for scanner.Scan() {
		line := scanner.Text()
		coordStrings := strings.Split(line, " -> ")
		for i := 1; i < len(coordStrings); i++ {
			start := parseCoord(coordStrings[i-1])
			end := parseCoord(coordStrings[i])
			cave.insertWall(start, end)
		}
	}

	part1(cave)
	cave.clearSand()
	part2(cave)
}

func part1(cave *sandCave) {
	cave.pourSand()

	restingAmount := cave.getRestingSandAmount()
	fmt.Println(
		"The amount of sand resting before a grain overflows is:",
		restingAmount,
	)
}

func part2(cave *sandCave) {
	cave.insertFloor()
	cave.pourSand()

	restingAmount := cave.getRestingSandAmount()
	fmt.Println(
		"The amount of sand resting once the entry hole is blocked is:",
		restingAmount,
	)
}
