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

type state struct {
	curr point
	prev point
}

func parsePipes(scanner *bufio.Scanner) (point, [][]rune) {
	pipes := [][]rune{}
	start := point{}
	y := 0
	for scanner.Scan() {
		pipe := []rune{}
		for x, s := range scanner.Text() {
			pipe = append(pipe, s)
			if s == 'S' {
				start.x = x
				start.y = y
			}
		}
		pipes = append(pipes, pipe)
		y++
	}
	return start, pipes
}

func nextPoint(prevP point, p point, nextPipe rune) point {
	switch nextPipe {
	case '|':
		if p.y-prevP.y < 0 {
			return point{x: p.x, y: p.y - 1}
		}
		return point{x: p.x, y: p.y + 1}
	case '-':
		if p.x-prevP.x > 0 {
			return point{x: p.x + 1, y: p.y}
		}
		return point{x: p.x - 1, y: p.y}
	case 'L':
		if p.x-prevP.x == 0 {
			return point{x: p.x + 1, y: p.y}
		}
		return point{x: p.x, y: p.y - 1}
	case 'J':
		if p.x-prevP.x == 0 {
			return point{x: p.x - 1, y: p.y}
		}
		return point{x: p.x, y: p.y - 1}
	case '7':
		if p.x-prevP.x == 0 {
			return point{x: p.x - 1, y: p.y}
		}
		return point{x: p.x, y: p.y + 1}
	case 'F':
		if p.x-prevP.x == 0 {
			return point{x: p.x + 1, y: p.y}
		}
		return point{x: p.x, y: p.y + 1}
	default:
		return p
	}
}

func createStartPath(start point, pipes [][]rune) state {
	northPossible := []rune{'|', 'F', '7'}
	southPossible := []rune{'|', 'L', 'J'}
	westPossible := []rune{'-', 'L', 'F'}
	eastPossible := []rune{'-', '7', 'J'}
	s := state{}
	// north
	if start.y-1 >= 0 && slices.Contains[[]rune](northPossible, pipes[start.y-1][start.x]) {
		s = state{prev: start, curr: point{x: start.x, y: start.y - 1}}
	}
	// south
	if start.y+1 < len(pipes) && slices.Contains[[]rune](southPossible, pipes[start.y+1][start.x]) {
		s = state{prev: start, curr: point{x: start.x, y: start.y + 1}}
	}
	// west
	if start.x-1 >= 0 && slices.Contains[[]rune](westPossible, pipes[start.y][start.x-1]) {
		s = state{prev: start, curr: point{x: start.x - 1, y: start.y}}
	}
	// east
	if start.x-1 < len(pipes) && slices.Contains[[]rune](eastPossible, pipes[start.y][start.x+1]) {
		s = state{prev: start, curr: point{x: start.x + 1, y: start.y}}
	}
	return s
}

func part1(start point, pipes [][]rune) (int, map[point]bool) {
	steps := 1
	state := createStartPath(start, pipes)
	path := map[point]bool{start: true}

	for pipes[state.curr.y][state.curr.x] != 'S' {
		path[state.curr] = true
		nextPoint := nextPoint(state.prev, state.curr, pipes[state.curr.y][state.curr.x])
		state.prev = state.curr
		state.curr = nextPoint
		steps++
	}
	return steps / 2, path
}

// uses ray casting to determine if a tile is inside (https://en.wikipedia.org/wiki/Point_in_polygon)
// (this works with my input because S would be an F, if S was L or J
// then you would use F and 7 or just replace S if you're not as lazy as me)
func containedMask(path map[point]bool, pipes [][]rune) [][]bool {
	xPipesMask := make([][]bool, len(pipes))
	for i := range xPipesMask {
		xPipesMask[i] = make([]bool, len(pipes[0]))
	}
	for y := range pipes {
		withinPipes := false
		for x := range pipes[y] {
			if path[point{x: x, y: y}] {
				if pipes[y][x] == '|' || pipes[y][x] == 'L' || pipes[y][x] == 'J' {
					withinPipes = !withinPipes
				}
			} else {
				xPipesMask[y][x] = withinPipes
			}
		}
	}
	return xPipesMask
}

func part2(path map[point]bool, pipes [][]rune) int {
	containedMask := containedMask(path, pipes)

	total := 0
	for y := range containedMask {
		for x := range containedMask[0] {
			if containedMask[y][x] {
				total++
			}
		}
	}

	return total
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
	startP, pipes := parsePipes(scanner)

	result1, path := part1(startP, pipes)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(path, pipes)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
