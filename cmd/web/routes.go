package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// NOTE: Middleware
	// [IN] (Log request) -> (Add Headers) -> (Serve mux)
	// [OUT] (Recover Panic)    <-			  (Serve mux)
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
