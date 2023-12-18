package main

type peeker[T any] struct {
	data            []T
	currentPosition int
}

func (pk *peeker[T]) peek() *T {
	nextPosition := pk.currentPosition + 1
	if nextPosition >= len(pk.data) {
		return nil
	}

	return &pk.data[nextPosition]
}

func (pk *peeker[T]) next() *T {
	nextPosition := pk.currentPosition + 1
	if nextPosition >= len(pk.data) {
		return nil
	}

	pk.currentPosition = nextPosition
	return &pk.data[nextPosition]
}

func newPeeker[T any](data []T) *peeker[T] {
	return &peeker[T]{data: data, currentPosition: -1}
}
