package jsonschema

import (
	"fmt"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/jsonschema/datagen"
)

func GenerateJSONDataByJSONSchema(schema []byte) ([]byte, error) {
	generator := datagen.New()

	data, err := generator.Generate(schema)
	if err != nil {
		return nil, fmt.Errorf("generate data based on json schema: %w", err)
	}

	return data, nil
}
