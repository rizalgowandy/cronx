package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTelemetry(t *testing.T) {
	type args struct {
		metric TelemetryClientItf
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Telemetry(tt.args.metric)
			assert.NotNil(t, got)
		})
	}
}
