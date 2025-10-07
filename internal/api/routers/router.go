package routers

import (
	"fmt"
	"net/http"
)

func MainRouter() *http.ServeMux {

	mux := http.NewServeMux()

	// eRouter := execsRouter()
	// tRouter := teachersRouter()
	// sRouter := studentsRouter()

	// sRouter.Handle("/", eRouter)
	// tRouter.Handle("/", sRouter)
	// return tRouter

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the API root 🚀")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK ✅")
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong 🏓")
	})

	return mux

}
