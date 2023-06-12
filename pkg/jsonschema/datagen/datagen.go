package datagen

import (
	"encoding/json"
	"fmt"

	gen "github.com/apicat/datagen"
)

type datagen struct{}

func New() datagen { return datagen{} }

// Generate allows to get randomized JSON samples based on provided JSON schema.
func (d datagen) Generate(schema []byte) ([]byte, error) {
	opt := gen.GenOption{}

	data, err := gen.JSONSchemaGen(schema, &opt)
	if err != nil {
		return nil, fmt.Errorf("call json schema gen: %w", err)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal json output: %w", err)
	}

	return dataBytes, nil
}
