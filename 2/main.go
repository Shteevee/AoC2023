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

const maxRed = 12
const maxBlue = 14
const maxGreen = 13

type CubeSet = map[string]int

type Game struct {
	id       int
	cubeSets []CubeSet
}

func stripGameId(line string) (int, string) {
	line = strings.TrimPrefix(line, "Game ")
	splitLine := strings.Split(line, ": ")
	id, _ := strconv.Atoi(splitLine[0])
	return id, splitLine[1]
}

func parseCubeSets(line string) []CubeSet {
	cubeSets := make([]CubeSet, 0)
	for _, set := range strings.Split(line, "; ") {
		cubeSet := make(map[string]int)
		for _, cubeCount := range strings.Split(set, ", ") {
			splitCount := strings.Split(cubeCount, " ")
			count, _ := strconv.Atoi(splitCount[0])
			cubeSet[splitCount[1]] = count
		}
		cubeSets = append(cubeSets, cubeSet)
	}
	return cubeSets
}

func parseGames(scanner *bufio.Scanner) []Game {
	games := make([]Game, 0)
	for scanner.Scan() {
		line := scanner.Text()
		game := Game{}
		game.id, line = stripGameId(line)
		game.cubeSets = parseCubeSets(line)
		games = append(games, game)
	}
	return games
}

func gameIsPossible(game Game) bool {
	possible := true
	for _, set := range game.cubeSets {
		possible = possible && set["red"] <= maxRed && set["green"] <= maxGreen && set["blue"] <= maxBlue
	}
	return possible
}

func sumPossibleGameIds(games []Game) int {
	total := 0
	for _, game := range games {
		if gameIsPossible(game) {
			total += game.id
		}
	}
	return total
}

func calcMinCubeSetPower(game Game) int {
	maxSet := make(map[string]int)
	for _, set := range game.cubeSets {
		for k, q := range set {
			if q > maxSet[k] {
				maxSet[k] = q
			}
		}
	}
	return maxSet["red"] * maxSet["blue"] * maxSet["green"]
}

func sumMinCubeSetPower(games []Game) int {
	total := 0
	for _, game := range games {
		total += calcMinCubeSetPower(game)
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
	games := parseGames(scanner)

	result1 := sumPossibleGameIds(games)
	fmt.Println("Part 1 result:", result1)

	result2 := sumMinCubeSetPower(games)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
