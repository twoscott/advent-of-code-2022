package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFile = "input.txt"

func getPriority(char rune) int {
	if char >= 'a' && char <= 'z' {
		return 1 + int(char) - 'a'
	}
	if char >= 'A' && char <= 'Z' {
		return 27 + int(char) - 'A'
	}

	return 0
}

func main() {
	part1()
	part2()
}

func part1() {
	prioritiesSum := 0

	bytes, err := os.ReadFile(inputFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := strings.TrimSpace(string(bytes))

	rucksacks := strings.Split(input, "\n")
	for _, r := range rucksacks {
		compartment1 := r[:len(r)/2]
		compartment2 := r[len(r)/2:]

		i := strings.IndexAny(compartment1, compartment2)
		if i == -1 {
			panic("couldn't find item type")
		}

		prioritiesSum += getPriority(rune(compartment1[i]))
	}

	fmt.Println("The sum of the incorrect items is:", prioritiesSum)
}

func part2() {
	prioritiesSum := 0

	file, err := os.Open(inputFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		var sack1, sack2, sack3 string
		_, err := fmt.Fscanf(file, "%s\n%s\n%s\n", &sack1, &sack2, &sack3)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panicln(err)
		}

		for _, c := range sack1 {
			if strings.ContainsRune(sack2, c) &&
				strings.ContainsRune(sack3, c) {

				prioritiesSum += getPriority(c)
				break
			}
		}
	}

	fmt.Println("The sum of the misplaced badges is:", prioritiesSum)
}
