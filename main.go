// Lee Panti
// Systems Programming Test 2
// 04/29/2023

package main

import (
	"log"
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
	log.Println("root route Handler successfully called.")
	w.Write([]byte("middlewares successfully executed"))
}

func main() {
	// router to handle the endpoints
	mux := http.NewServeMux()

	// nesting middlewares to see the order in which they run.
	homeHandler := http.HandlerFunc(home)
	mux.Handle("/", middleWareA(middleWareB(homeHandler)))

	// setting up the server
	log.Println("Starting server on http://localhost:9000...")
	err := http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
