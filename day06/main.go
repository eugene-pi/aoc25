package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)

type fnProcess func(ch chan Problem, chResult chan int)

type Problem struct {
	values []int
	op bool // true = sum, false = product
}

func (p Problem) Solve() int {
	val := p.values[0]
	for _, v := range p.values[1:] {
		if p.op {
			val += v
		} else {
			val *= v
		}
	}
	return val
}

func mySplit(s, sep string) []string {
	l := strings.Count(s, sep) + 1
	res := make([]string, 0, l)
	for part := range strings.SplitSeq(s, sep) {
		if part != "" {
			res = append(res, part)
		}
	}
	return res
}

// preparseInput reads the input file and returns a slice of strings,
// where each string represents a line from the file with values separated by commas.
func preparseInput(filename string) []string {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	s := mySplit(scanner.Text(), " ")
	for scanner.Scan() {
		s1 := mySplit(scanner.Text(), " ")
		for i := 0; i < len(s); i++ {
			s[i] += "," + s1[i]
		}
	}
	file.Close()
	return s
}

// read files to a slice of strings
func preparseInput2(filename string) []string {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	s := make([]string, 0)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	file.Close()
	return s
}

func transposeInput(s []string) []string {
	l := len(s[0])
	res := make([]string, l)
	for i := 0; i < l; i++ {
		var sb strings.Builder
		for j := 0; j < len(s); j++ {
			sb.WriteByte(s[j][i])
		}
		res[i] = sb.String()
	}
	return res
}

func process1(ch chan Problem, chResult chan int) {
	total := 0
	for p := range ch {
		total += p.Solve()
	}
	chResult <- total
	close(chResult)
}

func calcGrandTotal(filename string, fn fnProcess) {
	ch := make(chan Problem)
	chResult := make(chan int)
	s := preparseInput(filename)
	go fn( ch, chResult)
	for _, line := range s {
		strValues := strings.Split(line, ",")
		p := Problem{
			values: make([]int, len(strValues) - 1),
			op:     strings.HasSuffix(line,"+"),
		}
		for i := 0; i < len(p.values); i++ {
			v, _ := strconv.Atoi(strValues[i])
			p.values[i] = v
		}
		ch <- p
	}
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func calcVerticalGrandTotal(filename string, fn fnProcess) {
	ch := make(chan Problem)
	chResult := make(chan int)
	s1 := transposeInput(preparseInput2(filename))
	go fn( ch, chResult)
	p := Problem{
		values: make([]int, 0),
		op:     true,
	}
	for _, line := range s1 {
		if strings.HasSuffix(line, "+") {
			p.op = true
		} else if strings.HasSuffix(line, "*") {
			p.op = false
		}
		s := strings.TrimSpace(line[:len(line)-1])
		if s == "" {
			ch <- p
			p = Problem{
				values: make([]int, 0),
				op:     true,
			}
		} else {
			v, _ := strconv.Atoi(s)
			p.values = append(p.values, v)
		}
	}
	if len(p.values) > 0 {
		ch <- p
	}
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)
}

func main() {
	filename := "input.txt" // "test.txt" // 
	fmt.Println("Part1: Find The Grand Total")
	calcGrandTotal(filename, process1)
	fmt.Println("Part2: Find The Vertical Grand Total")
	calcVerticalGrandTotal(filename, process1)
}
