package generator

import (
	"encoding/json"
	"fmt"
)

type TemplateName string
type FieldName string

// Metadata represents the metadata of the generated code.
// FIXME: concurrent map is not safe
type Metadata map[TemplateName]Schema

func (m Metadata) Each(fn func(TemplateName, Schema) error) error {
	for tn, om := range m {
		if err := fn(tn, om); err != nil {
			return err
		}
	}
	return nil
}

type Schema map[FieldName]any

func (om Schema) GetDriver() (string, error) {
	if v, ok := om["db"]; ok {
		if s, ok := v.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("driver is not string")
	}
	return "", fmt.Errorf("driver not found")
}

func (om Schema) GetTable(dr string) (string, error) {
	switch dr {
	case "mysql":
		if v, ok := om["table"]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
			return "", fmt.Errorf("table is not string")
		}
	case "mysqlr":
		if v, ok := om["dbtable"]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
			return "", fmt.Errorf("table is not string")
		}
	}
	return "", fmt.Errorf("table not found")
}

type TMetadata struct {
	Namespace     string
	Pkg           string
	Meta          Metadata
	Input         string
	Output        string
	DisableRawSQL bool
}

func (tm TMetadata) IsDisableRawSQL() bool {
	return tm.DisableRawSQL
}

func (tm TMetadata) Encode() ([]byte, error) {
	return json.Marshal(tm)
}

func (tm *TMetadata) Decode(data []byte) error {
	return json.Unmarshal(data, &tm)
}

type Generator interface {
	Generate(TMetadata) error
	DriverName() string
}

func Render(meta TMetadata, generators ...Generator) error {

	for _, generator := range generators {
		if err := generator.Generate(meta); err != nil {
			return err
		}
	}
	return nil
}
