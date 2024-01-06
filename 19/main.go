// parsing could be better but solution is fine
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	COND = 0
	DIR  = 1
)

const (
	LESS = 0
	MORE = 1
)

type part = map[string]int
type rangePart = map[string]*Range

type rule struct {
	t      int
	attr   string
	cond   int
	value  int
	result string
}

type Range struct {
	lower int
	upper int
}

func parsePart(s string) part {
	part := part{}
	split := strings.Split(strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}"), ",")
	for _, v := range split {
		vSplit := strings.Split(v, "=")
		num, _ := strconv.Atoi(vSplit[1])
		part[vSplit[0]] = num
	}
	return part
}

func parseRules(s string) []rule {
	rules := []rule{}
	for _, r := range strings.Split(s, ",") {
		rule := rule{}
		if strings.Contains(r, ":") {
			re := regexp.MustCompile(`(\w+).(\w+):(\w+)`)
			matches := re.FindStringSubmatch(r)
			rule.t = COND
			rule.attr = matches[1]
			rule.value, _ = strconv.Atoi(matches[2])
			rule.result = matches[3]
			if strings.Contains(r, ">") {
				rule.cond = MORE
			} else {
				rule.cond = LESS
			}
		} else {
			rule.t = DIR
			rule.result = r
		}
		rules = append(rules, rule)
	}
	return rules
}

func parseInstr(scanner *bufio.Scanner) (map[string][]rule, []part) {
	instr := map[string][]rule{}
	parts := []part{}
	isInstr := true
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			isInstr = false
			continue
		}
		if isInstr {
			line := strings.Split(scanner.Text(), "{")
			line[1] = strings.TrimSuffix(line[1], "}")
			instr[line[0]] = parseRules(line[1])
		} else {
			parts = append(parts, parsePart(scanner.Text()))
		}
	}
	return instr, parts
}

func nextInstr(rules []rule, part part) string {
	next := "R"
	for _, rule := range rules {
		switch rule.t {
		case COND:
			switch rule.cond {
			case LESS:
				if part[rule.attr] < rule.value {
					return rule.result
				}
			case MORE:
				if part[rule.attr] > rule.value {
					return rule.result
				}
			}
		case DIR:
			return rule.result
		}
	}

	return next
}

func part1(instr map[string][]rule, parts []part) int {
	total := 0
	for _, part := range parts {
		next := "in"
		for next != "A" && next != "R" {
			next = nextInstr(instr[next], part)
		}
		if next == "A" {
			for _, v := range part {
				total += v
			}
		}
	}

	return total
}

func mapClone(src rangePart) rangePart {
	dst := make(rangePart, len(src))
	for k, v := range src {
		dst[k] = &Range{v.lower, v.upper}
	}
	return dst
}

func part2(next string, rp rangePart, instrMap map[string][]rule) int {
	if next == "A" {
		total := 1
		for k := range rp {
			total *= (rp[k].upper - rp[k].lower + 1)
		}
		return total
	}
	total := 0
	clonedMap := mapClone(rp)
	for _, rule := range instrMap[next] {
		switch rule.t {
		case COND:
			switch rule.cond {
			case LESS:
				nextMap := mapClone(clonedMap)
				nextMap[rule.attr].upper = rule.value - 1
				total += part2(rule.result, nextMap, instrMap)
				clonedMap[rule.attr].lower = rule.value
			case MORE:
				nextMap := mapClone(clonedMap)
				nextMap[rule.attr].lower = rule.value + 1
				total += part2(rule.result, nextMap, instrMap)
				clonedMap[rule.attr].upper = rule.value
			}
		case DIR:
			total += part2(rule.result, clonedMap, instrMap)
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
	instr, parts := parseInstr(scanner)

	result1 := part1(instr, parts)
	fmt.Println("Part 1 result:", result1)

	result2 := part2(
		"in",
		rangePart{"x": &Range{1, 4000}, "m": &Range{1, 4000}, "a": &Range{1, 4000}, "s": &Range{1, 4000}},
		instr,
	)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
