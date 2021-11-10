package server

import "net/http"

func runHealthAccepter() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm alive."))
	})
	return http.ListenAndServe("0.0.0.0:6688", nil)
}
