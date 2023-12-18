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
		return nil, createErrorWrongToken("'{'", nextToken)
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
			return nil, createErrorWrongToken("':'", nextToken)
		}
		peeker.next()

		mapValue, err = parseValue(peeker)
		if err != nil {
			return nil, err
		}

		nextToken = peeker.peek()
		result[mapKey] = mapValue

		if nextToken.kind == curlyBracket && nextToken.value == closeCurlyBracket {
			peeker.next()
			return result, nil
		} else if nextToken.kind == comma {
			peeker.next()
			nextToken = peeker.peek()
			continue
		} else {
			return nil, createErrorWrongToken("'}' or ','", nextToken)
		}
	}

	return nil, errors.ErrUnsupported
}

func parseSlice(peeker *peeker[token]) ([]interface{}, error) {
	result := make([]interface{}, 0)
	peeker.next()
	nextToken := peeker.peek()

	var value interface{}
	var err error
	for nextToken != nil {
		if nextToken.kind == squareBracket && nextToken.value == closeSquareBracket {
			peeker.next()
			return result, nil
		}

		value, err = parseValue(peeker)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
		nextToken = peeker.peek()

		if (nextToken.kind != squareBracket && nextToken.value != closeSquareBracket) &&
			nextToken.kind != comma {
			return nil, createErrorWrongToken("']' or ','", nextToken)
		}

		if nextToken.kind == comma {
			peeker.next()
		}

		nextToken = peeker.peek()
	}

	return nil, errors.New("expected ]")
}

func parseString(peeker *peeker[token]) (string, error) {
	var str string
	peeker.next()
	nextToken := peeker.peek()

	if nextToken.kind != literal {
		return "", createErrorWrongToken("literal", nextToken)
	}

	str = nextToken.value
	peeker.next()
	nextToken = peeker.peek()

	if nextToken.kind != quotes {
		return "", createErrorWrongToken("'\"'", nextToken)
	}
	peeker.next()

	return str, nil
}

func parseValue(peeker *peeker[token]) (interface{}, error) {
	nextToken := peeker.peek()

	var value interface{}
	var err error
	if nextToken.kind == quotes {
		value, err = parseString(peeker)
	} else if nextToken.kind == curlyBracket && nextToken.value == openCurlyBracket {
		value, err = parseObject(peeker)
	} else if nextToken.kind == literal {
		value, err = parseUnqoutedLiteral(peeker)
	} else if nextToken.kind == squareBracket && nextToken.value == openSquareBracket {
		value, err = parseSlice(peeker)
	} else {
		return nil, errors.New("couldn't parse any value")
	}

	if err != nil {
		return nil, err
	}

	return value, nil
}

func parseUnqoutedLiteral(peeker *peeker[token]) (interface{}, error) {
	nextToken := peeker.peek()

	var value interface{}
	var err error
	if nextToken.value == "true" {
		value = true
	} else if nextToken.value == "false" {
		value = false
	} else if nextToken.value == "null" {
		value = nil
	} else {
		value, err = strconv.ParseFloat(nextToken.value, 64)
		if err != nil {
			return nil, createErrorWrongToken("true, false, null or number", nextToken)
		}
	}

	peeker.next()
	return value, nil
}

func createErrorWrongToken(possibleTokens string, token *token) error {
	errorMessage := fmt.Sprintf(
		"was expecting %s token at position row %d col %d, got %s instead",
		possibleTokens, token.positionRow+1, token.positionCol+1, token.value)
	return errors.New(errorMessage)
}
