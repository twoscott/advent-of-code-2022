package main

import (
	"bytes"
	"log"
	"os"
)

const (
	startOfPacketSize  = 4
	startOfMessageSize = 14
)

func findUniquePacket(input []byte, packetSize int) int {
	for i := packetSize; i < len(input); i++ {
		total := 0
		chunk := input[i-packetSize : i]
		for _, c := range chunk {
			total += bytes.Count(chunk, []byte{c})
		}
		if total == packetSize {
			return i
		}
	}
	return 0
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	startOfPacket := findUniquePacket(input, startOfPacketSize)
	startOfMessage := findUniquePacket(input, startOfMessageSize)

	log.Println("The start of the packet is at:", startOfPacket)
	log.Println("The start of the message is at:", startOfMessage)
}
