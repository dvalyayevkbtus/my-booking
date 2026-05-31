package utils

import (
	"io"
	"net/http"
)

func MethodNotAllowed(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(rw, "Method not allowed!")
}

func InternalServerError(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusInternalServerError)
	io.WriteString(rw, "Internal server error!")
}

func BadRequest(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusBadRequest)
	io.WriteString(rw, "Bad request!")
}

func NotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	io.WriteString(rw, "Not found")
}

func Accepted(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusAccepted)
	io.WriteString(rw, "Accepted")
}

func SuccessString(rw http.ResponseWriter, body string) {
	io.WriteString(rw, body)
}
