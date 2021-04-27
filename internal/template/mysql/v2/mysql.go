package v2

import (
	_ "embed"

	"github.com/ezbuy/ezorm/v2/internal/template"
)

var (
	//go:embed config.tpl
	configTemplate []byte
)

func init() {
	template.RegisterDriver(MySQL{})
}

type MySQL struct {
}

func (m MySQL) Name() string {
	return "mysql.v2"
}

func (m MySQL) Templates() ([]template.Template, error) {
	return []template.Template{
		{
			Data:     configTemplate,
			Filename: "mysql.v2.config.go",
		},
	}, nil
}
