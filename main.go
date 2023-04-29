// Lee Panti
// Systems Programming Test 2
// 04/29/2023

package main

import (
	"log"
	"mime"
	"net/http"
)

/* ---------------------------------------------------------------- */

// demonstate middleware stacking by using two middlewares with basic functionality.

func middleWareA(next http.Handler) http.Handler {
	// simply logs out that middleware is executing
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running middleware A...")
		next.ServeHTTP(w, r)
		log.Print("Running middleware A returning...\n\n")
	})
}

func middleWareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// checks the request path and doesn't continue if it is a certain path.
		log.Println("Running middleware B...")
		if r.URL.Path == "/icons" {
			log.Println("middleware B failed. Cannot execute home handler.")
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Running middleware B returning...")
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("root route Handler was successfully called.")
	w.Write([]byte("middlewares successfully executed"))
}

/* ---------------------------------------------------------------- */

// applying middleware to enforce content-type to be "application/json"
func enforceJASONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("content-type")
		errMessage := "Middleware failed. Cannot continue."
		log.Println("Running enforceJSONHandler middleware...")

		if headerContentType != "" {
			mediaType, _, err := mime.ParseMediaType(headerContentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type", http.StatusBadRequest)
				log.Println(errMessage)
				return
			}
			if mediaType != "application/json" {
				http.Error(w, "content type must be 'application/json'", http.StatusUnsupportedMediaType)
				log.Println(errMessage)
				return
			}
		} else {
			log.Println(errMessage)
			log.Println("Content-Type header must be provided.")
			return
		}
		next.ServeHTTP(w, r)
		log.Print("Returning through enforceJSONHandler middleware... \n\n")
	})
}

func contentTypeHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("headers route Handler was successfully called.")
	w.Write([]byte("Middleware successfully Executed."))
}

/* ---------------------------------------------------------------- */

func main() {
	// router to handle the endpoints
	mux := http.NewServeMux()

	// nesting middlewares to see the order in which they run.
	homeHandler := http.HandlerFunc(home)
	mux.Handle("/", middleWareA(middleWareB(homeHandler)))

	// application of middleware to enforce content-type header
	contentTypeHandler := http.HandlerFunc(headersHandler)
	mux.Handle("/headers", contentTypeHeaders(enforceJASONHandler(contentTypeHandler)))

	// setting up the server
	log.Println("Starting server on http://localhost:9000...")
	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
