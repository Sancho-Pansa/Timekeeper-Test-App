package api

import(
	"encoding/json"
	"net/http"
)

func writeJson(
	w http.ResponseWriter,
	status int,
	data any,
) {
	body, err := json.Marshal(data)

	if err != nil {
		http.Error(
			w,
			"Internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(body)
}