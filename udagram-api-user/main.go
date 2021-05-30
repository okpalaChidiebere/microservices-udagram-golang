package main

import (
	"fmt"
	"net/http"
	"os"

	uh "user-service/controllers/v0/users"

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
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
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
	usersRouter := indexRouter.PathPrefix("/users").Subrouter().StrictSlash(true)

	authRouter := usersRouter.PathPrefix("/auth").Subrouter()

	// Define "Subrouter" routes using indexRouter
	indexRouter.Methods("GET").Path("/").HandlerFunc(indexRouterHandler)

	// Define "root" routes using r
	r.Methods("GET").Path("/").HandlerFunc(index)

	// Define "Subrouter" routes using usersRouter, prefix is /api/v0/users/...
	usersRouter.Methods("GET").Path("/{id}").HandlerFunc(uh.GetUserHandler)

	authRouter.Methods("GET").Path("").HandlerFunc(authIndex)                  // ../api/v0/users/auth it was a little bit odd that i did not specify the index route "/" in the Path() method. But it works :)
	authRouter.Methods("POST").Path("/").HandlerFunc(uh.RegisterUserHandler)   // ../api/v0/users/auth/
	authRouter.Methods("POST").Path("/login").HandlerFunc(uh.LoginUserHandler) // ../api/v0/users/auth/login
	authRouter.Methods("GET").Path("/verification").HandlerFunc(uh.Adapt(
		http.HandlerFunc(uh.VerificationHandler),
		uh.RequireAuthHandler(),
	).ServeHTTP) // ../api/v0/users/auth/verification
	authRouter.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}) //you actually have to define the OPTIONS method for CORS preflight

	//http.ListenAndServe(":"+port, r)
	http.Handle("/favicon.ico", http.NotFoundHandler()) //nice to have :)
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/api/v0/")
}

func indexRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "v0")
}

func authIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "auth")
}
