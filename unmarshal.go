package tame

import (
	"encoding/json"
	"errors"
	"fmt"
)

type recipeByteReader interface {
	Unmarshal([]byte) error
}
type recipeStringReader interface {
	Unmarshal(string) error
}

// Unmarshal reads data from document to provided targets
// Read operation is performed sequentially, first error occured
// will cancel rest of process
func Unmarshal(d Document, targets ...interface{}) error {
	if d == nil {
		return errors.New("nil document")
	}

	for i, j := range targets {
		switch j.(type) {
		case recipeByteReader:
			x := j.(recipeByteReader)
			if err := x.Unmarshal(d.GetBody()); err != nil {
				return err
			}
		case recipeStringReader:
			x := j.(recipeStringReader)
			if err := x.Unmarshal(string(d.GetBody())); err != nil {
				return err
			}
		case json.Unmarshaler:
			x := j.(json.Unmarshaler)
			if err := x.UnmarshalJSON(d.GetBody()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported target %t found at position %d", j, i)
		}
	}
	return nil
}
