package main

import (
	"errors"
	"log"
	"net/http"
)

var ErrInvalidRoute = errors.New("invalid route")

func LoadRoutes(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range routes {
		// avoid closures

		log.Printf("Loading route for %s with template %s\n", route.Path, route.Template)

		// Create a handler
		handler := createRouteHandler(route)
		mux.HandleFunc(route.Path, handler)
	}

	return mux
}

// create a handler for closures
func createRouteHandler(route Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check exact url path instead
		if r.URL.Path != route.Path {
			// root handler
			if route.Path == "/" {
				http.NotFound(w, r)
				return
			}
			// not found
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Serving - Path: %s (requested: %s) Template: %s\n", route.Path, r.URL.Path, route.Template)

		// get the template file
		templateFile, err := GetTemplateFile(route.Template)
		if err != nil {
			log.Printf("Template file error for %s: %v\n", route.Template, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get template data
		templateData, err := GetTemplateData(route.Template)
		if err != nil {
			log.Printf("Template data error for %s: %v\n", route.Template, err)
			// Don't error out if template data doesn't exist, just use empty map
			templateData = make(map[string]any)
		}

		err = templateFile.Execute(w, templateData)
		if err != nil {
			log.Printf("Template execution error: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddRoute(route string) error {
	if route[0] != byte('/') {
		return ErrInvalidRoute
	}

	err := AddToRoutes(route, "default")
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

	log.Printf("Loaded config with routes: %+v\n", config.Routes)

	mux := LoadRoutes(config.Routes)

	// return 404 for  routes
	notFoundMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := mux.Handler(r)
		log.Printf("Request %s -> Handler pattern: '%s'\n", r.URL.Path, pattern)

		if pattern == "" {
			log.Printf("No pattern found for %s, returning 404\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})

	// attach middlewares
	wrappedMux := AttachMiddlewares(notFoundMux, config.Middlewares)

	log.Println("Server starting on :8080")
	return http.ListenAndServe(":8080", wrappedMux)
}
