package httpresponse_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
	"github.com/vickiliou/challenge-wex/internal/transaction"
)

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	someErr := errors.New("somme error")

	httpresponse.RespondWithError(w, http.StatusBadRequest, someErr)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	want := httpresponse.ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Message:    someErr.Error(),
	}

	var got httpresponse.ErrorResponse
	gotErr := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, gotErr)
	assert.Equal(t, want, got)
}

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()

	someStruct := transaction.RequestResponse{
		ID: "1",
	}

	httpresponse.RespondJSON(w, http.StatusOK, someStruct)
	assert.Equal(t, http.StatusOK, w.Code)

	var got transaction.RequestResponse
	gotErr := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, gotErr)
	assert.Equal(t, someStruct, got)
}
