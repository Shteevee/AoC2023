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

type Card struct {
	winningNums  []int
	selectedNums []int
}

func powInt(x int, y int) int {
	total := 1
	for i := 0; i < y; i++ {
		total = total * x
	}
	return total
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

func parseCards(scanner *bufio.Scanner) []Card {
	cards := make([]Card, 0)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, ": ")[1]
		cardInfoString := strings.Split(line, " | ")
		cards = append(cards, Card{winningNums: parseNumList(cardInfoString[0]), selectedNums: parseNumList(cardInfoString[1])})
	}
	return cards
}

func calcCardWins(card Card) int {
	cardWins := 0
	for _, selectedNum := range card.selectedNums {
		for _, winningNum := range card.winningNums {
			if selectedNum == winningNum {
				cardWins++
			}
		}
	}
	return cardWins
}

func calcWinningScore(cards []Card) int {
	total := 0
	for _, card := range cards {
		cardWins := calcCardWins(card)
		if cardWins > 0 {
			total += powInt(2, cardWins-1)
		}
	}
	return total
}

func initCardOccurrences(length int) []int {
	cardOccurrences := make([]int, length)
	for i := range cardOccurrences {
		cardOccurrences[i] = 1
	}
	return cardOccurrences
}

func calcTotalCards(cards []Card) int {
	cardWinsLookup := make([]int, 0)
	for _, card := range cards {
		cardWinsLookup = append(cardWinsLookup, calcCardWins(card))
	}
	cardOccurrences := initCardOccurrences(len(cardWinsLookup))
	for i, cardWins := range cardWinsLookup {
		for j := 1; j <= cardWins; j++ {
			cardOccurrences[i+j] += cardOccurrences[i]
		}
	}
	total := 0
	for _, num := range cardOccurrences {
		total += num
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
	cards := parseCards(scanner)

	result1 := calcWinningScore(cards)
	fmt.Println("Part 1 result:", result1)

	result2 := calcTotalCards(cards)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
