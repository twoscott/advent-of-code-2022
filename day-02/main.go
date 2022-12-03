package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := string(bytes)

	part1(input)
	part2(input)
}

func part1(input string) {
	var totalScore int32 = 0

	for _, rps := range strings.Split(input, "\n") {
		if rps == "" {
			continue
		}

		var opponent, response rune
		fmt.Sscanf(rps, "%c %c", &opponent, &response)
		totalScore += getMatchResult(opponent, response)
	}

	log.Println(
		"The total score from the rock paper scissors matches is:",
		totalScore,
	)
}

func getMatchResult(opponent, response rune) int32 {
	opponent -= 65
	response -= 88

	score := response
	switch {
	case response == (opponent+1)%3: //win
		return 1 + score + 6
	case response == opponent: // draw
		return 1 + score + 3
	case response == (opponent+2)%3: // lose
		return 1 + score
	default: // unknown
		return 0
	}
}

func part2(input string) {
	var totalScore int32 = 0

	for _, rps := range strings.Split(input, "\n") {
		if rps == "" {
			continue
		}

		var rpsVal, result rune
		fmt.Sscanf(rps, "%c %c", &rpsVal, &result)
		totalScore += simulateMatch(rpsVal, result)
	}

	log.Println(
		"The total score from simulating the winning matches is:",
		totalScore,
	)
}

func simulateMatch(rpsVal, result rune) int32 {
	rpsVal -= 65

	switch result {
	case 'X': // lose
		return 1 + (rpsVal+2)%3
	case 'Y': // draw
		return 1 + 3 + rpsVal
	case 'Z': // win
		return 1 + 6 + (rpsVal+1)%3
	default: // unknown
		return 0
	}
}
