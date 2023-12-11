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

type pair struct {
	a point
	b point
}

func parseGalaxy(scanner *bufio.Scanner) [][]rune {
	galaxy := [][]rune{}
	y := 0
	for scanner.Scan() {
		row := []rune{}
		for _, s := range scanner.Text() {
			row = append(row, s)
		}
		galaxy = append(galaxy, row)
		y++
	}
	return galaxy
}

func findPlanets(galaxy [][]rune, blankSpacing int) []point {
	xOffset := 0
	xMapping := make([]int, len(galaxy[0]))
	for x := range galaxy[0] {
		isBlank := true
		for y := range galaxy {
			if galaxy[y][x] == '#' {
				isBlank = false
			}
		}
		xMapping[x] = x + xOffset
		if isBlank {
			xOffset += blankSpacing - 1
		}
	}

	planets := []point{}
	yOffset := 0
	for y := range galaxy {
		isBlank := true
		for x, s := range galaxy[y] {
			if s == '#' {
				isBlank = false
				planets = append(planets, point{x: xMapping[x], y: y + yOffset})
			}
		}
		if isBlank {
			yOffset += blankSpacing - 1
		}
	}
	return planets
}

func generatePairs(planets []point) []pair {
	pairs := []pair{}
	for i, p1 := range planets {
		for j := i + 1; j < len(planets); j++ {
			pairs = append(pairs, pair{a: p1, b: planets[j]})
		}
	}
	return pairs
}

func manhattenDist(a point, b point) int {
	dist := 0
	if a.x > b.x {
		dist += a.x - b.x
	} else {
		dist += b.x - a.x
	}

	if a.y > b.y {
		dist += a.y - b.y
	} else {
		dist += b.y - a.y
	}

	return dist
}

func part1(planets []point) int {
	pairs := generatePairs(planets)

	total := 0
	for _, pair := range pairs {
		total += manhattenDist(pair.a, pair.b)
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
	galaxy := parseGalaxy(scanner)
	planets := findPlanets(galaxy, 2)

	result1 := part1(planets)
	fmt.Println("Part 1 result:", result1)

	oldPlanets := findPlanets(galaxy, 1000000)
	result2 := part1(oldPlanets)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
