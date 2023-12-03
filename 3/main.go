package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

type Pos struct {
	x int
	y int
}

type EnginePart struct {
	number   int
	startPos Pos
	lastPos  Pos
}

type Symbol struct {
	value rune
	pos   Pos
}

func createEnginePart(partNumber string, startIndex int, lineNum int) EnginePart {
	number, _ := strconv.Atoi(partNumber)
	return EnginePart{
		number:   number,
		startPos: Pos{x: startIndex, y: lineNum},
		lastPos:  Pos{x: startIndex + len(partNumber) - 1, y: lineNum},
	}
}

func parseSchematicLine(lineNum int, line string) ([]EnginePart, []Symbol) {
	engineParts := make([]EnginePart, 0)
	symbols := make([]Symbol, 0)
	collector := ""
	collectorStartIndex := -1
	for i, c := range line {
		if c != '.' && !unicode.IsDigit(c) {
			symbols = append(symbols, Symbol{value: c, pos: Pos{x: i, y: lineNum}})
		}
		if unicode.IsDigit(c) {
			if collectorStartIndex == -1 {
				collectorStartIndex = i
			}
			collector += string(c)
		}
		if !unicode.IsDigit(c) && len(collector) > 0 {
			engineParts = append(engineParts, createEnginePart(collector, collectorStartIndex, lineNum))
			collector = ""
			collectorStartIndex = -1
		}
	}

	if len(collector) > 0 {
		engineParts = append(engineParts, createEnginePart(collector, collectorStartIndex, lineNum))
	}

	return engineParts, symbols
}

func parseSchematic(scanner *bufio.Scanner) ([]EnginePart, []Symbol) {
	engineParts := make([]EnginePart, 0)
	symbols := make([]Symbol, 0)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineEngineParts, lineSymbols := parseSchematicLine(lineNum, line)
		engineParts = append(engineParts, lineEngineParts...)
		symbols = append(symbols, lineSymbols...)
		lineNum++
	}
	return engineParts, symbols
}

func isAdjacent(symbol Symbol, enginePart EnginePart) bool {
	return (symbol.pos.x == enginePart.startPos.x-1 && symbol.pos.y == enginePart.startPos.y) ||
		(symbol.pos.x == enginePart.lastPos.x+1 && symbol.pos.y == enginePart.startPos.y) ||
		(symbol.pos.x >= enginePart.startPos.x-1 && symbol.pos.x <= enginePart.lastPos.x+1 && symbol.pos.y == enginePart.startPos.y+1) ||
		(symbol.pos.x >= enginePart.startPos.x-1 && symbol.pos.x <= enginePart.lastPos.x+1 && symbol.pos.y == enginePart.startPos.y-1)
}

func sumNeededEngineParts(engineParts []EnginePart, symbols []Symbol) int {
	neededEngineParts := map[EnginePart]struct{}{}
	for _, symbol := range symbols {
		for _, enginePart := range engineParts {
			if isAdjacent(symbol, enginePart) {
				neededEngineParts[enginePart] = struct{}{}
			}
		}
	}

	total := 0
	for k := range neededEngineParts {
		total += k.number
	}

	return total
}

func findGears(symbols []Symbol) []Symbol {
	gears := make([]Symbol, 0)
	for _, symbol := range symbols {
		if symbol.value == '*' {
			gears = append(gears, symbol)
		}
	}
	return gears
}

func sumGearRatios(engineParts []EnginePart, symbols []Symbol) int {
	gears := findGears(symbols)
	gearEngineParts := make(map[Symbol][]EnginePart)
	for _, gear := range gears {
		for _, enginePart := range engineParts {
			if isAdjacent(gear, enginePart) {
				gearEngineParts[gear] = append(gearEngineParts[gear], enginePart)
			}
		}
	}

	total := 0
	for _, engineParts := range gearEngineParts {
		if len(engineParts) == 2 {
			total += engineParts[0].number * engineParts[1].number
		}
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
	engineParts, symbols := parseSchematic(scanner)

	result1 := sumNeededEngineParts(engineParts, symbols)
	fmt.Println("Part 1 result:", result1)

	result2 := sumGearRatios(engineParts, symbols)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
