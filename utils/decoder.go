package utils

import (
	"encoding/json"
	"io"
)

const (
	// MaxBodySize is the maximum size of request body in bytes (1MB)
	MaxBodySize = 1 << 20
)

var (
	// ErrBodyTooLarge indicates that the request body is too large
	ErrBodyTooLarge = &json.SyntaxError{Offset: 0}
)

// LimitedDecoder wraps json.Decoder with size limit
type LimitedDecoder struct {
	*json.Decoder
	reader io.Reader
}

// NewLimitedDecoder creates a new limited JSON decoder
func NewLimitedDecoder(reader io.Reader) *LimitedDecoder {
	limitedReader := io.LimitReader(reader, MaxBodySize)
	return &LimitedDecoder{
		Decoder: json.NewDecoder(limitedReader),
		reader:  limitedReader,
	}
}
