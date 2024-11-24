package main

import (
	"github.com/gorilla/csrf"
	"net/http"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	csrfToken := csrf.Token(r)
	w.Write([]byte(`
		<html>
			<form action="/api/tasks" method="POST">
				<input type="hidden" name="gorilla.csrf.Token" value="` + csrfToken + `">
				<input type="text" name="title" placeholder="Task Title">
				<textarea name="description" placeholder="Task Description"></textarea>
				<button type="submit">Create Task</button>
			</form>
		</html>
	`))
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Form submitted successfully!"))
}
