package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type fnProcess func(ch chan string, chResult chan ManifoldState)

type ManifoldState struct {
	started bool
	splitterHitCount int
	tracksAvailable int
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

func process1(ch chan string, chResult chan ManifoldState) {
	state := ManifoldState{started: false, splitterHitCount: 0}
	var tracks []int
	prev := []byte{}
	for s := range ch {
		if !state.started {
			prev = []byte(strings.ReplaceAll(s, "S", "|"))
			state.started = true
			tracks = make([]int, len(prev))
			tracks[strings.IndexByte(s, 'S')] = 1
		} else {
			next := []byte(s)
			nextTracks := make([]int, len(next))
			for idx := 0; idx < len(next); idx++ {
				if prev[idx] == '|' {
					if next[idx] == '^' {
						state.splitterHitCount++
						next[idx-1] = '|'
						next[idx+1] = '|'
						nextTracks[idx-1] += tracks[idx]
						nextTracks[idx+1] += tracks[idx]
					} else {
						next[idx] = '|'
						nextTracks[idx] += tracks[idx]
					}
				}
			}
			prev = next
			tracks = nextTracks
		}
	}
	for _, v := range tracks {
		state.tracksAvailable += v
	}
	chResult <- state
	close(chResult)
}

func calcHitCount(filename string, fn fnProcess) {
	ch := make(chan string)
	chResult := make(chan ManifoldState)
	
	go fn( ch, chResult)
	parseInput(filename, ch)
	result := <- chResult
	fmt.Println("Part1: Splitter hit count = ", result.splitterHitCount)
	fmt.Println("Part2: Tracks count = ", result.tracksAvailable)
}

func main() {
	filename := "input.txt" // "test.txt" // 
	calcHitCount(filename, process1)
}
