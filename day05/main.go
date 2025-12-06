package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

type fnProcess func(fr FreshRanges, ch chan int, chResult chan int)

type Range struct {
	start int
	end int
}

func (r Range) Power() int {
	return r.end - r.start + 1
}

type FreshRanges struct {
	ranges []Range
}

func (fr *FreshRanges) isFresh(id int) bool {
	idxHigh := len(fr.ranges) - 1
	for i := 0; i < len(fr.ranges); i++ {
		if fr.ranges[i].start > id {
			idxHigh = i
			break
		}
	}
	for i := idxHigh; i >= 0; i-- {
		if fr.ranges[i].end >= id && fr.ranges[i].start <= id {
			return true
		}
	}
	return false
}

func createFreshRanges(lines []string) FreshRanges {
	ranges1 := make([]Range, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		ranges1[i] = Range{
			start: start,
			end:   end,
		}
	}
	slices.SortFunc(ranges1, func(a, b Range) int {
		if a.start < b.start {
			return -1
		} else if a.start > b.start {
			return 1
		}
		return 0	
	})
	ranges := make([]Range, 0, len(ranges1))
	current := ranges1[0]
	for i := 1; i < len(ranges1); i++ {
		if ranges1[i].start <= current.end+1 {
			current.end = max(current.end, ranges1[i].end)
		} else {
			ranges = append(ranges, current)
			current = ranges1[i]
		}
	}
	ranges = append(ranges, current)
	return FreshRanges{
		ranges: ranges,
	}
}

func preparseInput(filename string) ([]string, []string) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	file.Close()
	s := strings.Split(string(content), "\r\n\r\n")
	return strings.Split(s[0], "\r\n"), strings.Split(s[1], "\r\n")
}

func process1(fr FreshRanges, ch chan int, chResult chan int) {
	total := 0
	for id := range ch {
		if fr.isFresh(id) {
			total++
		}
	}
	chResult <- total
	close(chResult)
}

func process2(fr FreshRanges, ch chan int, chResult chan int) {
	total := 0
	for range ch {
	}
	for _, r := range fr.ranges {
		total += r.Power()
	}
	chResult <- total
	close(chResult)
}


func calcFreshCount(filename string, fn fnProcess) {
	ch := make(chan int)
	chResult := make(chan int)
	s1, s2 := preparseInput(filename)
	go fn(createFreshRanges(s1), ch, chResult)
	for _, line := range s2 {
		id, _ := strconv.Atoi(line)
		ch <- id
	}
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func main() {
	filename := "input.txt" //"test.txt" //  
	fmt.Println("Part1: Count fresh products")
	calcFreshCount(filename, process1)
	fmt.Println("Part1: Count all possible fresh products")
	calcFreshCount(filename, process2)

}
