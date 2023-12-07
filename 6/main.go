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

type race struct {
	time       int
	distRecord int
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

func parseKerning(s string) int {
	s = strings.Replace(s, " ", "", -1)
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func parseRaces(scanner *bufio.Scanner) ([]race, race) {
	scanner.Scan()
	times := parseNumList(strings.TrimPrefix(scanner.Text(), "Time:"))
	bigTime := parseKerning(strings.TrimPrefix(scanner.Text(), "Time:"))
	scanner.Scan()
	distances := parseNumList(strings.TrimPrefix(scanner.Text(), "Distance:"))
	bigDistance := parseKerning(strings.TrimPrefix(scanner.Text(), "Distance:"))
	races := make([]race, 0)
	for i := range times {
		races = append(races, race{time: times[i], distRecord: distances[i]})
	}
	return races, race{time: bigTime, distRecord: bigDistance}
}

func calcRecordBreaksMag(r race) int {
	breaks := 0
	for i := 1; i < r.time; i++ {
		dist := i * (r.time - i)
		if dist > r.distRecord {
			breaks++
		}
	}
	return breaks
}

func part1(races []race) int {
	total := 1
	for _, race := range races {
		recordBreaks := calcRecordBreaksMag(race)
		total = total * recordBreaks
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
	races, bigRace := parseRaces(scanner)
	result1 := part1(races)
	fmt.Println("Part 1 result:", result1)

	result2 := part1([]race{bigRace})
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
