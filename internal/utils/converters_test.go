package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gqlgen-starter/internal/oops"
	"strings"
	"testing"
)

func TestID64(t *testing.T) {
	testCases := []struct {
		Input       string
		ExpectError bool
		Message     string
	}{
		{
			Input:       "10",
			ExpectError: false,
			Message:     "",
		},
		{
			Input:       "a",
			ExpectError: true,
			Message:     "a is not a valid number",
		},
	}

	for _, tc := range testCases {
		output, err := ID64(tc.Input)
		if tc.ExpectError {
			assert.ErrorContainsf(t, err, "not parsable as int64", "%s not parsable as int64", tc.Input)
			var codedErr *oops.CodedError
			if errors.As(err, &codedErr) {
				assert.True(t, strings.Contains(codedErr.HumanMessage, "not a valid number"))
			}

		} else {
			assert.Equal(t, output, int64(10))
		}
	}
}
