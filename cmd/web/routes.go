package main

import (
	"net/http"

	"github.com/markaya/snippetbox/ui"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	// TODO: Find a way to use it only on handfull of requests, not all of them
	dynamic := func(handler http.Handler) http.Handler {
		return app.sessionManager.LoadAndSave(app.authenticate(handler))
	}

	// NOTE: Protected routes
	protected := func(handler http.Handler) http.Handler {
		return dynamic(app.requireAuthentication(handler))
	}

	fileServer := http.FileServer(http.FS(ui.Files))

	// NOTE: Match all path route
	//
	// NOTE: When using embeded files we do not need to strip prefix
	mux.Handle("GET /static/{filepath...}", fileServer)

	/*
		NOTE: There is one last bit of syntax. As we showed above, patterns ending in a slash,
		like /posts/, match all paths beginning with that string.
		To match only the path with the trailing slash, you can write /posts/{$}.
		That will match /posts/ but not /posts or /posts/234.
	*/

	mux.HandleFunc("GET /ping", app.ping)

	// NOTE: Regular Session
	mux.Handle("GET /{$}", dynamic(http.HandlerFunc(app.home)))
	mux.Handle("GET /about", dynamic(http.HandlerFunc(app.aboutView)))
	mux.Handle("GET /snippet/view/{id}", dynamic(http.HandlerFunc(app.snippetView)))
	mux.Handle("GET /user/signup", dynamic(http.HandlerFunc(app.userSignup)))
	mux.Handle("POST /user/signup", dynamic(http.HandlerFunc(app.userSignupPost)))
	mux.Handle("GET /user/login", dynamic(http.HandlerFunc(app.userLogin)))
	mux.Handle("POST /user/login", dynamic(http.HandlerFunc(app.userLoginPost)))

	// NOTE: Auth Session
	mux.Handle("GET /snippet/create", protected(dynamic(http.HandlerFunc(app.snippetCreate))))
	mux.Handle("POST /snippet/create", protected(dynamic(http.HandlerFunc(app.snippetCreatePost))))
	mux.Handle("POST /user/logout", protected(dynamic(http.HandlerFunc(app.userLogoutPost))))
	mux.Handle("GET /account/view/", protected(dynamic(http.HandlerFunc(app.accountView))))
	mux.Handle("GET /account/password/update", protected(dynamic(http.HandlerFunc(app.accountPasswordUpdate))))
	mux.Handle("POST /account/password/update", protected(dynamic(http.HandlerFunc(app.accountPasswordUpdatePost))))

	// NOTE: How to handle a specific path with its own Middleware, not to be used
	// on all others
	// mux.Handle("/", sessMng(http.HandlerFunc(app.userLogin)))

	// Match everything else
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// NOTE: Middleware
	// [IN] (Log request) -> (Add Headers) -> (Serve mux)
	// [OUT] (Recover Panic)    <-			  (Serve mux)
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
