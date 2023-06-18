package actions

import (
	"encoding/json"
	"net/http"
)

func HealthAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "Ok",
	})
	if err != nil {
		panic(err)
	}
}
