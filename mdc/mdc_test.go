package mdc_test

import (
	"context"
	"testing"

	"github.com/go-log4g/core/mdc"
	"github.com/stretchr/testify/require"
)

func TestMDC(test *testing.T) {
	ctx := context.Background()
	ctx = mdc.Put(ctx, "requestId", "123")
	ctx = mdc.PutAll(ctx, map[string]string{
		"userId":  "456",
		"traceId": "789",
	})

	requestId, ok := mdc.Get(ctx, "requestId")
	require.True(test, ok)
	require.Equal(test, "123", requestId)

	require.Equal(test, map[string]string{
		"requestId": "123",
		"userId":    "456",
		"traceId":   "789",
	}, mdc.Values(ctx))

	ctx = mdc.Remove(ctx, "userId")
	_, ok = mdc.Get(ctx, "userId")
	require.False(test, ok)

	ctx = mdc.Clear(ctx)
	require.Empty(test, mdc.Values(ctx))
}

func TestMDCImmutable(test *testing.T) {
	ctx1 := mdc.Put(context.Background(), "requestId", "123")
	ctx2 := mdc.Put(ctx1, "userId", "456")

	require.Equal(test, map[string]string{"requestId": "123"}, mdc.Values(ctx1))
	require.Equal(test, map[string]string{"requestId": "123", "userId": "456"}, mdc.Values(ctx2))

	values := mdc.Values(ctx2)
	values["requestId"] = "changed"

	requestId, _ := mdc.Get(ctx2, "requestId")
	require.Equal(test, "123", requestId)
}
