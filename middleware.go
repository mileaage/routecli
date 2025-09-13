package main

import (
	"log"
	"net/http"
	"time"
)

// going to make some experimental middleware
type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

// create the wrapper
func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

/*
@TODO:
  - Add auth
  - Add cors
*/
func AttachMiddlewares(handler http.Handler, middlewares []string) http.Handler {

	for _, middleware := range middlewares {
		switch middleware {
		case "logger":
			handler = NewLogger(handler)
		}
	}

	return handler
}
