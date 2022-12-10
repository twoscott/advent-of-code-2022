package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	noopInstruction = "noop"
	addxInstruction = "addx"

	crtWidth  = 40
	crtHeight = 6
)

type register struct {
	value   int
	history map[uint]int
}

func (r *register) addValue(val int, cycle uint) {
	r.value += val
	r.history[cycle] = r.value
}

type deviceCPU struct {
	cycle     uint
	xRegister *register
}

func (c *deviceCPU) executeInstruction(inst instruction) {
	switch inst.name {
	case noopInstruction:
		c.cycle++
	case addxInstruction:
		c.cycle += 2
		c.xRegister.addValue(inst.value, c.cycle)
	default:
		panic("unknown instruction name")
	}
}

func (c deviceCPU) getValueAtCycle(cycle uint) int {
	val, ok := 0, false
	for !ok && cycle > 0 {
		val, ok = c.xRegister.history[cycle]
		cycle--
	}

	return val
}

func (c deviceCPU) getSignalStrength(cycle uint) int {
	return int(cycle) * c.getValueAtCycle(cycle)
}

func (c deviceCPU) renderImage() string {
	crt := strings.Builder{}
	crt.Grow(crtHeight*crtWidth + crtHeight)
	
	for i := 0; i < crtHeight; i++ {
		for j := 1; j <= crtWidth; j++ {
			cycle := uint((i * crtWidth) + j)
			spriteX := c.getValueAtCycle(cycle) + 1

			if j >= spriteX-1 && j <= spriteX+1 {
				crt.WriteByte('#')
			} else {
				crt.WriteByte('.')
			}
		}
		crt.WriteByte('\n')
	}

	return crt.String()
}

func newDeviceCPU() *deviceCPU {
	return &deviceCPU{
		cycle: 1,
		xRegister: &register{
			value: 1,
			history: map[uint]int{
				1: 1,
			},
		},
	}
}

type instruction struct {
	name  string
	value int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	cpu := newDeviceCPU()

	for fileScanner.Scan() {
		line := fileScanner.Text()
		var inst instruction
		fmt.Sscanf(line, "%s %d", &inst.name, &inst.value)

		cpu.executeInstruction(inst)
	}

	part1(cpu)
	part2(cpu)
}

func part1(cpu *deviceCPU) {
	signalStrengthSum := 0
	for c := 20; uint(c) < cpu.cycle; c += 40 {
		signalStrengthSum += cpu.getSignalStrength(uint(c))
	}

	fmt.Println("Signal strength sum for cycles:", signalStrengthSum)
}

func part2(cpu *deviceCPU) {
	crtImage := cpu.renderImage()
	fmt.Println("The following is the image rendered onto the CRT:")
	fmt.Print(crtImage)
}
