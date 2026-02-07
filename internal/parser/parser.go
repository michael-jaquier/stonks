package parser

import (
	"io"
)

type ParseSource interface {
	Parse(r io.Reader) error
}

type StreamParser[T any] struct {
	P      ParseSource
	Stream <-chan T
}

func NewStreamProcessor[T any](p ParseSource, stream chan T) *StreamParser[T] {
	return &StreamParser[T]{
		P:      p,
		Stream: stream,
	}
}
