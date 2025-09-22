package httpx

import (
	"encoding/json"
	"net/http"
)

// WriteJSON will write the given data to the response writer as JSON.
func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonb)
	return err
}
