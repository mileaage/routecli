package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var ErrInvalidRoute = errors.New("invalid route")

func LoadRoutes(routes []string) *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range routes {
		log.Printf("Loading route for %s\n", route)
		mux.HandleFunc(strings.Join([]string{"GET ", route}, ""), func(w http.ResponseWriter, r *http.Request) {
			if route == "/" {
				route = "/home page!"
			}
			fmt.Fprintf(w, "Route for %s", route[1:])
		})
	}

	return mux
}

func AddRoute(route string) error {
	if route[0] != byte('/') {
		return ErrInvalidRoute
	}

	err := AddToRoutes(route)
	if err != nil {
		return err
	}

	return nil
}

func StartRoutes() error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	mux := LoadRoutes(config.Routes)

	// attach middlewares
	wrappedMux := AttachMiddlewares(mux, config.Middlewares)

	http.ListenAndServe(":8080", wrappedMux)
	return nil
}
