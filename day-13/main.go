package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

const (
	divider1Num = 2.0
	divider2Num = 6.0
)

type packetPair [2][]interface{}

func compareOrders(packet1, packet2 interface{}) (comp int) {
	float1, p1IsFloat := packet1.(float64)
	float2, p2IsFloat := packet2.(float64)
	array1, p1IsArray := packet1.([]interface{})
	array2, p2IsArray := packet2.([]interface{})

	switch {
	case p1IsFloat && p2IsFloat:
		comp = int(float1 - float2)
	case p1IsArray && p2IsArray:
		for i := 0; comp == 0; i++ {
			if i >= len(array1) || i >= len(array2) {
				comp = len(array1) - len(array2)
				break
			}

			comp = compareOrders(array1[i], array2[i])
		}
	case p1IsFloat && p2IsArray:
		comp = compareOrders([]interface{}{float1}, array2)
	case p1IsArray && p2IsFloat:
		comp = compareOrders(array1, []interface{}{float2})
	default:
		log.Panicln(
			"invalid types found:",
			reflect.TypeOf(packet1),
			reflect.TypeOf(packet2),
		)
	}
	return
}

func sortPackets(packets [][]interface{}) [][]interface{} {
	n := len(packets)
	for i := 0; i < n; i++ {
		for j := 0; j < n-1-i; j++ {
			if compareOrders(packets[j], packets[j+1]) > 0 {
				packets[j], packets[j+1] = packets[j+1], packets[j]
			}
		}
	}

	return packets
}

func getDecoderKey(packets [][]interface{}, divNums ...float64) int {
	key := 1

	for i, p := range packets {
		if len(p) != 1 {
			continue
		}

		inner := p[0]
		if val, ok := inner.([]interface{}); ok {
			if len(val) != 1 {
				continue
			}
			inner = val[0]
			if val, ok := inner.(float64); ok {
				for _, div := range divNums {
					if val == div {
						key *= (i + 1)
					}
				}
			}
		}
	}

	return key
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	packets := make([][]interface{}, 0)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var packet []interface{}
		json.Unmarshal(line, &packet)
		packets = append(packets, packet)
	}

	packetPairs := make([]packetPair, 0, len(packets)/2)
	for i := 1; i < len(packets); i += 2 {
		left, right := packets[i-1], packets[i]
		packetPairs = append(packetPairs, packetPair{left, right})
	}

	part1(packetPairs)
	part2(packets)
}

func part1(pairs []packetPair) {
	correctIndicesSum := 0
	for i, pair := range pairs {
		pairNum := i + 1
		comp := compareOrders(pair[0], pair[1])
		if comp < 0 {
			correctIndicesSum += pairNum
		}
	}

	fmt.Println(
		"The sum of the indices of the correctly ordered pairs of packets is:",
		correctIndicesSum,
	)
}

func part2(packets [][]interface{}) {
	div1 := []interface{}{
		[]interface{}{divider1Num},
	}
	div2 := []interface{}{
		[]interface{}{divider2Num},
	}

	packets = append(packets, div1, div2)

	sortedPackets := sortPackets(packets)
	key := getDecoderKey(sortedPackets, divider1Num, divider2Num)
	fmt.Println("The decoder key for the distress signal is:", key)
}
