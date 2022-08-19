package handler

import "net/http"

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	payload := r.PostForm

	if !payload.Has("jwt") { //HAS email and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must contain a jwt.\n"))
		return
	}

	//TODO: Call jwt validation logic

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account is authorized."))
}
