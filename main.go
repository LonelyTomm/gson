package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	runes := readRunesFromFile("input.json")
	runePeeker := newPeeker[rune](runes)
	tokens, _ := tokenize(runePeeker)

	tkPeeker := newPeeker[token](tokens)
	res, err := parse(tkPeeker)
	fmt.Println("Number of tokens " + strconv.Itoa(len(res)) + err.Error())
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
