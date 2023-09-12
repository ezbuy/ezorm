package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	expect := map[string]any{
		"Pkg": "hello-driver",
		"Meta": map[string]any{
			"Hello": map[string]any{
				"db": "hello-driver",
			},
		},
		"Input":         "./e2e/plugins/hello-driver",
		"Output":        "./e2e/plugins/hello-driver",
		"DisableRawSQL": false,
		"Namespace":     "",
		"args": map[string]string{
			"key": "value",
		},
	}
	b, err := json.Marshal(expect)
	assert.NoError(t, err)
	got, err := os.ReadFile("metadata.json")
	assert.NoError(t, err)
	assert.JSONEq(t, string(b), string(got))
}
