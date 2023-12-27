package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type plan struct {
	dir    rune
	dist   int
	colour string
}

func parseDigPlans(scanner *bufio.Scanner) []plan {
	plans := []plan{}
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		dist, _ := strconv.Atoi(split[1])
		colour := strings.TrimSuffix(strings.TrimPrefix(split[2], "("), ")")
		plans = append(plans, plan{dir: rune(split[0][0]), dist: dist, colour: colour})
	}
	return plans
}

func nextPoint(p point, plan plan) point {
	dx, dy := 0, 0
	switch plan.dir {
	case 'R':
		dx = 1
	case 'L':
		dx = -1
	case 'U':
		dy = -1
	case 'D':
		dy = 1
	}
	for i := 1; i <= plan.dist; i++ {
		p.x += dx
		p.y += dy
	}
	return p
}

func getVertices(plans []plan) []point {
	vertices := []point{{0, 0}}
	lastPos := point{0, 0}
	for i := 0; i < len(plans)-1; i++ {
		lastPos = nextPoint(lastPos, plans[i])
		vertices = append(vertices, lastPos)
	}
	return vertices
}

// shoelace formula
func calcArea(vertices []point, perimeter int) int {
	a, b := 0, 0
	for i := range vertices {
		if i == len(vertices)-1 {
			a += vertices[i].x * vertices[0].y
			b += vertices[i].y * vertices[0].x
		} else {
			a += vertices[i].x * vertices[i+1].y
			b += vertices[i].y * vertices[i+1].x
		}
	}

	return (a-b)/2 + (perimeter / 2) + 1
}

func part1(plans []plan) int {
	vertices := getVertices(plans)
	perimeter := 0
	for _, plan := range plans {
		perimeter += plan.dist
	}
	return calcArea(vertices, perimeter)
}

func dirFromColour(s string) rune {
	var res rune
	switch s[len(s)-1] {
	case '0':
		res = 'R'
	case '1':
		res = 'D'
	case '2':
		res = 'L'
	case '3':
		res = 'U'
	}
	return res
}

func findNewPlans(plans []plan) []plan {
	newPlans := []plan{}
	for _, p := range plans {
		dist, _ := strconv.ParseInt(p.colour[1:len(p.colour)-1], 16, 0)
		newPlans = append(
			newPlans,
			plan{dir: dirFromColour(p.colour), dist: int(dist)},
		)
	}
	return newPlans
}

func part2(plans []plan) int {
	newPlans := findNewPlans(plans)
	return part1(newPlans)
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
	plans := parseDigPlans(scanner)

	result1 := part1(plans)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(plans)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
