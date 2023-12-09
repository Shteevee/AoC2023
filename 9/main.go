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

func parseOasisHistories(scanner *bufio.Scanner) [][]int {
	histories := [][]int{}
	for scanner.Scan() {
		history := []int{}
		splitLine := strings.Split(scanner.Text(), " ")
		for _, s := range splitLine {
			num, _ := strconv.Atoi(s)
			history = append(history, num)
		}
		histories = append(histories, history)
	}
	return histories
}

func allZeros(xs []int) bool {
	res := true
	for _, x := range xs {
		res = res && x == 0
	}
	return res
}

func nextNumInSeq(hist []int) int {
	result := hist[len(hist)-1]
	nextHist := []int{}
	for !allZeros(hist) {
		for i := 0; i < len(hist)-1; i++ {
			nextHist = append(nextHist, hist[i+1]-hist[i])
		}
		result += nextHist[len(nextHist)-1]
		hist = nextHist
		nextHist = []int{}
	}
	return result
}

func part1(histories [][]int) int {
	total := 0
	for _, hist := range histories {
		total += nextNumInSeq(hist)
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
	histories := parseOasisHistories(scanner)

	result1 := part1(histories)
	fmt.Println("Part 1 result:", result1)

	for i := range histories {
		slices.Reverse(histories[i])
	}
	result2 := part1(histories)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
