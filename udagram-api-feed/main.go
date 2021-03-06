package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	fh "feed-service/controllers/v0/feed"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		/*
			By default, Elastic Beanstalk configures the nginx proxy to forward
			requests to your application on port 5000. You can override the default
			port by setting the PORT system property to the port on which your main
			application listens.
		*/
		port = "8080"
	}

	r := mux.NewRouter()

	//CORS Should be restricted
	r.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", os.Getenv("URL"))
				w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
				w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,PUT,PATCH,POST,DELETE")
				next.ServeHTTP(w, r)
			})
		},
	)

	/*
		Create a "Subrouter" dedicated to /api which will use the PathPrefix
		more on nesting routes with gorrila mux here https://stackoverflow.com/questions/25107763/nested-gorilla-mux-router-does-not-work
		https://binx.io/blog/2018/11/27/go-gorilla/
	*/
	indexRouter := r.PathPrefix("/api/v0").Subrouter().StrictSlash(true)

	// This step is where we connect our "index" SubRouter to Feed SubRouter and Users SubRouter
	feedRouter := indexRouter.PathPrefix("/feed").Subrouter().StrictSlash(true)

	// Define "Subrouter" routes using indexRouter
	indexRouter.Methods("GET").Path("/").HandlerFunc(indexRouterHandler)

	// Define "root" routes using r
	r.Methods("GET").Path("/").HandlerFunc(index)
	r.Methods("GET").Path("/health").HandlerFunc(fh.CheckDbConnectionHandler)

	// Define "Subrouter" routes using feedRouter, prefix is /api/v0/feed/...
	authFeedRouter := feedRouter.PathPrefix("").Subrouter() //this subRouter under feedRouter
	/*
		What we are doing here is we are type casting custom type Adapter to a mux.Middleware type
		the root type for our Adapter is func(http.Handler) http.Handler which is thesame to Middleware type. This is why we are able to do this
		https://pkg.go.dev/github.com/gorilla/mux@v1.8.0#MiddlewareFunc
	*/
	authFeedRouter.Use(mux.MiddlewareFunc(fh.RequireAuthHandler())) //we make the subrouter from feedRouter to be protected

	//Define not protected feedRouter. The RequireAuthHandler middleare will not apply to them
	feedRouter.Methods("GET").Path("").HandlerFunc(fh.IndexHandler)
	feedRouter.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	//Define protected feedRouter
	authFeedRouter.Methods("POST").Path("").HandlerFunc(fh.CreateFeedItemHandler)
	authFeedRouter.Methods("GET").Path("/{id}").HandlerFunc(fh.GetFeedItemHandler)
	authFeedRouter.Methods("GET").Path("/signed-url/{fileName}").HandlerFunc(fh.GetGetSignedUrlHandler)

	//http.ListenAndServe(":"+port, r)
	http.Handle("/favicon.ico", http.NotFoundHandler()) //nice to have :)
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
	log.Printf(`server running %s`, os.Getenv("URL"))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/api/v0/")
}

func indexRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "v0")
}
