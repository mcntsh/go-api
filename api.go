package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// Response represents an API response body that can write
// the output buffer of an HTTP response.
type Response struct {
	rw http.ResponseWriter

	Status *ResponseStatus `json:"status"`
	Body   interface{}     `json:"body"`
}

// ResponseStatus is a JSON-format representation of the
// HTTP status.
type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// WriteJSON marshals the Response into JSON format and attaches
// it to the HTTP response body.
func (r *Response) WriteJSON() {
	json, err := json.Marshal(r)
	if err != nil {
		logrus.Fatal("Could not marshal JSON in the API writer")
	}

	r.rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	r.rw.Header().Set("Content-Length", strconv.Itoa(len(json)))
	r.rw.WriteHeader(r.Status.Code)

	r.rw.Write(json)
}

func statusCodeFatal(code int) bool {
	fatalCodes := []int{500, 501, 502, 503, 504, 505, 511}
	for _, c := range fatalCodes {
		if c == code {
			return true
		}
	}

	return false
}

// WriteErrorResponse attaches the appropriate error code to the
// ResponseStatus and writes an error response. If the code is within
// a fatal error range (500, 501, 502...) it will log the error.
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	res := &Response{rw: w}

	res.Status = &ResponseStatus{
		Code:    code,
		Message: http.StatusText(code),
		Error:   err.Error(),
	}

	if statusCodeFatal(code) == true {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
			"code":   code,
			"error":  err,
		}).Warn("API Handler returned an error!")
	}

	res.WriteJSON()
}

// WriteResponse appends data to the Response and writes
// a successful (200) code to the output.
func WriteResponse(w http.ResponseWriter, data interface{}) {
	res := &Response{rw: w}

	res.Status = &ResponseStatus{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}

	res.Body = data

	res.WriteJSON()
}
