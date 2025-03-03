package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// TODO: Find a way to use it only on handfull of requests, not all of them
	sessMng := app.sessionManager.LoadAndSave

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// NOTE: Match all path route
	mux.Handle("GET /static/{filepath...}", http.StripPrefix("/static", fileServer))

	/*
		NOTE: There is one last bit of syntax. As we showed above, patterns ending in a slash,
		like /posts/, match all paths beginning with that string.
		To match only the path with the trailing slash, you can write /posts/{$}.
		That will match /posts/ but not /posts or /posts/234.
	*/
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Match everything else
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// NOTE: Middleware
	// [IN] (Log request) -> (Add Headers) -> (Serve mux)
	// [OUT] (Recover Panic)    <-			  (Serve mux)
	return app.recoverPanic(app.logRequest(sessMng(secureHeaders(mux))))
}
