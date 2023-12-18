package main

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	openCurlyBracket  = "'{'"
	closeCurlyBracket = "'}'"

	openSquareBracket  = "'['"
	closeSquareBracket = "']'"
)

func parse(peeker *peeker[token]) (map[string]interface{}, error) {
	return parseObject(peeker)
}

func parseObject(peeker *peeker[token]) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 0)
	nextToken := peeker.peek()

	if nextToken.kind != curlyBracket || nextToken.value != openCurlyBracket {
		errorMessage := fmt.Sprintf("Was expecting '{' token at position row %d col %d, got %s instead", nextToken.positionRow, nextToken.positionCol, nextToken.value)
		return nil, errors.New(errorMessage)
	}

	peeker.next()
	nextToken = peeker.peek()

	var mapKey string
	var mapValue interface{}
	var err error
	for nextToken != nil {
		mapKey, err = parseString(peeker)
		if err != nil {
			return nil, err
		}
		nextToken = peeker.peek()

		if nextToken.kind != colon {
			return nil, errors.New("")
		}
		peeker.next()
		nextToken = peeker.peek()

		if nextToken.kind == quotes {
			mapValue, err = parseString(peeker)
			if err != nil {
				return nil, err
			}
		} else if nextToken.kind == curlyBracket && nextToken.value == openCurlyBracket {
			mapValue, err = parseObject(peeker)
			if err != nil {
				return nil, err
			}
		} else if nextToken.kind == literal {
			mapValue, err = strconv.ParseFloat(nextToken.value, 64)
			if err != nil {
				return nil, errors.New("")
			}
			peeker.next()
		} else if nextToken.kind == squareBracket && nextToken.value == openSquareBracket {
			mapValue, err = parseSlice(peeker)
			if err != nil {
				return nil, errors.New("")
			}
		}
		nextToken = peeker.peek()
		result[mapKey] = mapValue

		if nextToken.kind == curlyBracket && nextToken.value == closeCurlyBracket {
			return result, nil
		} else if nextToken.kind == comma {
			peeker.next()
			nextToken = peeker.peek()
			continue
		} else {
			return nil, errors.New("Unexpected symbol")
		}
	}

	return nil, errors.ErrUnsupported
}

func parseSlice(peeker *peeker[token]) ([]interface{}, error) {
	return nil, nil
}

func parseString(peeker *peeker[token]) (string, error) {
	var str string
	nextToken := peeker.peek()
	if nextToken.kind != quotes {
		return "", errors.New("")
	}
	peeker.next()
	nextToken = peeker.peek()

	if nextToken.kind != literal {
		return "", errors.New("")
	}

	str = nextToken.value
	peeker.next()
	nextToken = peeker.peek()

	if nextToken.kind != quotes {
		return "", errors.New("")
	}
	peeker.next()

	return str, nil
}
