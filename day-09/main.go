package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

const ropeKnotsCount = 10

type moveDirection byte

const (
	leftDirection  moveDirection = 'L'
	rightDirection moveDirection = 'R'
	upDirection    moveDirection = 'U'
	downDirection  moveDirection = 'D'
)

type coord struct {
	x, y int
}

func (c coord) diff(c2 coord) coord {
	return coord{
		x: c.x - c2.x,
		y: c.y - c2.y,
	}
}

type knot struct {
	position coord
	visited  map[coord]bool
}

func (k *knot) moveTo(newPos coord) {
	k.position = newPos
	k.visited[k.position] = true
}

func (k *knot) moveBy(posDiff coord) {
	newPos := k.position
	newPos.x += posDiff.x
	newPos.y += posDiff.y
	k.moveTo(newPos)
}

func (k *knot) moveDirection(dir moveDirection) {
	newPos := k.position

	switch dir {
	case leftDirection:
		newPos.x--
	case rightDirection:
		newPos.x++
	case upDirection:
		newPos.y++
	case downDirection:
		newPos.y--
	default:
		panic("unable to determine direction")
	}

	k.moveTo(newPos)
}

func (k knot) isAdjacent(comp knot) bool {
	diff := k.position.diff(comp.position)
	diff.x = int(math.Abs(float64(diff.x)))
	diff.y = int(math.Abs(float64(diff.y)))

	return diff.x <= 1 && diff.y <= 1
}

func newKnot() *knot {
	start := coord{0, 0}
	return &knot{
		position: start,
		visited:  map[coord]bool{start: true},
	}
}

type bridgeRope []*knot

func (r bridgeRope) head() *knot {
	return r[0]
}

func (r bridgeRope) tail() bridgeRope {
	return r[1:]
}

func (r bridgeRope) last() *knot {
	return r[len(r)-1]
}

func (r bridgeRope) String() string {
	var grid string
	minCoord := r.head().position
	maxCoord := r.last().position

	for _, knot := range r {
		if knot.position.x < minCoord.x {
			minCoord.x = knot.position.x
		}
		if knot.position.x > maxCoord.x {
			maxCoord.x = knot.position.x
		}
		if knot.position.y < minCoord.y {
			minCoord.y = knot.position.y
		}
		if knot.position.y > maxCoord.y {
			maxCoord.y = knot.position.y
		}
	}

	for i := maxCoord.y; i >= minCoord.y; i-- {
		row := ""
		for j := minCoord.x; j <= maxCoord.x; j++ {
			var knotAtPosNum string
			for num, knot := range r {
				if knot.position.x == j && knot.position.y == i {
					if num == 0 {
						knotAtPosNum = "H"
					} else {
						knotAtPosNum = strconv.Itoa(num)
					}
					break
				}
			}
			if knotAtPosNum != "" {
				row += knotAtPosNum
			} else {
				row += "+"
			}
		}
		grid += row + "\n"
	}

	return grid
}

type moveInstruction struct {
	direction moveDirection
	amount    int
}

func (i moveInstruction) execute(rope bridgeRope) {
	head := rope.head()
	tails := rope.tail()

	for m := 0; m < i.amount; m++ {
		head.moveDirection(i.direction)

		leadingKnot := head
		for _, tail := range tails {
			if tail.isAdjacent(*leadingKnot) {
				break
			}
			diff := leadingKnot.position.diff(tail.position)
			move := coord{}
			if diff.x > 0 {
				move.x = 1
			}
			if diff.x < 0 {
				move.x = -1
			}
			if diff.y > 0 {
				move.y = 1
			}
			if diff.y < 0 {
				move.y = -1
			}

			tail.moveBy(move)
			leadingKnot = tail
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	instructions := make([]moveInstruction, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		var inst moveInstruction
		fmt.Sscanf(line, "%c %d", &inst.direction, &inst.amount)

		instructions = append(instructions, inst)
	}

	part1(instructions)
	part2(instructions)
}

func part1(instructions []moveInstruction) {
	head := newKnot()
	tail := newKnot()
	rope := bridgeRope{head, tail}

	for _, inst := range instructions {
		inst.execute(rope)
	}

	log.Printf(
		"The tail knot in the two-knot rope moved to %d different positions.\n",
		len(tail.visited),
	)
}

func part2(instructions []moveInstruction) {
	rope := make(bridgeRope, ropeKnotsCount)
	for i := range rope {
		rope[i] = newKnot()
	}

	for _, inst := range instructions {
		inst.execute(rope)
	}

	tail := rope.last()
	log.Printf(
		"The tail knot in the ten-knot long rope moved to %d different positions.\n",
		len(tail.visited),
	)
}
