package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type token struct {
	kind        tokenKind
	positionCol int
	positionRow int
	value       string
}

type tokenKind uint

const (
	curlyBracket tokenKind = iota
	squareBracket
	colon
	comma
	quotes
	literal
)

var currentCol = 0
var currentRow = 0

func tokenize(peeker *peeker[rune]) ([]token, error) {
	currentCol = 0
	currentRow = 0

	var tk *token
	var tks []token
	var count int

	tokens := make([]token, 0)

	var nextRune = peeker.peek()
	for nextRune != nil {
		if *nextRune == '\n' {
			currentRow++
			currentCol = 0
		}

		skipWhitespacesAndNewLines(peeker)
		if peeker.peek() == nil {
			break
		}

		tk, count = readBrackets(peeker)
		if tk != nil {
			tokens = append(tokens, *tk)
			currentCol += count
			continue
		}

		tk, count = readColonOrComma(peeker)
		if tk != nil {
			tokens = append(tokens, *tk)
			currentCol += count
			continue
		}

		tks, count = readLiteral(peeker)
		if tks != nil {
			tokens = append(tokens, tks...)
			currentCol += count
			continue
		}

		return nil, fmt.Errorf(
			"wasn't able to parse token at position - row %d, col %d",
			currentRow+1, currentCol+1)
	}

	skipWhitespacesAndNewLines(peeker)
	return tokens, nil
}

func skipWhitespacesAndNewLines(peeker *peeker[rune]) {
	nextRune := peeker.peek()
	for nextRune != nil && (unicode.IsSpace(*nextRune)) {
		if *nextRune == '\n' {
			currentRow++
			currentCol = 0
		} else {
			currentCol++
		}

		peeker.next()
		nextRune = peeker.peek()
	}
}

func readBrackets(peeker *peeker[rune]) (*token, int) {
	countRead := 0
	nextRune := peeker.peek()

	if *nextRune == '{' || *nextRune == '}' {
		peeker.next()
		return &token{
			kind:        curlyBracket,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol,
			positionRow: currentRow,
		}, countRead + 1
	}

	if *nextRune == '[' || *nextRune == ']' {
		peeker.next()
		return &token{
			kind:        squareBracket,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol,
			positionRow: currentRow,
		}, countRead + 1
	}

	return nil, countRead
}

func readColonOrComma(peeker *peeker[rune]) (*token, int) {
	countRead := 0
	nextRune := peeker.peek()

	if *nextRune == ':' {
		peeker.next()
		return &token{
			kind:        colon,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol,
			positionRow: currentRow,
		}, countRead + 1
	}

	if *nextRune == ',' {
		peeker.next()
		return &token{
			kind:        comma,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol,
			positionRow: currentRow,
		}, countRead + 1
	}

	return nil, countRead
}

func readLiteral(peeker *peeker[rune]) ([]token, int) {
	countRead := 0
	nextRune := peeker.peek()

	tokens := make([]token, 0)

	if *nextRune == '"' {
		countRead++
		tokens = append(tokens, token{
			kind:        quotes,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol,
			positionRow: currentRow,
		})

		peeker.next()
		nextRune = peeker.peek()

		runes := make([]rune, 0)
		for *nextRune != '"' {
			runes = append(runes, *nextRune)

			countRead++
			peeker.next()
			nextRune = peeker.peek()
		}

		tokens = append(tokens, token{
			kind:        literal,
			value:       string(runes),
			positionCol: currentCol + countRead,
			positionRow: currentRow,
		})

		tokens = append(tokens, token{
			kind:        quotes,
			value:       strconv.QuoteRune(*nextRune),
			positionCol: currentCol + countRead,
			positionRow: currentRow,
		})

		countRead++
		peeker.next()
		return tokens, countRead
	}

	runes := make([]rune, 0)
	for unicode.IsDigit(*nextRune) || unicode.IsLetter(*nextRune) || *nextRune == '.' {
		runes = append(runes, *nextRune)
		countRead++
		peeker.next()
		nextRune = peeker.peek()
	}

	if len(runes) > 0 {
		return []token{{
			kind:        literal,
			value:       string(runes),
			positionCol: currentCol,
			positionRow: currentRow,
		}}, countRead
	}

	return nil, countRead
}
