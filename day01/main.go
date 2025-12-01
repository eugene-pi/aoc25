package main

import (
	"bufio"
	"fmt"
	"os"
)

// split returns actual dial increment and number of times zero is traversed
func split(x int) (int, int) {
	if x >= 100 {
		return x % 100, x / 100
	} else if x <= -100 {
		return x % 100, -1 * (x / 100)
	}
	return x, 0
}

type dialer struct {
	value int
	zeroCount int
	traverseZeroCount int
}	

func (d *dialer) dial(add int) {
	count, extra := split(add)
	d.traverseZeroCount += extra
	d.value += count
	if d.value >= 100 {
		d.value = d.value - 100
		if d.value > 0 {
			d.traverseZeroCount += 1
		}
	} else if d.value < 0 {
		if d.value > count {
			d.traverseZeroCount += 1
		}
		d.value = 100 + d.value
	}
	if d.value == 0 {
		d.zeroCount++
	}
}

func readInput(filename string, ch chan int) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s := scanner.Text()
		btDirection := byte(0)
		count := 0
		fmt.Sscanf(s,"%c%d", &btDirection, &count)
		if btDirection == 'L' {
			count =  -count
		}
		ch <- count
	}
	file.Close()
}

func process(ch chan int, chResult chan int) {
	d := dialer{value: 50, zeroCount: 0}
	for val := range ch {
		d.dial(val)
	}
	fmt.Println("Final dial value:", d.value)
	fmt.Println("Number of times dial hit zero:", d.zeroCount)
	fmt.Println("Number of times dial traversed zero:", d.traverseZeroCount)
	chResult <- d.zeroCount + d.traverseZeroCount
	close(chResult)
}

func main() {
	ch := make(chan int)
	chResult := make(chan int)
	go process(ch, chResult)

	filename := "test.txt" // "input.txt"
	readInput(filename, ch)
	close(ch)
	result := <- chResult
	fmt.Println("Result received:", result)

}
