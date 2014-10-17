package hut

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
)

func (s *Service) ErrorReply(err error, w http.ResponseWriter) {
	log.Println(err)
	debug.PrintStack()
	if s.Env.InProd() {
		// If we're in prod, we don't want to reveal internal server details
		s.HttpErrorReply(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		stack := make([]byte, 4096)
		b := bytes.NewBuffer(stack)
		bytesWritten := runtime.Stack(b.Bytes(), false)
		b.Truncate(bytesWritten)
		message := strings.Join([]string{err.Error(), b.String()}, "\n")
		s.HttpErrorReply(w, message, http.StatusInternalServerError)
	}
}

func (s *Service) HttpErrorReply(w http.ResponseWriter, message string, code int) {
	reply := map[string]interface{}{
		"status":  "error",
		"code":    code,
		"message": message,
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := encoder.Encode(reply)
	if err != nil {
		s.ErrorReply(err, w)
	}
}

func (s *Service) NotFound(w http.ResponseWriter) {
	s.HttpErrorReply(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Service) Reply(response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonEncoder := json.NewEncoder(w)
	if err := jsonEncoder.Encode(response); err != nil {
		s.ErrorReply(err, w)
		return
	}
	return
}

func (s *Service) SuccessReply(metadata map[string]interface{}, w http.ResponseWriter) {
	rep := map[string]interface{}{"status": "OK"}
	if metadata != nil {
		for k, v := range metadata {
			rep[k] = v
		}
	}
	s.Reply(rep, w)
}
