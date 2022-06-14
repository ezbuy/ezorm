package plugin

import (
	"encoding/json"
	"fmt"
)

type Metadata map[TemplateName]Schema
type FieldName string

type TMetadata struct {
	Pkg           string
	Meta          Metadata
	Input         string
	Output        string
	DisableRawSQL bool
}

type GeneratorRequest struct {
	meta *TMetadata
}

// Schema exports the internal generator Schema type.
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

// TemplateName exports the internal generator TemplateName type.
type TemplateName string

func Decode(data []byte) (*GeneratorRequest, error) {
	req := &GeneratorRequest{
		meta: &TMetadata{},
	}
	if err := json.Unmarshal(data, req.meta); err != nil {
		return nil, err
	}
	return req, nil
}

func (req *GeneratorRequest) GetPackage() string {
	return req.meta.Pkg
}

func (req *GeneratorRequest) GetInputPath() string {
	return req.meta.Input
}

func (req *GeneratorRequest) GetOutputPath() string {
	return req.meta.Output
}

func (req *GeneratorRequest) Each(fn func(TemplateName, Schema) error) error {
	return req.meta.Meta.Each(fn)
}

func (m Metadata) Each(fn func(TemplateName, Schema) error) error {
	for tn, om := range m {
		if err := fn(tn, om); err != nil {
			return err
		}
	}
	return nil
}
