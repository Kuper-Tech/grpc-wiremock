package schematomock

import (
	"fmt"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/jsonschema"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/mock"
)

type pathData struct {
	method string
	path   string

	op *openapi3.Operation
}

func Parse(descriptor *openapi3.T) ([]mock.Mock, error) {
	mocks := make([]mock.Mock, 0)
	components := descriptor.Components

	handleOperationF := func(operation *openapi3.Operation, dt pathData) {
		dt.op = operation
		mocks = append(mocks, createMocks(dt, components)...)
	}

	for path, pathItem := range descriptor.Paths {
		dt := pathData{path: path}

		if op := pathItem.Get; op != nil {
			dt.method = http.MethodGet
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Post; op != nil {
			dt.method = http.MethodPost
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Delete; op != nil {
			dt.method = http.MethodDelete
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Patch; op != nil {
			dt.method = http.MethodPatch
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Put; op != nil {
			dt.method = http.MethodPut
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Head; op != nil {
			dt.method = http.MethodHead
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Options; op != nil {
			dt.method = http.MethodOptions
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Trace; op != nil {
			dt.method = http.MethodTrace
			handleOperationF(op, dt)

			continue
		}

		if op := pathItem.Connect; op != nil {
			dt.method = http.MethodConnect
			handleOperationF(op, dt)

			continue
		}
	}

	return mocks, nil
}

func createMocks(data pathData, components *openapi3.Components) []mock.Mock {
	mocks := make([]mock.Mock, 0)
	responses := data.op.Responses
	definedStatusCodes := getDefinedStatusCodes(responses)

	for statusCodeRaw, response := range responses {
		if response == nil {
			log.Printf("parse openapi: path: %s, err: nil response (skip)\n", data.path)
			continue
		}

		if response.Value == nil {
			log.Printf("parse openapi: path: %s, err: nil response value (skip)\n", data.path)
			continue
		}

		responseBodyRaw, err := createMockResponseBody(response, components)
		if err != nil {
			log.Printf("parse openapi: path: %s, err: create mock body: %s (skip)\n", data.path, err.Error())
			continue
		}

		responseBody := string(responseBodyRaw)

		statuses, err := getStatusCodes(definedStatusCodes, statusCodeRaw)
		if err != nil {
			log.Printf("parse openapi: path: %s, err: incorrect status code (skip)\n", data.path)
			continue
		}

		var description string
		if response.Value.Description != nil {
			description = *response.Value.Description
		}

		for _, status := range statuses {
			mock := mock.Mock{
				RequestMethod:  data.method,
				RequestUrlPath: data.path,

				ResponseBody: responseBody,
				Description:  description,

				ResponseStatusCode: status,
			}

			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func createMockResponseBody(response *openapi3.ResponseRef, components *openapi3.Components) ([]byte, error) {
	var mockData []byte

	for _, sref := range response.Value.Content {
		schema := sref.Schema

		if err := resolveSchema(schema, components); err != nil {
			return nil, fmt.Errorf("resolve schema: %w", err)
		}

		jsonSchema, err := schema.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("extract json schema: %w", err)
		}

		mockData, err = jsonschema.GenerateJSONDataByJSONSchema(jsonSchema)
		if err != nil {
			return nil, fmt.Errorf("handle json schema: %w", err)
		}

		break
	}

	return mockData, nil
}
