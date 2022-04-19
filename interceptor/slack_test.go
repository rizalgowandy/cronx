package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifySlack(t *testing.T) {
	type args struct {
		serviceName string
		mode        string
		sc          SlackClientItf
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
			got := NotifySlack(tt.args.serviceName, tt.args.mode, tt.args.sc)
			assert.NotNil(t, got)
		})
	}
}
