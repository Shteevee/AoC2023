package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type point struct {
	x int
	y int
}

type qItem struct {
	pos   point
	dir   rune
	dist  int
	index int
	prev  *qItem
}

type neighbor struct {
	pos   point
	dir   rune
	steps int
}

type PriorityQueue []*qItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*qItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func parseBlocks(scanner *bufio.Scanner) [][]int {
	blocks := [][]int{}
	for y := 0; scanner.Scan(); y++ {
		row := []int{}
		for _, c := range scanner.Text() {
			i, _ := strconv.Atoi(string(c))
			row = append(row, i)
		}
		blocks = append(blocks, row)
	}
	return blocks
}

func validNeighbor(n neighbor, maxX, maxY int, lastDir rune) bool {
	return n.pos.x >= 0 && n.pos.x < maxX &&
		n.pos.y >= 0 && n.pos.y < maxY &&
		n.dir != lastDir &&
		!(n.pos.x == 0 && n.pos.y == 0)
}

func nextNeighbors(qItem qItem, blocks [][]int, minStep, maxStep int) []neighbor {
	pos := qItem.pos
	neighbors := []neighbor{}
	for i := minStep; i <= maxStep; i++ {
		neighbors = append(neighbors,
			[]neighbor{
				{pos: point{x: pos.x - i, y: pos.y}, dir: 'h', steps: -i},
				{pos: point{x: pos.x + i, y: pos.y}, dir: 'h', steps: i},
				{pos: point{x: pos.x, y: pos.y - i}, dir: 'v', steps: -i},
				{pos: point{x: pos.x, y: pos.y + i}, dir: 'v', steps: i},
			}...)
	}

	n := 0
	for _, neighbor := range neighbors {
		if validNeighbor(neighbor, len(blocks[0]), len(blocks), qItem.dir) {
			neighbors[n] = neighbor
			n++
		}
	}
	return neighbors[:n]
}

func calcDist(item *qItem, n neighbor, blocks [][]int) int {
	curr, next, dist := item.pos, n.pos, item.dist
	switch n.dir {
	case 'h':
		if curr.x > next.x {
			for i := curr.x - 1; i >= next.x; i-- {
				dist += blocks[curr.y][i]
			}
		} else {
			for i := curr.x + 1; i <= next.x; i++ {
				dist += blocks[curr.y][i]
			}
		}
	case 'v':
		if curr.y > next.y {
			for i := curr.y - 1; i >= next.y; i-- {
				dist += blocks[i][curr.x]
			}
		} else {
			for i := curr.y + 1; i <= next.y; i++ {
				dist += blocks[i][curr.x]
			}
		}
	}
	return dist
}

func findPath(blocks [][]int, minStep, maxStep int) int {
	q := &PriorityQueue{&qItem{pos: point{x: 0, y: 0}}}
	heap.Init(q)
	seen := map[neighbor]bool{}
	end := point{x: len(blocks[0]) - 1, y: len(blocks) - 1}

	for len(*q) > 0 {
		u := heap.Pop(q).(*qItem)
		neighbors := nextNeighbors(*u, blocks, minStep, maxStep)
		if u.pos == end {
			return u.dist
		}
		for _, n := range neighbors {
			if seen[n] {
				continue
			}
			seen[n] = true
			nDist := calcDist(u, n, blocks)
			heap.Push(q, &qItem{pos: n.pos, dir: n.dir, dist: nDist, prev: u})
		}
	}

	return -1
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
	blocks := parseBlocks(scanner)

	result1 := findPath(blocks, 1, 3)
	fmt.Println("Part 1 result:", result1)

	result2 := findPath(blocks, 4, 10)
	fmt.Println("Part 2 result:", result2)

	log.Printf("Time taken: %s", time.Since(start))
}
