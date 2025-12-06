package main

import (
	"io"
	"fmt"
	"os"
	"strings"
)

type fnProcess func(ch chan Grid, chResult chan int)

type Grid struct {
	cells [][]byte
	size int
}

func (g *Grid) String() string {
	var sb strings.Builder
	for i := 0; i < g.size; i++ {
		sb.WriteString(string(g.cells[i]))
		sb.WriteString("\n")
	}
	return sb.String()
}

func dumpMap(mp [][]int) {
	var sb strings.Builder
	l := len(mp)
	for i := 0; i < l; i++ {
		for j := 0; j < l; j++ {
			sb.WriteString(fmt.Sprintf("%d", mp[i][j]))
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}


func (g *Grid) countNeighbors(x int, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if g.cells[x+i][y+j] == '@' {
				count++
			}
		}
	}
	return count - 1
}

func expandLines(lines []string) []string {
	l := len(lines)
	res := make([]string, l+2)
	for i := 0; i < l; i++ {
		res[i+1] = "." + lines[i] + "."
	}
	res[0] = strings.Repeat(".", l+2)
	res[l+1] = strings.Repeat(".", l+2)
	return res
}

func createGrid(lines []string) Grid {
	lines = expandLines(lines)
	size := len(lines)
	cells := make([][]byte, size)
	for i := 0; i < size; i++ {
		cells[i] = []byte(lines[i])
	}
	return Grid{cells: cells, size: size}
}

func readInput(filename string, ch chan Grid) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	ch <- createGrid(strings.Split(string(content), "\r\n"))
}

func makeNeighborMap(grid Grid) [][]int {
	l := grid.size
	res := make([][]int, l)
	for i := 0; i < l; i++ {
		res[i] = make([]int, l)
	}
	for i := 1; i < l-1; i++ {
		for j := 1; j < l-1; j++ {
			if grid.cells[i][j] == '@' {
				res[i][j] = grid.countNeighbors(i,j)
			}
		}
	}
	return res
}

func processNeighnorhoodMap(neighborMap [][]int, grid *Grid) int {
	total := 0
	l := grid.size
	for i := 1; i < l-1; i++ {
		for j := 1; j < l-1; j++ {
			if grid.cells[i][j] == '@' {
				if neighborMap[i][j] < 4 {
					total++
					// I can filter it now, but I will do that later
					// otherwise I have to process the neighbors again
					// which nat lead to extra checks
					grid.cells[i][j] = 'P' // P = pending
				}
			}
		}
	}
	return total
}

func RemoveRollsFromMap(neighborMap [][]int, grid *Grid) {
	l := grid.size
	for i := 1; i < l-1; i++ {
		for j := 1; j < l-1; j++ {
			if grid.cells[i][j] == 'P' {
				grid.cells[i][j] = 'R' // R = removed
				for ii := i - 1; ii <= i + 1; ii++ {
					for jj := j - 1; jj <= j +1; jj++ {
						if grid.cells[ii][jj] == '@' {
							neighborMap[ii][jj] -= 1
						}
					}
				}
			}
		}
	}
}


func process1(ch chan Grid, chResult chan int) {
	total := 0
	grid := <- ch
	mp := makeNeighborMap(grid)
	total = processNeighnorhoodMap(mp, &grid)
	chResult <- total
	close(chResult)
}

func process2(ch chan Grid, chResult chan int) {
	total := 0
	grid := <- ch
	mp := makeNeighborMap(grid)
	removed := 1
	for removed > 0 {
		removed = processNeighnorhoodMap(mp, &grid)
		if removed > 0 {
			fmt.Println("Removed in iteration:", removed)
			RemoveRollsFromMap(mp, &grid)
			total += removed
		}
	}
	chResult <- total
	close(chResult)
}


func calcRollsCount(filename string, fn fnProcess) {
	ch := make(chan Grid)
	chResult := make(chan int)
	go fn(ch, chResult)
	readInput(filename, ch)
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func main() {
	filename := "input.txt" //"test.txt" //  
	fmt.Println("Part1: Count accessible rolls")
	calcRollsCount(filename, process1)
	fmt.Println("Part2: Count all accessible rolls")
	calcRollsCount(filename, process2)
}
