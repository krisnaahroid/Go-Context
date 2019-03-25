package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", StatusPage)
	mux.HandleFunc("/login", LoginPage)
	mux.HandleFunc("/logout", LogoutPage)

	ContextMux := AddContext(mux)

	log.Println("Starting on port 9090")
	log.Fatal(http.ListenAndServe(":9090", ContextMux))
}

func StatusPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This page will show the context username once the context is added \n"))

	if username := r.Context().Value("Username"); username != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello " + username.(string) + "\n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Logged In"))
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(1 * 1 * time.Hour)
	cookie := http.Cookie{Name: "username", Value: "ahroidlife@gmail.com", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{Name: "username", Value: "ahroidlife@gmail.com", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "-", r.RequestURI)
		cookie, _ := r.Cookie("username")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "Username", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
