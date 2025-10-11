package src_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"users/src"
	httpApi "users/src/apis/http"
	"users/src/domain"
	"users/src/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleGetUser(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	userA := &domain.User{}

	var (
		config = &src.Config{
			HttpSessionCookieDomain: "abc",
			SessionDuration:         5 * time.Minute,
		}

		logBuf = bytes.NewBuffer([]byte{})
		logger = slog.New(slog.NewJSONHandler(logBuf, nil))

		databaseRepo, databaseClient = NewDatabaseRepository(t, userFoo)
		cacheRepo, cacheClient       = NewCacheRepository(t, userFoo)

		svc = service.New(config, databaseRepo, cacheRepo)
		api = httpApi.New(config, logger, svc)
	)

	tt := []struct {
		name string

		urlVals        url.Values
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "",
			urlVals:        url.Values{"id": []string{"abc"}},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id": "abc"}`,
		},
	}

	for _, tc := range tt {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/users?%s", tc.urlVals.Encode()), nil)
		api.HandleGetUser(recorder, request)

		body, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)

		assert.Equal(t, tc.expectedStatus, recorder.Code)
		assert.Equal(t, tc.expectedBody, body)
	}
}
