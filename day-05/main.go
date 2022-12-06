package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type craneInstruction struct {
	source, dest, amount int
}

func parseCraneInstructions(input string) []craneInstruction {
	lines := strings.Split(input, "\n")
	instructions := make([]craneInstruction, 0, len(lines))

	for _, line := range lines {
		var instruction craneInstruction
		fmt.Sscanf(
			line,
			"move %d from %d to %d",
			&instruction.amount, &instruction.source, &instruction.dest,
		)

		instructions = append(instructions, instruction)
	}

	return instructions
}

type crateStack struct {
	crates [][]byte
}

func (s crateStack) String() string {
	return string(bytes.Join(s.crates, []byte("\n")))
}

func (s crateStack) move9000(inst craneInstruction) {
	sIdx := inst.source - 1
	dIdx := inst.dest - 1

	for i := 0; i < inst.amount; i++ {
		c := s.crates[sIdx][len(s.crates[sIdx])-1]
		s.crates[sIdx] = s.crates[sIdx][:len(s.crates[sIdx])-1]

		s.crates[dIdx] = append(s.crates[dIdx], c)
	}
}

func (s crateStack) move9001(inst craneInstruction) {
	sIdx := inst.source - 1
	dIdx := inst.dest - 1

	c := s.crates[sIdx][len(s.crates[sIdx])-inst.amount:]
	s.crates[sIdx] = s.crates[sIdx][:len(s.crates[sIdx])-inst.amount]
	s.crates[dIdx] = append(s.crates[dIdx], c...)
}

func (s crateStack) getTopCrates() string {
	chars := []byte{}

	for _, column := range s.crates {
		columnChars := bytes.Trim(column, "\x00")
		if len(columnChars) < 1 {
			continue
		}

		chars = append(chars, columnChars[len(columnChars)-1])
	}

	return string(chars)
}

func parseCrates(input string) [][]byte {
	inputRows := strings.Split(input, "\n")

	rows := [][]byte{}
	for _, inputRow := range inputRows {
		if inputRow[0] != '[' {
			break
		}

		row := []byte{}
		for i := 0; i < len(inputRow); i += 4 {
			var crate byte
			fmt.Sscanf(inputRow[i:], "[%c]", &crate)
			row = append(row, crate)
		}

		rows = append(rows, row)
	}

	columns := make([][]byte, len(rows[0]))
	for i := range columns {
		columns[i] = make([]byte, len(rows))
	}

	for i := range rows {
		for j := range rows[i] {
			idx := len(rows) - i - 1
			columns[j][idx] = rows[i][j]
		}
	}

	for i := range columns {
		columns[i] = bytes.Trim(columns[i], "\x00")
	}

	return columns
}

func main() {
	bytes, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	input := strings.TrimSpace(string(bytes))

	sections := strings.SplitAfterN(input, "9", 2)
	if len(sections) < 2 {
		panic("not enough input sections")
	}

	columns := parseCrates(sections[0])
	columns2 := make([][]byte, len(columns))
	for i, c := range columns {
		columns2[i] = make([]byte, len(columns[i]))
		copy(columns2[i], c)
	}

	stack := crateStack{crates: columns}
	stack2 := crateStack{crates: columns2}

	moveSteps := parseCraneInstructions(strings.TrimSpace(sections[1]))

	part1(stack, moveSteps)
	part2(stack2, moveSteps)
}

func part1(stack crateStack, moveSteps []craneInstruction) {
	for _, step := range moveSteps {
		stack.move9000(step)
	}

	topCrates := stack.getTopCrates()
	log.Println("The top crates after sorting are:", topCrates)
}

func part2(stack crateStack, moveSteps []craneInstruction) {
	for _, step := range moveSteps {
		stack.move9001(step)
	}

	topCrates := stack.getTopCrates()
	log.Println("The top crates after sorting are:", topCrates)
}
