package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type spring struct {
	layout     []rune
	conditions []int
}

func parseNums(s string) []int {
	nums := []int{}
	split := strings.Split(s, ",")
	for _, sNum := range split {
		num, _ := strconv.Atoi(sNum)
		nums = append(nums, num)
	}
	return nums
}

func parseSprings(scanner *bufio.Scanner) []spring {
	springs := []spring{}
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		springs = append(
			springs,
			spring{layout: []rune(split[0]), conditions: parseNums(split[1])},
		)
	}
	return springs
}

func conditionMet(layout []rune, conditions []int, cache map[string]int) int {
	if len(conditions) == 0 {
		return 0
	}

	if len(layout) < conditions[0] {
		return 0
	}
	for i := 0; i < conditions[0]; i++ {
		if layout[i] == '.' {
			return 0
		}
	}
	if len(layout) == conditions[0] {
		if len(conditions) == 1 {
			return 1
		}
		return 0
	}
	if layout[conditions[0]] == '#' {
		return 0
	}
	return calcSpringArrangements(layout[conditions[0]+1:], conditions[1:], cache)
}

// had to rework this for part2 with help from
// https://www.reddit.com/r/adventofcode/comments/18hbbxe/2023_day_12python_stepbystep_tutorial_with_bonus/
func calcSpringArrangements(layout []rune, conditions []int, cache map[string]int) int {
	if val, ok := cache[string(layout)+fmt.Sprint(conditions)]; ok {
		return val
	}
	if len(layout) == 0 {
		if len(conditions) == 0 {
			return 1
		}
		return 0
	}

	switch layout[0] {
	case '.':
		result := calcSpringArrangements(layout[1:], conditions, cache)
		cache[string(layout[1:])+fmt.Sprint(conditions)] = result
		return result
	case '#':
		return conditionMet(layout, conditions, cache)
	case '?':
		result := calcSpringArrangements(layout[1:], conditions, cache)
		cache[string(layout[1:])+fmt.Sprint(conditions)] = result
		return result + conditionMet(layout, conditions, cache)
	default:
		panic("unreachable")
	}
}

func part1(springs []spring) int {
	total := 0
	cache := map[string]int{}
	for _, spring := range springs {
		total += calcSpringArrangements(spring.layout, spring.conditions, cache)
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
	springs := parseSprings(scanner)

	result1 := part1(springs)
	fmt.Println("Part 1 result:", result1)

	for i := range springs {
		layout := slices.Clone(springs[i].layout)
		conditions := slices.Clone(springs[i].conditions)
		for j := 0; j < 4; j++ {
			layout = append(layout, '?')
			layout = append(layout, springs[i].layout...)
			conditions = append(conditions, springs[i].conditions...)
		}
		springs[i].layout = layout
		springs[i].conditions = conditions
	}
	result2 := part1(springs)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
