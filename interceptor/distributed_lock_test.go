package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistributedLock(t *testing.T) {
	type args struct {
		serviceName string
		mode        string
		dl          DistributedLockItf
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
			got := DistributedLock(tt.args.serviceName, tt.args.mode, tt.args.dl, tt.args.sc)
			assert.NotNil(t, got)
		})
	}
}
