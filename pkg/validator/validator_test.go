package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/pkg/validator"
)

type input struct {
	URL string `validate:"required,url"`
}

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{name: "valid url", url: "https://example.com"},
		{name: "empty", url: "", wantErr: true},
		{name: "not a url", url: "not-a-url", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Validate(input{URL: tc.url})

			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
