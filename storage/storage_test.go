package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryMetadataValueAndScan(t *testing.T) {
	t.Parallel()

	input := &HistoryMetadata{
		MachineID:  "machine-1",
		EntryID:    10,
		Wave:       2,
		TotalWave:  5,
		IsLastWave: false,
	}

	val, err := input.Value()
	require.NoError(t, err)

	var got HistoryMetadata
	err = got.Scan(val)
	require.NoError(t, err)
	assert.Equal(t, *input, got)
}

func TestHistoryMetadataScanInvalidType(t *testing.T) {
	t.Parallel()

	var got HistoryMetadata
	err := got.Scan("not-bytes")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestErrorDetailValueAndScan(t *testing.T) {
	t.Parallel()

	input := &ErrorDetail{
		Err:     "something failed",
		Message: "wrapped message",
	}

	val, err := input.Value()
	require.NoError(t, err)

	var got ErrorDetail
	err = got.Scan(val)
	require.NoError(t, err)
	assert.Equal(t, *input, got)
}

func TestErrorDetailScanInvalidType(t *testing.T) {
	t.Parallel()

	var got ErrorDetail
	err := got.Scan(42)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestNoopClient(t *testing.T) {
	t.Parallel()

	c := NewNoopClient()
	err := c.WriteHistory(context.Background(), &History{})
	require.NoError(t, err)

	data, err := c.ReadHistories(context.Background(), &HistoryFilter{})
	require.NoError(t, err)
	assert.Nil(t, data)
}
