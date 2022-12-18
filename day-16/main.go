package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	infinity             = math.MaxInt
	startValve           = "AA"
	startMinutes         = 30
	startMinutesElephant = startMinutes - 4
)

var inputLineRegex = regexp.MustCompile(
	`Valve (\w{2}) has flow rate=(\d+); tunnels? leads? to valves? ([\w, ]+)`,
)

type tunnel struct {
	lengthInMinutes int
	destination     *caveValve
}

type caveValve struct {
	label       string
	flowRate    int
	connections []*caveValve
	tunnels     []tunnel
}

type tunnelSystem struct {
	valves           map[string]*caveValve
	visited          map[string]int
	maxPressureSoFar int
	queue            []caveValve
}

func (s *tunnelSystem) findValveDistance(start, end caveValve) int {
	dist := make(map[string]int)
	s.queue = make([]caveValve, 0)

	for label, valve := range s.valves {
		dist[label] = infinity
		s.queue = append(s.queue, *valve)
	}
	dist[start.label] = 0

	for len(s.queue) > 0 {
		valve := s.findClosestInQueue(dist)
		d := dist[valve.label]

		if d == infinity {
			break
		}
		if valve.label == end.label {
			return dist[valve.label]
		}

		for _, v := range s.queue {
			for _, n := range valve.connections {
				if v.label == n.label {
					alt := d + 1
					if alt < dist[n.label] {
						dist[n.label] = alt
					}
				}
			}
		}
	}

	return infinity
}

func (s *tunnelSystem) findClosestInQueue(dist map[string]int) caveValve {
	closest := s.queue[0]
	closestDistance := dist[closest.label]
	closestIndex := 0
	for i, v := range s.queue {
		d := dist[v.label]
		if d < closestDistance {
			closest = v
			closestDistance = d
			closestIndex = i
		}
	}
	s.queue = append(s.queue[:closestIndex], s.queue[closestIndex+1:]...)

	return closest
}

func (s *tunnelSystem) findMaxPressureRelease(
	valve caveValve, released, minsRemaining int) int {

	v := s.visited[valve.label]
	s.visited[valve.label] = v + 1

	if valve.flowRate != 0 {
		minsRemaining--
		released += (valve.flowRate * minsRemaining)
	}

	maxRelease := released
	for _, t := range valve.tunnels {
		if t.lengthInMinutes >= minsRemaining {
			continue
		}

		dest := t.destination
		if dest.flowRate == 0 {
			continue
		}

		ok := s.visited[dest.label]
		if ok > 0 {
			continue
		}

		r := s.findMaxPressureRelease(
			*dest, released, minsRemaining-t.lengthInMinutes,
		)
		if r > maxRelease {
			maxRelease = r
		}
	}

	v = s.visited[valve.label]
	s.visited[valve.label] = v - 1
	return maxRelease
}

func (s *tunnelSystem) findMaxPressureReleaseElephant(
	valve, eValve caveValve,
	minsRemaining, eMinsRemaining, released int) int {

	meVisited := s.visited[valve.label]
	elVisited := s.visited[eValve.label]

	if valve.flowRate != 0 && meVisited < 1 {
		minsRemaining--
		released += (valve.flowRate * minsRemaining)
	}
	if eValve.flowRate != 0 && elVisited < 1 {
		eMinsRemaining--
		released += (eValve.flowRate * eMinsRemaining)
	}

	v := s.visited[valve.label]
	s.visited[valve.label] = v + 1
	v = s.visited[eValve.label]
	s.visited[eValve.label] = v + 1

	isElephant := false
	chosenValve := &valve
	if eMinsRemaining > minsRemaining {
		isElephant = true
		chosenValve = &eValve
	}

	maxRelease := released
	for _, t := range chosenValve.tunnels {
		if isElephant {
			if t.lengthInMinutes >= eMinsRemaining {
				continue
			}
		} else {
			if t.lengthInMinutes >= minsRemaining {
				continue
			}
		}

		dest := t.destination
		if dest.flowRate == 0 {
			continue
		}

		ok := s.visited[dest.label]
		if ok > 0 {
			continue
		}

		var r int
		if isElephant {
			r = s.findMaxPressureReleaseElephant(
				valve, *dest,
				minsRemaining, eMinsRemaining-t.lengthInMinutes,
				released,
			)
		} else {
			r = s.findMaxPressureReleaseElephant(
				*dest, eValve,
				minsRemaining-t.lengthInMinutes, eMinsRemaining,
				released,
			)
		}

		if r > maxRelease {
			maxRelease = r
		}
	}

	v = s.visited[valve.label]
	s.visited[valve.label] = v - 1
	v = s.visited[eValve.label]
	s.visited[eValve.label] = v - 1

	if maxRelease > s.maxPressureSoFar {
		s.maxPressureSoFar = maxRelease
	}

	return maxRelease
}

func newTunnelSystem() *tunnelSystem {
	return &tunnelSystem{
		valves:  make(map[string]*caveValve),
		visited: make(map[string]int),
	}
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	tunnels := newTunnelSystem()
	matches := inputLineRegex.FindAllStringSubmatch(string(bytes), -1)
	for _, m := range matches {
		var (
			label      = m[1]
			rateString = m[2]
		)

		flowRate, _ := strconv.Atoi(rateString)

		tunnels.valves[label] = &caveValve{
			label:       label,
			flowRate:    flowRate,
			connections: make([]*caveValve, 0),
		}
	}

	for _, m := range matches {
		var (
			label             = m[1]
			connectionsString = m[3]
		)

		connections := make([]*caveValve, 0)
		for _, c := range strings.Split(connectionsString, ", ") {
			v := tunnels.valves[c]
			connections = append(connections, v)
		}

		valve := tunnels.valves[label]
		valve.connections = connections
	}

	for _, v1 := range tunnels.valves {
		for _, v2 := range tunnels.valves {
			if v1.label == v2.label {
				continue
			}
			if v2.flowRate == 0 {
				continue
			}

			d := tunnels.findValveDistance(*v1, *v2)
			t1 := tunnel{
				lengthInMinutes: d,
				destination:     v2,
			}

			v1.tunnels = append(v1.tunnels, t1)
		}
	}

	part1(tunnels)
	part2(tunnels)
}

func part1(tunnels *tunnelSystem) {
	start := time.Now()

	bestPressureRelease := tunnels.findMaxPressureRelease(
		*tunnels.valves[startValve], 0, startMinutes,
	)
	fmt.Println(
		"The most pressure that can be released in 30 minutes is:",
		bestPressureRelease,
	)

	fmt.Printf("Part 1 took: %dms\n", time.Since(start).Milliseconds())
}

func part2(tunnels *tunnelSystem) {
	start := time.Now()

	meStart := tunnels.valves[startValve]
	meMins := startMinutesElephant
	elStart := tunnels.valves[startValve]
	elMins := startMinutesElephant

	v := tunnels.valves[startValve]
	if len(v.tunnels) == 2 {
		meStart = v.tunnels[0].destination
		meMins -= v.tunnels[0].lengthInMinutes
		elStart = v.tunnels[1].destination
		elMins -= v.tunnels[1].lengthInMinutes
	}

	bestPressureRelease := tunnels.findMaxPressureReleaseElephant(
		*meStart,
		*elStart,
		meMins,
		elMins,
		0,
	)
	fmt.Println(
		"The most pressure that can be released in 26 minutes with an elephant is:",
		bestPressureRelease,
	)

	fmt.Printf("Part 2 took: %.2fs\n", time.Since(start).Seconds())
}
