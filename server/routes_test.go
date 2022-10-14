package server

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/eine-doodka/twoStage/cache"
	"github.com/eine-doodka/twoStage/logic"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_commit(t *testing.T) {
	ctx := context.Background()
	redisCLient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	cacheModule := cache.NewImpl(redisCLient, 30*time.Second)
	logicModule := logic.NewImpl(cacheModule, 4)
	handlersModule := NewHandlers(logicModule)
	routes := NewServer(handlersModule)
	//init
	staticId := uuid.New()
	staticCode := "1234"
	err := cacheModule.Set(ctx, staticId, staticCode)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"id":   staticId.String(),
				"code": staticCode,
			},
			expectedCode: http.StatusFound,
		},
		{
			name: "wrong-code",
			payload: map[string]string{
				"id":   staticId.String(),
				"code": "asdf",
			},
			expectedCode: http.StatusExpectationFailed,
		},
		{
			name: "wrong-id",
			payload: map[string]string{
				"id":   uuid.New().String(),
				"code": staticCode,
			},
			expectedCode: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.payload)
		req := httptest.NewRequest(http.MethodPost, "/commit", b)
		routes.ServeHTTP(rec, req)
		assert.Equal(t, tc.expectedCode, rec.Code)
	}
}
