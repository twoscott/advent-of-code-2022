package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

const (
	minHeight = 0
	maxHeight = 9
)

type coord struct {
	x, y int
}

type treeFarm struct {
	trees [][]int
}

func (f *treeFarm) parseTrees(file *os.File) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	i := 0
	for fileScanner.Scan() {
		row := fileScanner.Text()
		f.trees = append(f.trees, make([]int, len(row)))
		for j := range row {
			f.trees[i][j], _ = strconv.Atoi(string(row[j]))

		}

		i++
	}
}

func (f treeFarm) findVisibleTrees() int {
	visible := make(map[coord]bool)

	for h := minHeight; h <= maxHeight; h++ {
		for i, row := range f.trees {
			for j, height := range row {
				if height == h {
					visible[coord{x: j, y: i}] = true
				}
				if height >= h {
					break
				}
			}
			for j := len(row) - 1; j >= 0; j-- {
				height := row[j]
				if height == h {
					visible[coord{x: j, y: i}] = true
				}
				if height >= h {
					break
				}
			}
		}
		for i := 0; i < len(f.trees[0]); i++ {
			for j := 0; j < len(f.trees); j++ {
				height := f.trees[j][i]
				if height == h {
					visible[coord{x: i, y: j}] = true
				}
				if height >= h {
					break
				}
			}
			for j := len(f.trees) - 1; j >= 0; j-- {
				height := f.trees[j][i]
				if height == h {
					visible[coord{x: i, y: j}] = true
				}
				if height >= h {
					break
				}
			}
		}
	}

	return len(visible)
}

func (f treeFarm) findMostScenicScore() int {
	max := 0

	for i, row := range f.trees {
		for j, height := range row {
			score := f.getTreeScenicScore(j, i, height)
			if score > max {
				max = score
			}
		}
	}

	return max
}

func (f treeFarm) getTreeScenicScore(x, y, treeHeight int) int {
	leftScore := 0
	for i := x - 1; i >= 0; i-- {
		currHeight := f.trees[y][i]
		leftScore++
		if currHeight >= treeHeight {
			break
		}
	}

	rightScore := 0
	for i := x + 1; i < len(f.trees[y]); i++ {
		currHeight := f.trees[y][i]
		rightScore++
		if currHeight >= treeHeight {
			break
		}
	}

	topScore := 0
	for i := y - 1; i >= 0; i-- {
		currHeight := f.trees[i][x]
		topScore++
		if currHeight >= treeHeight {
			break
		}
	}

	bottomScore := 0
	for i := y + 1; i < len(f.trees); i++ {
		currHeight := f.trees[i][x]
		bottomScore++
		if currHeight >= treeHeight {
			break
		}
	}

	return leftScore * rightScore * topScore * bottomScore
}

func newTreeFarm() *treeFarm {
	return &treeFarm{
		trees: make([][]int, 0),
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	farm := newTreeFarm()
	farm.parseTrees(file)
	visible := farm.findVisibleTrees()
	log.Println(farm.getTreeScenicScore(0, 0, farm.trees[0][0]))
	mostScenic := farm.findMostScenicScore()

	log.Println("The total number of visible trees is:", visible)
	log.Println("The highest scenic score of the farm is:", mostScenic)
}
