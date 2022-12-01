package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	maxCal := 0

	inventories := strings.Split(input, "\n\n")
	for _, inv := range inventories {
		totalCal := 0

		calories := strings.Split(inv, "\n")

		for _, cal := range calories {
			intCal, _ := strconv.Atoi(cal)
			totalCal += intCal
		}

		if totalCal > maxCal {
			maxCal = totalCal
		}
	}

	fmt.Println("Most calories carried by an elf:", maxCal)
}

func part2(input string) {
	topCals := make([]int, 3)

	inventories := strings.Split(input, "\n\n")
	for _, inv := range inventories {
		totalCal := 0

		calories := strings.Split(inv, "\n")
		for _, cal := range calories {
			intCal, _ := strconv.Atoi(cal)
			totalCal += intCal
		}

		for i, cal := range topCals {
			if totalCal > cal {
				topCals = append(
					topCals[:i],
					append([]int{totalCal}, topCals[i:2]...)...,
				)
				break
			}
		}
	}

	topThreeSum := 0
	for _, cal := range topCals {
		topThreeSum += cal
	}

	fmt.Println(
		"Total calories from the 3 elves carrying the most calories:",
		topThreeSum,
	)
}
