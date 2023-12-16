package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type point struct {
	x int
	y int
}

func parseRocks(scanner *bufio.Scanner) ([]point, map[point]bool, int) {
	movableRocks := []point{}
	stationaryRocks := map[point]bool{}
	y := 0
	for scanner.Scan() {
		for x, c := range scanner.Text() {
			if c == '#' {
				stationaryRocks[point{x: x, y: y}] = true
			} else if c == 'O' {
				movableRocks = append(movableRocks, point{x: x, y: y})
			}
		}
		y++
	}
	return movableRocks, stationaryRocks, y
}

func moveRockNorth(rock point, sRocks map[point]bool) point {
	for rock.y > 0 && !sRocks[point{x: rock.x, y: rock.y - 1}] {
		rock.y -= 1
	}
	return rock
}

func part1(mRocks []point, sRocks map[point]bool, maxHeight int) int {
	for i := range mRocks {
		mRocks[i] = moveRockNorth(mRocks[i], sRocks)
		sRocks[mRocks[i]] = true
	}

	total := 0
	for _, rock := range mRocks {
		total += maxHeight - rock.y
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
	movableRocks, stationaryRocks, maxHeight := parseRocks(scanner)

	result1 := part1(movableRocks, stationaryRocks, maxHeight)
	fmt.Println("Part 1 result:", result1)

	// result2 := part2(patterns)
	// fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
