package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	runes := readRunesFromFile("input.json")
	res, err := Decode(runes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(res.([]interface{})[0].(float64) + res.([]interface{})[1].(float64))
}

func readRunesFromFile(fileName string) []rune {
	file, err := os.Open(fileName)
	if err != nil {
		panic("Error reading file " + fileName + ":" + err.Error())
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	var runes []rune = make([]rune, 0)
	for scanner.Scan() {
		runes = append(runes, []rune(scanner.Text())[0])
	}

	return runes
}
