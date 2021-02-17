package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func fromFile(fileName string) ([]float64, []float64) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	XValues := []float64{}
	YValues := []float64{}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		x, y, _ := decomposeLine(line)
		XValues = append(XValues, x)
		YValues = append(YValues, y)
	}

	return XValues, YValues
}

func decomposeLine(line string) (float64, float64, int64) {
	res := strings.Split(line, " ")
	if len(res) != 3 {
		log.Fatal(fmt.Errorf("Invalid file format"))
	}
	x, err := strconv.ParseFloat(res[0], 64)
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.ParseFloat(res[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	status, err := strconv.ParseInt(res[1], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return x, y, status
}

func mean(numbers []float64) float64 {
	total := 0.0
	for _, v := range numbers {
		total += v
	}
	return math.Round(total / float64(len(numbers)))
}

func median(numbers []float64) float64 {
	sort.Float64s(numbers)
	mNumber := len(numbers) / 2
	if len(numbers)%2 == 1 {
		return numbers[mNumber]
	}
	return (numbers[mNumber-1] + numbers[mNumber]) / 2
}
