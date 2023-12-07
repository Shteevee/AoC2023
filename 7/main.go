package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	fiveOfAKind  = 7
	fourOfAKind  = 6
	fullHouse    = 5
	threeOfAKind = 4
	twoPair      = 3
	onePair      = 2
	highCard     = 1
)

type occurrenceMap = map[rune]int

type round struct {
	cardOccurrences occurrenceMap
	hand            string
	handStrength    int
	bid             int
}

func parseRounds(scanner *bufio.Scanner) []round {
	rounds := make([]round, 0)
	for scanner.Scan() {
		splitRound := strings.Split(scanner.Text(), " ")
		cardOccurrences := occurrenceMap{}
		for _, card := range splitRound[0] {
			cardOccurrences[card]++
		}
		bid, _ := strconv.Atoi(splitRound[1])
		rounds = append(rounds, round{hand: splitRound[0], cardOccurrences: cardOccurrences, bid: bid})
	}
	return rounds
}

func cardValuesMap() map[rune]int {
	return map[rune]int{
		'A': 13, 'K': 12, 'Q': 11, 'J': 10, 'T': 9, '9': 8,
		'8': 7, '7': 6, '6': 5, '5': 4, '4': 3, '3': 2, '2': 1,
	}
}

// assumes list is sorted
func handStrength(occurrenceList []int) int {
	switch {
	case occurrenceList[len(occurrenceList)-1] == 5:
		return fiveOfAKind
	case occurrenceList[len(occurrenceList)-1] == 4:
		return fourOfAKind
	case occurrenceList[len(occurrenceList)-1] == 3 && occurrenceList[len(occurrenceList)-2] == 2:
		return fullHouse
	case occurrenceList[len(occurrenceList)-1] == 3:
		return threeOfAKind
	case occurrenceList[len(occurrenceList)-1] == 2 && occurrenceList[len(occurrenceList)-2] == 2:
		return twoPair
	case occurrenceList[len(occurrenceList)-1] == 2:
		return onePair
	default:
		return highCard
	}
}

func calcHandStrength(cardOccurrences occurrenceMap) int {
	occurrenceList := []int{}
	for _, co := range cardOccurrences {
		occurrenceList = append(occurrenceList, co)
	}
	slices.Sort(occurrenceList)
	return handStrength(occurrenceList)
}

func calcHandStrength2(cardOccurrences occurrenceMap) int {
	occurrenceList := []int{}
	for c, co := range cardOccurrences {
		if c != 'J' {
			occurrenceList = append(occurrenceList, co)
		}
	}
	// all jokers case
	if len(occurrenceList) == 0 {
		return fiveOfAKind
	}
	slices.Sort(occurrenceList)
	occurrenceList[len(occurrenceList)-1] += cardOccurrences['J']
	return handStrength(occurrenceList)
}

func compare(r1 round, r2 round, cardValues map[rune]int) int {
	score := cmp.Compare(r1.handStrength, r2.handStrength)
	if score == 0 {
		for i := range r1.hand {
			cardScore := cardValues[rune(r1.hand[i])] - cardValues[rune(r2.hand[i])]
			if cardScore != 0 {
				score = cardScore
				break
			}
		}
	}
	return score
}

func part1(rounds []round) int {
	for i, round := range rounds {
		rounds[i].handStrength = calcHandStrength(round.cardOccurrences)
	}
	cardValues := cardValuesMap()
	slices.SortFunc(rounds, func(r1, r2 round) int {
		return compare(r1, r2, cardValues)
	})

	total := 0
	for i := range rounds {
		total += rounds[i].bid * (i + 1)
	}
	return total
}

func part2(rounds []round) int {
	for i, round := range rounds {
		rounds[i].handStrength = calcHandStrength2(round.cardOccurrences)
	}
	cardValues := cardValuesMap()
	cardValues['J'] = 0
	slices.SortFunc(rounds, func(r1, r2 round) int {
		return compare(r1, r2, cardValues)
	})

	total := 0
	for i := range rounds {
		total += rounds[i].bid * (i + 1)
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
	rounds := parseRounds(scanner)

	result1 := part1(rounds)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(rounds)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
