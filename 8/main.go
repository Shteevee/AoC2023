package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type maps = map[string]map[rune]string

func parseMaps(scanner *bufio.Scanner) (string, maps) {
	scanner.Scan()
	directions := scanner.Text()
	scanner.Scan()

	maps := maps{}
	for scanner.Scan() {
		splitMap := strings.Split(scanner.Text(), " = ")
		dirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(splitMap[1], "("), ")"), ", ")
		aMap := map[rune]string{}
		aMap['L'] = dirs[0]
		aMap['R'] = dirs[1]
		maps[splitMap[0]] = aMap
	}
	return directions, maps
}

func part1(directions string, maps maps) int {
	steps := 0
	nextLocation := "AAA"
	for nextLocation != "ZZZ" {
		nextLocation = maps[nextLocation][rune(directions[steps%len(directions)])]
		steps++
	}
	return steps
}

func findLocationsEndingWithA(maps maps) []string {
	locations := []string{}
	for loc := range maps {
		if loc[len(loc)-1] == 'A' {
			locations = append(locations, loc)
		}
	}
	return locations
}

func allCyclesFilled(cycles []int) bool {
	result := true
	for _, cycle := range cycles {
		result = result && cycle != 0
	}
	return result
}

func gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func findLCM(xs []int) int {
	result := xs[0]
	for i := 1; i < len(xs); i++ {
		result = (xs[i] * result) / (gcd(xs[i], result))
	}

	return result
}

func part2(directions string, maps maps) int {
	steps := 0
	nextLocations := findLocationsEndingWithA(maps)
	cycles := make([]int, len(nextLocations))
	for !allCyclesFilled(cycles) {
		for i := range nextLocations {
			if nextLocations[i][len(nextLocations[i])-1] == 'Z' {
				cycles[i] = steps
			}
			nextLocations[i] = maps[nextLocations[i]][rune(directions[steps%len(directions)])]
		}
		steps++
	}
	return findLCM(cycles)
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
	directions, maps := parseMaps(scanner)

	result1 := part1(directions, maps)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(directions, maps)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
