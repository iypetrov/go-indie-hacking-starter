package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Toast struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (t Toast) Error() string {
	return fmt.Sprintf("custom error: %s", t.Message)
}

type alert struct {
	Toast Toast `json:"add-toast"`
}

func AddToast(w http.ResponseWriter, t Toast) {
	a := alert{
		Toast: t,
	}

	res, err := json.Marshal(a)
	if err != nil {
		return
	}
	w.Header().Set("HX-Trigger", string(res))
}
