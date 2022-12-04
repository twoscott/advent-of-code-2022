package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type elfPair struct {
	range1 [2]int
	range2 [2]int
}

func (p elfPair) isContained() bool {
	return p.range1[0] >= p.range2[0] && p.range1[1] <= p.range2[1] ||
		p.range2[0] >= p.range1[0] && p.range2[1] <= p.range1[1]
}

func (p elfPair) isOverlapped() bool {
	return p.range1[1] >= p.range2[0] && p.range1[0] <= p.range2[1]
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := strings.TrimSpace(string(bytes))
	pairs := getPairs(input)

	part1(pairs)
	part2(pairs)
}

func getPairs(input string) []elfPair {
	lines := strings.Split(input, "\n")
	pairs := make([]elfPair, 0, len(lines))

	for _, line := range lines {
		var pair elfPair
		fmt.Sscanf(
			line,
			"%d-%d,%d-%d",
			&pair.range1[0], &pair.range1[1], &pair.range2[0], &pair.range2[1],
		)

		pairs = append(pairs, pair)
	}

	return pairs
}

func part1(pairs []elfPair) {
	totalContained := 0

	for _, pair := range pairs {
		if pair.isContained() {
			totalContained++
		}
	}

	fmt.Println(
		"The total number of contained assignment pairs is:",
		totalContained,
	)
}

func part2(pairs []elfPair) {
	totalOverlapped := 0

	for _, pair := range pairs {
		if pair.isOverlapped() {
			totalOverlapped++
		}
	}

	fmt.Println(
		"The total number of overlapping assignment pairs is:",
		totalOverlapped,
	)
}
