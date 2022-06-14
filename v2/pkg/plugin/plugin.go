package plugin

import (
	"encoding/json"

	"github.com/ezbuy/ezorm/v2/internal/generator"
)

type GeneratorRequest struct {
	meta *generator.TMetadata
}

// Schema exports the internal generator Schema type.
type Schema = generator.Schema

// TemplateName exports the internal generator TemplateName type.
type TemplateName = generator.TemplateName

func Decode(data []byte) (*GeneratorRequest, error) {
	req := &GeneratorRequest{
		meta: &generator.TMetadata{},
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
