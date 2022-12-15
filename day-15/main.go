package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

const (
	part1Row = 2_000_000
	part2Min = 0
	part2Max = 4_000_000
)

func intDiff(i1, i2 int) int {
	return int(math.Abs(float64(i1 - i2)))
}

func manhattanDistance(c1, c2 coord) int {
	return intDiff(c1.x, c2.x) + intDiff(c1.y, c2.y)
}

type coord struct {
	x, y int
}

type caveSensor struct {
	coord
	nearestBeacon coord
}

func (s caveSensor) distanceToBeacon() int {
	return manhattanDistance(s.coord, s.nearestBeacon)
}

type distressCave struct {
	sensors    map[coord]caveSensor
	emptyTiles map[coord]bool
}

func (c distressCave) countEmptyTilesInRow(y int) (count int) {
	var minX, maxX int
	var minSet bool

	for pos, sensor := range c.sensors {
		x1 := pos.x - sensor.distanceToBeacon()
		x2 := pos.x + sensor.distanceToBeacon()
		if !minSet {
			minX = x1
			maxX = x2
			minSet = true
			continue
		}

		if x1 < minX {
			minX = x1
		}
		if x2 > maxX {
			maxX = x2
		}
	}

	x := minX
	for x <= maxX {
		outOfRangeCount := 0

		for pos, sensor := range c.sensors {
			if manhattanDistance(sensor.coord, coord{x: x, y: y}) >
				sensor.distanceToBeacon() {

				outOfRangeCount++
				continue
			}

			distance := sensor.distanceToBeacon()
			diff := intDiff(pos.y, y)

			distanceDiff := distance - diff
			newX := pos.x + distanceDiff + 1
			count += (newX - x)
			x = newX
		}

		if outOfRangeCount == len(c.sensors) {
			x++
		}
	}

	negatedBeacons := make(map[coord]bool)
	for _, s := range c.sensors {
		if s.nearestBeacon.y == y {
			if _, ok := negatedBeacons[s.nearestBeacon]; !ok {
				count--
				negatedBeacons[s.nearestBeacon] = true
			}
		}
	}

	return count
}

func (c distressCave) findDistressSignal() coord {
	x := part2Min
	y := part2Min

	for x <= part2Max && y <= part2Max {
		outOfRangeCount := 0

		for pos, sensor := range c.sensors {
			if manhattanDistance(sensor.coord, coord{x: x, y: y}) >
				sensor.distanceToBeacon() {

				outOfRangeCount++
				continue
			}

			distance := sensor.distanceToBeacon()
			diff := intDiff(pos.y, y)

			distanceDiff := distance - diff
			x = pos.x + distanceDiff + 1
			if x > part2Max {
				x = 0
				y++
			}
		}

		if outOfRangeCount == len(c.sensors) {
			return coord{x: x, y: y}
		}
	}

	return coord{}
}

func newDistressCave() *distressCave {
	return &distressCave{
		sensors:    make(map[coord]caveSensor),
		emptyTiles: make(map[coord]bool),
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

	cave := newDistressCave()
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		var sensorCoord coord
		var beaconCoord coord
		fmt.Sscanf(
			line,
			"Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&sensorCoord.x, &sensorCoord.y,
			&beaconCoord.x, &beaconCoord.y,
		)

		cave.sensors[sensorCoord] = caveSensor{
			coord: coord{
				x: sensorCoord.x, y: sensorCoord.y,
			},
			nearestBeacon: beaconCoord,
		}
	}

	part1(*cave)
	part2(*cave)
}

func part1(cave distressCave) {
	start := time.Now()

	tilesCount := cave.countEmptyTilesInRow(part1Row)
	fmt.Printf(
		"The number of tiles that cannot contain a beacon in row %d are: %d\n",
		part1Row, tilesCount,
	)

	fmt.Printf("Part 1 took: %dms\n", time.Since(start).Milliseconds())
}

func part2(cave distressCave) {
	start := time.Now()

	signalLocation := cave.findDistressSignal()
	tuningFrequency := (signalLocation.x * part2Max) + signalLocation.y
	fmt.Println(
		"The tuning frequency for the distress beacon is:",
		tuningFrequency,
	)

	fmt.Printf("Part 2 took: %dms\n", time.Since(start).Milliseconds())
}
