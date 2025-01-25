package fluffyjson_test

import (
	"testing"

	fluffyjson "github.com/hayash1/go-fluffy-json"
)

func TestAccess(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		accessor fluffyjson.Accessor
		err      error
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO
		})
	}
}
