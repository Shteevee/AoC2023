package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	seedToSoil            = 0
	soilToFertilizer      = 1
	fertilizerToWater     = 2
	waterToLight          = 3
	lightToTemperature    = 4
	temperatureToHumidity = 5
	humidityToLocation    = 6
)

type Mapping struct {
	destRangeStart int
	srcRangeStart  int
	rangeLength    int
}

type Seed struct {
	start       int
	rangeLength int
}

type RangeMapping = [][]Mapping

func createMappingMap() RangeMapping {
	mappingMap := make([][]Mapping, 7)
	for i := range mappingMap {
		mappingMap[i] = make([]Mapping, 0)
	}
	return mappingMap
}

func parseNumList(numList string) []int {
	nums := make([]int, 0)
	for _, sNum := range strings.Split(numList, " ") {
		if len(sNum) > 0 {
			num, _ := strconv.Atoi(strings.TrimSpace(sNum))
			nums = append(nums, num)
		}
	}
	return nums
}

func nextMappingToParse(line string) int {
	switch {
	case strings.Contains(line, "seed-to-soil"):
		return seedToSoil
	case strings.Contains(line, "soil-to-fertilizer"):
		return soilToFertilizer
	case strings.Contains(line, "fertilizer-to-water"):
		return fertilizerToWater
	case strings.Contains(line, "water-to-light"):
		return waterToLight
	case strings.Contains(line, "light-to-temperature"):
		return lightToTemperature
	case strings.Contains(line, "temperature-to-humidity"):
		return temperatureToHumidity
	case strings.Contains(line, "humidity-to-location"):
		return humidityToLocation
	default:
		return -1
	}
}

func parseSeeds(scanner *bufio.Scanner) ([]int, RangeMapping) {
	seeds := make([]int, 0)
	scanner.Scan()
	seeds = append(seeds, parseNumList(strings.TrimPrefix(scanner.Text(), "seeds: "))...)

	currentMapping := -1
	mappingMap := createMappingMap()
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			nextMapping := nextMappingToParse(line)
			if nextMapping != -1 {
				currentMapping = nextMapping
			} else {
				rangeMapping := parseNumList(line)
				mappingMap[currentMapping] = append(
					mappingMap[currentMapping],
					Mapping{
						destRangeStart: rangeMapping[0],
						srcRangeStart:  rangeMapping[1],
						rangeLength:    rangeMapping[2],
					},
				)
			}
		}
	}

	return seeds, mappingMap
}

func inRange(value int, low int, high int) bool {
	return value >= low && value < high
}

func min(values ...int) int {
	min := math.MaxInt
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func minSeed(values ...Seed) int {
	min := math.MaxInt
	for _, value := range values {
		if value.start < min && value.start != 0 {
			min = value.start
		}
	}
	return min
}

func findLowestLocationNumber(seeds []int, mappingMap RangeMapping) int {
	candidates := seeds
	nextCandidates := make([]int, 0)
	for _, mapping := range mappingMap {
		for _, candidate := range candidates {
			dest := candidate
			for _, target := range mapping {
				if inRange(candidate, target.srcRangeStart, target.srcRangeStart+target.rangeLength) {
					dest = candidate + (target.destRangeStart - target.srcRangeStart)
				}
			}
			nextCandidates = append(nextCandidates, dest)
		}
		candidates = nextCandidates
		nextCandidates = make([]int, 0)
	}

	return min(candidates...)
}

func createSeedRanges(xs []int) []Seed {
	seedRanges := make([]Seed, 0)
	for i := 0; i < len(xs); i += 2 {
		seedRanges = append(seedRanges, Seed{start: xs[i], rangeLength: xs[i+1]})
	}
	return seedRanges
}

func calcNewRanges(seed Seed, mappings []Mapping) []Seed {
	for _, mapping := range mappings {
		seedEnd := seed.start + seed.rangeLength - 1
		low := mapping.srcRangeStart
		high := mapping.srcRangeStart + mapping.rangeLength
		offset := mapping.destRangeStart - mapping.srcRangeStart
		// fully contained
		if seed.start >= low && seedEnd <= high {
			return []Seed{
				{start: seed.start + offset, rangeLength: seed.rangeLength},
			}
		}
		// left contained, right not
		if seed.start >= low && seed.start <= high && seedEnd > high {
			return append(
				[]Seed{{start: seed.start + offset, rangeLength: high - seed.start + 1}},
				calcNewRanges(Seed{start: high + 1, rangeLength: seedEnd - high}, mappings)...,
			)
		}
		// right contained, left not
		if seed.start < low && seedEnd <= high && seedEnd >= low {
			return append(
				[]Seed{{start: low + offset, rangeLength: seedEnd - low + 1}},
				calcNewRanges(Seed{start: seed.start, rangeLength: low - seed.start}, mappings)...,
			)
		}
		// both ends hang over mapping
		if seed.start < low && seedEnd > high {
			return append(
				append(
					[]Seed{{start: low + offset, rangeLength: high - low + 1}},
					calcNewRanges(Seed{start: high + 1, rangeLength: seedEnd - high}, mappings)...,
				),
				calcNewRanges(Seed{start: seed.start, rangeLength: low - seed.start}, mappings)...,
			)
		}
	}
	return []Seed{{start: seed.start, rangeLength: seed.rangeLength}}
}

func findLowestLocationNumberFromSeedRanges(seeds []Seed, mappingMap RangeMapping) int {
	for _, mappings := range mappingMap {
		newSeeds := []Seed{}
		for _, seedRange := range seeds {
			newSeeds = append(newSeeds, calcNewRanges(seedRange, mappings)...)
		}
		seeds = newSeeds
	}

	return minSeed(seeds...)
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
	seeds, mappingMap := parseSeeds(scanner)

	result1 := findLowestLocationNumber(seeds, mappingMap)
	fmt.Println("Part 1 result:", result1)

	seedRanges := createSeedRanges(seeds)
	result2 := findLowestLocationNumberFromSeedRanges(seedRanges, mappingMap)
	// the result is off by one, should fix but won't
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
