package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(filename string, ch chan string) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
	file.Close()
}

func findNlargest(s string, start int, n int) (pos int, val byte) {
	pos = start
	maxValue := s[pos]
	for i := start + 1; i < len(s) - n; i++ {
		if s[i] > maxValue {
			maxValue = s[i]
			pos = i
		}
	}
	val = maxValue
	return
}

func process2(ch chan string, chResult chan int) {
	total := 0
	for s := range ch {
		res := ""
		pos := -1
		val := byte(0)
		for i := 11; i >=0; i-- {
			pos, val = findNlargest(s, pos + 1, i)
			res = res + string(val)
		}
		joltage, _ := strconv.Atoi(res)
		fmt.Println(s, " becomes ", joltage)
		total += joltage
	}
	fmt.Println("Accumulated value:", total)
	chResult <- total
	close(chResult)
}

func process1(ch chan string, chResult chan int) {
	total := 0
	for s := range ch {
		cells := strings.Split(s, "")
		// find first largest cell
		idx0 := 0
		maxValue0 := s[0]
		for i := 1; i < len(cells) -1; i++ {
			if s[i] > maxValue0 {
				maxValue0 = s[i]
				idx0 = i
			}
		}
		idx1 := idx0 + 1
		maxValue2 := s[idx1]
		for i := idx0 + 2; i < len(cells); i++ {
			if s[i] > maxValue2 {
				maxValue2 = s[i]
				idx1 = i
			}
		}
		res := 10*(maxValue0 - '0') + (maxValue2 - '0')
		total += int(res)
	}
	fmt.Println("Accumulated value:", total)
	chResult <- total
	close(chResult)
}

func part1(filename string) {
	ch := make(chan string)
	chResult := make(chan int)
	go process1(ch, chResult)
	readInput(filename, ch)
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func part2(filename string) {
	ch := make(chan string)
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
