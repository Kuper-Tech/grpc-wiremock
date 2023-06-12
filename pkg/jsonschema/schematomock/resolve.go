package schematomock

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

// resolveSchema recursively traverses all nested schemas and replaces references to values.
func resolveSchema(sref *openapi3.SchemaRef, components *openapi3.Components) error {
	if sref == nil {
		return fmt.Errorf("provided schema ref is empty")
	}

	ref := sref.Ref
	value := sref.Value

	if isRef(ref) {
		if err := handleAsRef(sref, components); err != nil {
			return fmt.Errorf("handle ref: %w", err)
		}

		return nil
	}

	if err := handleAsValue(value, components); err != nil {
		return fmt.Errorf("handle value: %w", err)
	}

	if sref.Value.Type == "" {
		sref.Value.Type = "object"
	}

	return nil
}

func handleAsValue(value *openapi3.Schema, components *openapi3.Components) error {
	if value == nil {
		return nil
	}

	if value.Properties == nil {
		value.Properties = make(openapi3.Schemas)
	}

	for _, property := range value.Properties {
		if err := resolveSchema(property, components); err != nil {
			return fmt.Errorf("property: %w", err)
		}
	}

	if items := value.Items; items != nil {
		value.Items = nil

		if err := resolveSchema(items, components); err != nil {
			return fmt.Errorf("item: %w", err)
		}

		if items.Value.Type == "" {
			items.Value.Type = "object"
		}

		value.Items = items
	}

	return nil
}

func handleAsRef(sref *openapi3.SchemaRef, components *openapi3.Components) error {
	schemas := components.Schemas

	if schemas == nil || len(schemas) == 0 {
		return fmt.Errorf("provided set of components is empty")
	}

	ref := sref.Ref
	refKey := filepath.Base(ref)

	referencedSchemaRef, ok := schemas[refKey]
	if !ok {
		return fmt.Errorf("schema referenced not found, reference path: %s", ref)
	}

	sref.Ref = ""

	if err := resolveSchema(referencedSchemaRef, components); err != nil {
		return fmt.Errorf("resolve schema ref: %w", err)
	}

	if err := copySchemaRef(sref, referencedSchemaRef); err != nil {
		return fmt.Errorf("copy schema ref: %w", err)
	}

	return nil
}

func isRef(ref string) bool {
	return ref != ""
}

func copySchemaRef(targetSchema, sourceSchema *openapi3.SchemaRef) error {
	if sourceSchema.Value == nil {
		targetSchema.Value = nil
		return nil
	}

	targetSchema.Value = openapi3.NewSchema()

	sourceData, err := json.Marshal(sourceSchema.Value)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err = json.Unmarshal(sourceData, targetSchema.Value); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	return nil
}
