package main

func Decode(runes []rune) (map[string]interface{}, error) {
	var err error

	runePeeker := newPeeker[rune](runes)
	tokens, err := tokenize(runePeeker)
	if err != nil {
		return nil, err
	}

	tkPeeker := newPeeker[token](tokens)
	res, err := parse(tkPeeker)
	if err != nil {
		return nil, err
	}

	return res, err
}
