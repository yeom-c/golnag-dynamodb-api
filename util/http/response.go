package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

func newResponse(status int, data interface{}) *response {
	return &response{
		Status: status,
		Result: data,
	}
}

func (resp *response) bytes() []byte {
	data, _ := json.Marshal(resp)
	return data
}

func (resp *response) string() string {
	return string(resp.bytes())
}

func (resp *response) sendResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.Status)
	_, _ = w.Write(resp.bytes())
	log.Println(resp.string())
}

func StatusOK(w http.ResponseWriter, r *http.Request, data interface{}) {
	newResponse(http.StatusOK, data).sendResponse(w, r)
}

func StatusNoContent(w http.ResponseWriter, r *http.Request) {
	newResponse(http.StatusOK, nil).sendResponse(w, r)
}

func StatusBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	newResponse(http.StatusBadRequest, data).sendResponse(w, r)
}

func StatusNotFound(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	newResponse(http.StatusNotFound, data).sendResponse(w, r)
}

func StatusMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	newResponse(http.StatusMethodNotAllowed, nil).sendResponse(w, r)
}

func StatusConflict(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	newResponse(http.StatusConflict, data).sendResponse(w, r)
}

func StatusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	newResponse(http.StatusInternalServerError, data).sendResponse(w, r)
}
