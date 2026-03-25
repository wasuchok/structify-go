package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func Parse(reader io.Reader) (any, error) {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if err := ensureSingleValue(decoder); err != nil {
		return nil, err
	}

	return value, nil
}

func ParseBytes(data []byte) (any, error) {
	return Parse(bytes.NewReader(data))
}

func ensureSingleValue(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		return fmt.Errorf("invalid JSON: %w", err)
	}

	return errors.New("invalid JSON: input contains multiple JSON values")
}
