package cache

import (
	"context"
	"github.com/eine-doodka/twoStage/customerrors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryCache_GetSet(t *testing.T) {
	ctx := context.Background()
	cache := NewInMemoryCache()
	baseId := uuid.New()
	testCases := []struct {
		getId     uuid.UUID
		setId     uuid.UUID
		setCode   string
		getCode   string
		isSuccess bool
	}{
		{
			getId:     baseId,
			getCode:   "123",
			setId:     baseId,
			setCode:   "123",
			isSuccess: true,
		},
		{
			getId:     uuid.New(),
			getCode:   "123",
			setId:     baseId,
			setCode:   "123",
			isSuccess: false,
		},
		{
			getId:     baseId,
			getCode:   "1234",
			setId:     baseId,
			setCode:   "123",
			isSuccess: false,
		},
	}
	for _, tc := range testCases {
		err := cache.Set(ctx, tc.setId, tc.setCode)
		if err != nil {
			t.Fatal(err)
		}
		res, err := cache.Get(ctx, tc.getId)
		if err == customerrors.ErrNotFound {
			assert.Equal(t, tc.isSuccess, false)
		} else if err != nil {
			t.Fatal(err)
		} else {
			assert.Equal(t, res == tc.getCode, tc.isSuccess)
		}
	}
}
