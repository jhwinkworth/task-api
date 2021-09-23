package main

import (
	"log"
	"net/http"
	"task-api/router"
	"task-api/server"
)

func main() {
	r := &router.Router{}

	r.HandleFunc("/tasks/pending", http.MethodGet, server.GetPendingTasks)
	r.HandleFunc("/tasks", http.MethodPost, server.AddTask)
	r.HandleFunc(`/tasks/\d`, http.MethodPut, server.UpdateTask)

	log.Fatal(http.ListenAndServe(":8080", r))
}
