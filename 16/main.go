package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

type point struct {
	x int
	y int
}

type laser struct {
	x         int
	y         int
	direction rune
}

type Queue []laser

func (queue *Queue) enqueue(l laser) {
	*queue = append(*queue, l)
}

func (queue *Queue) enqueueSlice(l []laser) {
	*queue = append(*queue, l...)
}

func (queue *Queue) dequeue() laser {
	l := (*queue)[0]
	*queue = (*queue)[1:]
	return l
}

func parseMirrors(scanner *bufio.Scanner) [][]rune {
	mirrors := [][]rune{}
	for scanner.Scan() {
		mirrors = append(mirrors, []rune(scanner.Text()))
	}
	return mirrors
}

func validLaser(l laser, maxX int, maxY int) bool {
	return l.x >= 0 && l.x < maxX && l.y >= 0 && l.y < maxY
}

func filterValidLasers(ls []laser, maxX int, maxY int) []laser {
	n := 0
	for _, l := range ls {
		if validLaser(l, maxX, maxY) {
			ls[n] = l
			n++
		}
	}
	return ls[:n]
}

func nextLasers(l laser, mirrors [][]rune, encounteredSplitters map[point]bool) []laser {
	newLasers := []laser{}
	switch mirrors[l.y][l.x] {
	case '.':
		switch l.direction {
		case 'r':
			newLasers = append(newLasers, laser{x: l.x + 1, y: l.y, direction: l.direction})
		case 'l':
			newLasers = append(newLasers, laser{x: l.x - 1, y: l.y, direction: l.direction})
		case 'u':
			newLasers = append(newLasers, laser{x: l.x, y: l.y - 1, direction: l.direction})
		case 'd':
			newLasers = append(newLasers, laser{x: l.x, y: l.y + 1, direction: l.direction})
		}
	case '/':
		switch l.direction {
		case 'r':
			newLasers = append(newLasers, laser{x: l.x, y: l.y - 1, direction: 'u'})
		case 'l':
			newLasers = append(newLasers, laser{x: l.x, y: l.y + 1, direction: 'd'})
		case 'u':
			newLasers = append(newLasers, laser{x: l.x + 1, y: l.y, direction: 'r'})
		case 'd':
			newLasers = append(newLasers, laser{x: l.x - 1, y: l.y, direction: 'l'})
		}
	case '\\':
		switch l.direction {
		case 'r':
			newLasers = append(newLasers, laser{x: l.x, y: l.y + 1, direction: 'd'})
		case 'l':
			newLasers = append(newLasers, laser{x: l.x, y: l.y - 1, direction: 'u'})
		case 'u':
			newLasers = append(newLasers, laser{x: l.x - 1, y: l.y, direction: 'l'})
		case 'd':
			newLasers = append(newLasers, laser{x: l.x + 1, y: l.y, direction: 'r'})
		}
	case '-':
		switch l.direction {
		case 'r':
			newLasers = append(newLasers, laser{x: l.x + 1, y: l.y, direction: l.direction})
		case 'l':
			newLasers = append(newLasers, laser{x: l.x - 1, y: l.y, direction: l.direction})
		case 'u', 'd':
			p := point{x: l.x, y: l.y}
			if !encounteredSplitters[p] {
				newLasers = append(
					newLasers,
					laser{x: l.x - 1, y: l.y, direction: 'l'},
					laser{x: l.x + 1, y: l.y, direction: 'r'},
				)
				encounteredSplitters[p] = true
			}
		}
	case '|':
		switch l.direction {
		case 'u':
			newLasers = append(newLasers, laser{x: l.x, y: l.y - 1, direction: l.direction})
		case 'd':
			newLasers = append(newLasers, laser{x: l.x, y: l.y + 1, direction: l.direction})
		case 'l', 'r':
			p := point{x: l.x, y: l.y}
			if !encounteredSplitters[p] {
				newLasers = append(
					newLasers,
					laser{x: l.x, y: l.y - 1, direction: 'u'},
					laser{x: l.x, y: l.y + 1, direction: 'd'},
				)
				encounteredSplitters[p] = true
			}

		}
	}

	return filterValidLasers(newLasers, len(mirrors[0]), len(mirrors))
}

func part1(start laser, mirrors [][]rune) int {
	pointSet := map[point]bool{}
	encounteredSplitters := map[point]bool{}
	laserQueue := Queue{start}
	for len(laserQueue) != 0 {
		l := laserQueue.dequeue()
		p := point{x: l.x, y: l.y}
		pointSet[p] = true
		laserQueue.enqueueSlice(nextLasers(l, mirrors, encounteredSplitters))
	}

	return len(pointSet)
}

func findEdgeLasers(maxX, maxY int) []laser {
	lasers := []laser{}
	for x := 0; x < maxX; x++ {
		lasers = append(lasers, laser{x: x, y: 0, direction: 'd'})
	}
	for x := 0; x < maxX; x++ {
		lasers = append(lasers, laser{x: x, y: maxY - 1, direction: 'u'})
	}
	for y := 0; y < maxY; y++ {
		lasers = append(lasers, laser{x: 0, y: y, direction: 'r'})
	}
	for y := 0; y < maxY; y++ {
		lasers = append(lasers, laser{x: maxX - 1, y: y, direction: 'l'})
	}
	return lasers
}

func part2(mirrors [][]rune) int {
	edgeLasers := findEdgeLasers(len(mirrors[0]), len(mirrors))
	scores := []int{}

	for _, l := range edgeLasers {
		scores = append(scores, part1(l, mirrors))
	}

	return slices.Max(scores)
}

func main() {
	start := time.Now()
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	mirrors := parseMirrors(scanner)

	result1 := part1(laser{x: 0, y: 0, direction: 'r'}, mirrors)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(mirrors)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
