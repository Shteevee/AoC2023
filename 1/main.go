package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func createNumberTextMapping() map[string]rune {
	return map[string]rune{
		"one":   '1',
		"two":   '2',
		"three": '3',
		"four":  '4',
		"five":  '5',
		"six":   '6',
		"seven": '7',
		"eight": '8',
		"nine":  '9',
	}
}

func parseText(scanner *bufio.Scanner) []string {
	text := make([]string, 0)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text
}

func findFirstAndLastNumericChar(line string) string {
	var first rune
	var last rune
	for _, char := range line {
		if unicode.IsDigit(char) {
			if first == 0 {
				first = char
			}
			last = char
		}
	}
	return string(first) + string(last)
}

func sumFirstAndLastDigits(text []string) int {
	total := 0
	for _, line := range text {
		firstAndLastString := findFirstAndLastNumericChar(line)
		firstAndLast, _ := strconv.Atoi(firstAndLastString)
		total += int(firstAndLast)
	}

	return total
}

func findFirstAndLastNumbers(line string, numberMap map[string]rune) string {
	var first rune
	var last rune
	firstIndex := -1
	lastIndex := -1

	for i, char := range line {
		if unicode.IsDigit(char) {
			if firstIndex == -1 {
				firstIndex = i
				first = char
			}
			lastIndex = i
			last = char
		}
	}

	if firstIndex == -1 {
		firstIndex = len(line) - 1
	}

	for number := range numberMap {
		i := strings.Index(line, number)
		j := strings.LastIndex(line, number)
		if i != -1 && (i < firstIndex) {
			firstIndex = i
			first = numberMap[number]
		}
		if j != -1 && j > lastIndex {
			lastIndex = j
			last = numberMap[number]
		}
	}

	return string(first) + string(last)
}

func sumFirstAndLastNumbers(text []string) int {
	total := 0
	numbersMap := createNumberTextMapping()
	for _, line := range text {
		firstAndLastString := findFirstAndLastNumbers(line, numbersMap)
		firstAndLast, _ := strconv.Atoi(firstAndLastString)
		total = total + firstAndLast
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
	text := parseText(scanner)

	result1 := sumFirstAndLastDigits(text)
	fmt.Println("Part 1 result:", result1)

	result2 := sumFirstAndLastNumbers(text)
	fmt.Println("Part 2 result:", result2)
	log.Printf("Time taken: %s", time.Since(start))
}
