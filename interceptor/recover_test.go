package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Recover()
			assert.NotNil(t, got)
		})
	}
}
