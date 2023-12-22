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

const (
	REMOVE = 0
	EQUAL  = 1
)

type instruction struct {
	label     string
	focalLen  int
	hash      int
	intsrType int
}

func parseInstr(scanner *bufio.Scanner) []string {
	scanner.Scan()
	return strings.Split(scanner.Text(), ",")
}

func HASH(s string) int {
	current := 0
	for _, c := range s {
		current = ((current + int(c)) * 17) % 256
	}
	return current
}

func part1(instr []string) int {
	total := 0
	for _, s := range instr {
		total += HASH(s)
	}
	return total
}

func createInstructions(instr []string) []instruction {
	instructions := []instruction{}
	for _, s := range instr {
		instruction := instruction{}
		if split := strings.Split(s, "="); len(split) == 2 {
			instruction.label = split[0]
			instruction.hash = HASH(split[0])
			instruction.focalLen, _ = strconv.Atoi(split[1])
			instruction.intsrType = EQUAL
		} else {
			split := strings.Split(s, "-")
			instruction.label = split[0]
			instruction.hash = HASH(split[0])
			instruction.intsrType = REMOVE
		}
		instructions = append(instructions, instruction)
	}

	return instructions
}

func fillBoxes(instructions []instruction) [][]instruction {
	boxes := make([][]instruction, 256)
	for _, instrct := range instructions {
		switch instrct.intsrType {
		case EQUAL:
			replaceIndex := slices.IndexFunc[[]instruction](
				boxes[instrct.hash],
				func(i instruction) bool { return i.label == instrct.label },
			)
			if replaceIndex == -1 {
				boxes[instrct.hash] = append(boxes[instrct.hash], instrct)
			} else {
				boxes[instrct.hash][replaceIndex] = instrct
			}
		case REMOVE:
			removeIndex := slices.IndexFunc[[]instruction](
				boxes[instrct.hash],
				func(i instruction) bool { return i.label == instrct.label },
			)
			if removeIndex != -1 {
				boxes[instrct.hash] = append(
					boxes[instrct.hash][:removeIndex],
					boxes[instrct.hash][removeIndex+1:]...,
				)
			}
		}
	}
	return boxes
}

func part2(instr []string) int {
	instructions := createInstructions(instr)
	boxes := fillBoxes(instructions)

	total := 0
	for boxNum, box := range boxes {
		for lensNum, lens := range box {
			total += (boxNum + 1) * (lensNum + 1) * lens.focalLen
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
	instr := parseInstr(scanner)

	result1 := part1(instr)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(instr)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
