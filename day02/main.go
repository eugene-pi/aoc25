package main

import (
	"io"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	start int
	stop int
	half int 
}	

func createRange(l, r string) Range {
	start, _ := strconv.Atoi(l)
	stop, _ := strconv.Atoi(r)
	rv := Range{start: start, stop: stop}
	size := len(l)
	if size % 2 == 1 {
		half := 1
		for i := 0; i < size/2; i++ {
			half *= 10
		}
		rv.half = half
	} else {
		rv.half, _ = strconv.Atoi(l[:size/2])
	}
	return rv
}

func (r *Range) getNext() int {
	r.half += 1
	s := strconv.Itoa(r.half)
	rv, _ := strconv.Atoi(s + s)
	return rv
}

func checkRange(r *Range) (total, count int) {
	total = 0
	count = 0
	s := strconv.Itoa(r.half)
	current, _ := strconv.Atoi(s + s)
	for current <= r.stop {
		if current >= r.start {
			total += current
			count += 1
		}
		current = r.getNext()
	}
	return
}

func readInput(filename string, ch chan Range) {
	file, _ := os.Open(filename)
	content, _ := io.ReadAll(file)
	ranges := strings.Split(string(content), ",")

	for _, line := range ranges {
		borders := strings.Split(line, "-")
		r := createRange(borders[0], borders[1])
		ch <- r
	}
	file.Close()
}

func process1(ch chan Range, chResult chan int) {
	total, count := 0, 0
	for r := range ch {
		t, c := checkRange(&r)
		total += t
		count += c
	}
	fmt.Println("Accumulated value:", total)
	fmt.Println("Number of times invalid ids:", count)
	chResult <- total
	close(chResult)
}

func checkPattern(value int) bool {
	s := strconv.Itoa(value)
	l := len(s)
	for repeats := 2; repeats <= l; repeats++ {
		if l % repeats == 0 {
			s0 := strings.Repeat(s[0 : l/repeats], repeats)
			if s0 == s {
				return true
			}
		}
	}
	return false
}

func process2(ch chan Range, chResult chan int) {
	total, count := 0, 0
	for r := range ch {
		for value := r.start; value <= r.stop; value++ {
			if checkPattern(value) {
				total += value
				count += 1
			}
		}
	}
	fmt.Println("Accumulated value:", total)
	fmt.Println("Number of times invalid ids:", count)
	chResult <- total
	close(chResult)
}

func part1(filename string) {
	ch := make(chan Range)
	chResult := make(chan int)
	go process1(ch, chResult)
	readInput(filename, ch)
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func part2(filename string) {
	ch := make(chan Range)
	chResult := make(chan int)
	go process2(ch, chResult)
	readInput(filename, ch)
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func main() {
	filename := "input.txt" // "test.txt" //  
	part1(filename)
	part2(filename)
}
