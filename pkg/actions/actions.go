package actions

import "net/http"

func jsonInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(`{"status": "Internal Server Error"}`))
	if err != nil {
		return
	}
}
