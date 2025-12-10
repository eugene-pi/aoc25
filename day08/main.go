package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)
type Result struct {
	part1 int
	part2 int
}

type JBox struct {
	x,y,z int
}

type Edge struct {
	from, to int
	distance float64
}

type Playground struct {
	boxes []JBox
	edges []Edge
	nodes []*TreeNode
}

type TreeNode struct {
	parent *TreeNode
	boxIndex int
	size int
}

func findRoot(node *TreeNode) *TreeNode {
	if node.parent != node {
		return findRoot(node.parent)
	}
	return node
}

// we don't need actual connecion, the goal is to group the nodes
// into "circuits", so I can merge the root nodes instead
func connect(t1, t2 *TreeNode) {
	r1 := findRoot(t1)
	r2 := findRoot(t2)
	if r1 == r2 {
		return
	}
	r1.size += r2.size
	r2.size = 0
	r2.parent = r1
}

func calcDistance(a, b JBox) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return math.Sqrt(float64(dx*dx) + float64(dy*dy) + float64(dz*dz))
}

// createPlayground creates all possible edges and sort them by distance
func createPlayground(boxes []JBox) Playground {
	p := Playground{
		boxes: boxes,
		edges: make([]Edge, 0),
		nodes: make([]*TreeNode, len(boxes)),
	}
	for i := 0; i < len(p.boxes); i++ {
		for j := i+1; j < len(p.boxes); j++ {
			edge := Edge{from: i, to: j}
			edge.distance = calcDistance(p.boxes[i], p.boxes[j])
			p.edges = append(p.edges, edge)
		}
	}
	for i := 0; i < len(p.boxes); i++ {
		node := &TreeNode{boxIndex: i, size: 1}
		node.parent = node
		p.nodes[i] = node
	}
	slices.SortFunc(p.edges, func(a, b Edge) int {
		if a.distance < b.distance {
			return -1
		} else if a.distance > b.distance {
			return 1
		}
		return 0
	})
	return p
}

func parseInput(filename string, ch chan string) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	close(ch)
	file.Close()
}

func process(connectionCount int, ch chan string, chResult chan Result) {
	// parse coords
	boxes := make([]JBox, 0)
	for s := range ch {
		s2 := strings.Split(s,",")
		var jb JBox
		fmt.Sscanf(s2[0], "%d", &jb.x)
		fmt.Sscanf(s2[1], "%d", &jb.y)
		fmt.Sscanf(s2[2], "%d", &jb.z)
		boxes = append(boxes, jb)
	}
	p := createPlayground(boxes)
	fmt.Printf("Top edge %+v\n", p.edges[0])
	for i:=range connectionCount {
		connect(p.nodes[p.edges[i].from], p.nodes[p.edges[i].to])
	}
	sizes := make([]int, len(p.nodes))
	for i, node := range p.nodes {
		sizes[i] = node.size
	}
	slices.SortFunc(sizes, func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	})
	Result := Result{
		part1: sizes[0] * sizes[1] * sizes[2],
	}
	for i:=connectionCount; i<len(p.edges); i++ {
		connect(p.nodes[p.edges[i].from], p.nodes[p.edges[i].to])
		if findRoot(p.nodes[0]).size == len(p.boxes) {
			Result.part2 = p.boxes[p.edges[i].from].x * p.boxes[p.edges[i].to].x
			chResult <- Result
			close(chResult)
			return
		}
	}
}

func calcTop3(filename string, connectionCount int) {
	ch := make(chan string)
	chResult := make(chan Result)
	
	go process(connectionCount, ch, chResult)
	parseInput(filename, ch)
	result := <- chResult
	fmt.Println("Part1", result.part1)
	fmt.Println("Part2", result.part2)
}

func main() {
	filename := "test.txt" //"input.txt" //  
	calcTop3(filename, 10)
	calcTop3("input.txt", 1000)
}
