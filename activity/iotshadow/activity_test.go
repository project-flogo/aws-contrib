package iotshadow

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayloadToOut(t *testing.T) {

	payload := `{
		"state": {
			"desired": {
				"attribute1": 1,
				"attribute2": "string2"
			},
			"reported": {
				"attribute1": 2,
				"attribute2": "string1"
			},
			"delta": {
				"attribute3": 1,
				"attribute5": "stringY"
			}
		},
		"metadata": {
			"desired": {
				"attribute1": {
					"timestamp": 123
				},
				"attribute2": {
					"timestamp": 123
				}
			},
			"reported": {
				"attribute1": {
					"timestamp": 123
				},
				"attribute2": {
					"timestamp": 123
				}
			}
		},
		"timestamp": 123,
		"clientToken": "token",
		"version": 1
	}`
	out := &Output{}
	err := json.Unmarshal([]byte(payload), &out.Result)
	assert.Nil(t, err)
}
