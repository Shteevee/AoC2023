package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type pattern = [][]rune

func parsePatterns(scanner *bufio.Scanner) []pattern {
	patterns := []pattern{}
	pattern := [][]rune{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 {
			pattern = append(pattern, []rune(line))
		} else {
			patterns = append(patterns, pattern)
			pattern = [][]rune{}
		}
	}
	patterns = append(patterns, pattern)
	return patterns
}

func transpose[T any](a [][]T) [][]T {
	t := make([][]T, len(a[0]))
	for i := range t {
		t[i] = make([]T, len(a))
	}
	for x := range a[0] {
		for y := range a {
			t[x][y] = a[y][x]
		}
	}
	return t
}

func findReflectionIndex(pattern pattern, targetDiff int) int {
	for i := 1; i < len(pattern); i++ {
		diffs := 0
		for j := 0; j < i && j < len(pattern)-i; j++ {
			for k := 0; k < len(pattern[i-j-1]); k++ {
				if pattern[i-j-1][k] != pattern[i+j][k] {
					diffs++
				}
			}
		}
		if diffs == targetDiff {
			return i
		}
	}
	return 0
}

func findReflectionValue(pattern pattern, targetDiff int) int {
	yReflection := findReflectionIndex(pattern, targetDiff)
	if yReflection != 0 {
		return yReflection * 100
	}
	return findReflectionIndex(transpose(pattern), targetDiff)
}

func part1(patterns []pattern) int {
	total := 0
	for _, pattern := range patterns {
		total += findReflectionValue(pattern, 0)
	}
	return total
}

func part2(patterns []pattern) int {
	total := 0
	for _, pattern := range patterns {
		total += findReflectionValue(pattern, 1)
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
	patterns := parsePatterns(scanner)

	result1 := part1(patterns)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(patterns)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
