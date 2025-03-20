package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	ww := NewCustomResponseWriter(w)
	ww.Header().Set("test", "value")
	_, err := ww.Write([]byte("abcd"))
	assert.NoError(t, err) //nolint: testifylint // ignore error
	assert.NotEmpty(t, ww.Header())
	assert.Equal(t, "value", ww.Header().Get("test"))
	assert.Equal(t, http.StatusOK, ww.StatusCode())
	ww.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, ww.StatusCode())
}
