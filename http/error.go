package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	cabby "github.com/pladdy/cabby2"
	log "github.com/sirupsen/logrus"
)

func errorStatus(w http.ResponseWriter, title string, err error, status int) {
	errString := fmt.Sprintf("%v", err)

	te := cabby.Error{Title: title, Description: errString, HTTPStatus: status}

	log.WithFields(log.Fields{
		"error":       err,
		"title":       title,
		"http status": status,
	}).Error("Returning error in response")

	w.Header().Set("Content-Type", cabby.TaxiiContentType)
	w.WriteHeader(status)
	io.WriteString(w, resourceToJSON(te))
}

func badRequest(w http.ResponseWriter, err error) {
	errorStatus(w, "Bad Request", err, http.StatusBadRequest)
}

func forbidden(w http.ResponseWriter, err error) {
	errorStatus(w, "Forbidden", err, http.StatusForbidden)
}

func internalServerError(w http.ResponseWriter, err error) {
	errorStatus(w, "Internal Server Error", err, http.StatusInternalServerError)
}

func methodNotAllowed(w http.ResponseWriter, err error) {
	errorStatus(w, "Method Not Allowed", err, http.StatusMethodNotAllowed)
}

func resourceNotFound(w http.ResponseWriter, err error) {
	errorStatus(w, "Resource Not Found", err, http.StatusNotFound)
}

func requestTooLarge(w http.ResponseWriter, rc, mc int64) {
	err := fmt.Errorf("content length is %v, content length can't be bigger than %v", rc, mc)
	errorStatus(w, "Request Too large", err, http.StatusRequestEntityTooLarge)
}

func rangeNotSatisfiable(w http.ResponseWriter, err error) {
	errorStatus(w, "Requested Range Not Satisfiable", err, http.StatusRequestedRangeNotSatisfiable)
}

func unauthorized(w http.ResponseWriter, err error) {
	w.Header().Set("WWW-Authenticate", "Basic realm=TAXII 2.0")
	errorStatus(w, "Unauthorized", err, http.StatusUnauthorized)
}

func unsupportedMediaType(w http.ResponseWriter, err error) {
	errorStatus(w, "Unsupported Media Type", err, http.StatusUnsupportedMediaType)
}

func recoverFromPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Error("Panic!  Printing out Stack...")
		debug.PrintStack()
		resourceNotFound(w, errors.New("Resource not found"))
	}
}
