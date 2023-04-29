// Lee Panti
// Systems Programming Test 2
// 04/29/2023

package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
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

// handler functions for third party middlewares

// handler function after third party authentication middleware
func landingPage(w http.ResponseWriter, r *http.Request) {
	log.Print("landing page route Handler successfully called. \n\n")
	w.Write([]byte("Middleware successfully Executed."))
}

// handler function after third party request logging middleware
func loggingFile(w http.ResponseWriter, r *http.Request) {
	log.Print("logging file route Handler successfully called. \n\n")
	w.Write([]byte("Middleware successfully Executed."))
}

// constructor function for logging middleware
func newLoggingHandler(file io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(file, h)
	}
}

func main() {
	// router to handle the endpoints
	mux := http.NewServeMux()

	// nesting middlewares to see the order in which they run.
	homeHandler := http.HandlerFunc(home)
	mux.Handle("/", middleWareA(middleWareB(homeHandler)))

	// application of middleware to enforce content-type header
	contentTypeHandler := http.HandlerFunc(headersHandler)
	mux.Handle("/headers", contentTypeHeaders(enforceJASONHandler(contentTypeHandler)))

	// third party middleware for authentication before calling the handler
	authMiddleware := httpauth.SimpleBasicAuth("lee", "pass")
	landingPageHandler := http.HandlerFunc(landingPage)
	mux.Handle("/signup", authMiddleware(landingPageHandler))

	// third party middleware for request logging before calling the handler
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	loggingHandler := http.HandlerFunc(loggingFile)
	mux.Handle("/log", handlers.LoggingHandler(logFile, loggingHandler))

	// logging requests with constructor function
	loggingFileHandler := http.HandlerFunc(loggingFile)
	handleLogs := newLoggingHandler(logFile)
	mux.Handle("/constructor", handleLogs(loggingFileHandler))

	// easily chaining middleware with alice package.
	middleWareChain := alice.New(middleWareA, middleWareB)

	mux.Handle("/easychain", middleWareChain.Then(http.HandlerFunc(home)))

	// setting up the server
	log.Println("Starting server on http://localhost:9000...")
	err = http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
