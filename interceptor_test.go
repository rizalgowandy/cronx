package cronx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	type args struct {
		interceptors []Interceptor
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.args.interceptors...)
			assert.NotNil(t, got)
		})
	}
}
