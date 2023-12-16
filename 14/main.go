package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

const CYCLES = 1000000000

type point struct {
	x int
	y int
}

func copyMap(m map[point]bool) map[point]bool {
	copy := map[point]bool{}
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

func parseRocks(scanner *bufio.Scanner) ([]point, map[point]bool, int, int) {
	movableRocks := []point{}
	stationaryRocks := map[point]bool{}
	x := 0
	y := 0
	for scanner.Scan() {
		for x, c := range scanner.Text() {
			if c == '#' {
				stationaryRocks[point{x: x, y: y}] = true
			} else if c == 'O' {
				movableRocks = append(movableRocks, point{x: x, y: y})
			}
		}
		x = len(scanner.Text())
		y++
	}
	return movableRocks, stationaryRocks, x, y
}

func moveRockNorth(rock point, sRocks map[point]bool) point {
	for rock.y > 0 && !sRocks[point{x: rock.x, y: rock.y - 1}] {
		rock.y -= 1
	}
	return rock
}

func moveRockSouth(rock point, sRocks map[point]bool, maxY int) point {
	for rock.y < maxY-1 && !sRocks[point{x: rock.x, y: rock.y + 1}] {
		rock.y += 1
	}
	return rock
}

func moveRockWest(rock point, sRocks map[point]bool) point {
	for rock.x > 0 && !sRocks[point{x: rock.x - 1, y: rock.y}] {
		rock.x -= 1
	}
	return rock
}

func moveRockEast(rock point, sRocks map[point]bool, maxX int) point {
	for rock.x < maxX-1 && !sRocks[point{x: rock.x + 1, y: rock.y}] {
		rock.x += 1
	}
	return rock
}

func part1(mRocks []point, sRocks map[point]bool, maxY int) int {
	mRocks = moveRocks(mRocks, sRocks, sortNorth, moveRockNorth)

	total := 0
	for _, rock := range mRocks {
		total += maxY - rock.y
	}
	return total
}

func moveRocks(
	mRocks []point,
	sRocks map[point]bool,
	sort func(a, b point) int,
	move func(point, map[point]bool) point,
) []point {
	sRocksCopy := copyMap(sRocks)
	slices.SortFunc[[]point](mRocks, sort)
	for i := range mRocks {
		mRocks[i] = move(mRocks[i], sRocksCopy)
		sRocksCopy[mRocks[i]] = true
	}
	return mRocks
}

func sortNorth(a, b point) int { return a.y - b.y }
func sortWest(a, b point) int  { return a.x - b.x }
func sortSouth(a, b point) int { return b.y - a.y }
func sortEast(a, b point) int  { return b.x - a.x }

func performCycle(mRocks []point, sRocks map[point]bool, maxX int, maxY int) []point {
	mRocks = moveRocks(mRocks, sRocks, sortNorth, moveRockNorth)
	mRocks = moveRocks(mRocks, sRocks, sortWest, moveRockWest)
	mRocks = moveRocks(mRocks, sRocks, sortSouth, func(p point, m map[point]bool) point { return moveRockSouth(p, m, maxY) })
	mRocks = moveRocks(mRocks, sRocks, sortEast, func(p point, m map[point]bool) point { return moveRockEast(p, m, maxX) })
	return mRocks
}

func printRocks(mRocks []point, sRocks map[point]bool, maxX, maxY int) {
	res := make([][]rune, maxY)
	for y := range res {
		for x := 0; x < maxX; x++ {
			res[y] = append(res[y], '.')
		}
	}
	for _, rock := range mRocks {
		res[rock.y][rock.x] = 'O'
	}
	for rock := range sRocks {
		res[rock.y][rock.x] = '#'
	}
	for _, l := range res {
		fmt.Println(string(l))
	}
}

func part2(mRocks []point, sRocks map[point]bool, maxX int, maxY int) int {
	for i := 0; i < 3; i++ {
		mRocks = performCycle(mRocks, sRocks, maxX, maxY)
	}
	// printRocks(mRocks, sRocks, maxX, maxY)
	total := 0
	for _, rock := range mRocks {
		total += maxY - rock.y
	}
	return total
}

func main() {
	start := time.Now()
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	movableRocks, stationaryRocks, x, y := parseRocks(scanner)

	result1 := part1(movableRocks, stationaryRocks, y)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(movableRocks, stationaryRocks, x, y)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
