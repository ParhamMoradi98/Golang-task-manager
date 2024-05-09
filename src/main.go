package main

import (
    "net/http"
    "encoding/json"
)
type Task struct {
    ID    string `json:"id"`
    Title string `json:"title"`
}

var tasks = []Task{} // a simple in-memory store for tasks

// Make sure that IDs are unique
func taskExists(id string) bool {
    for _, t := range tasks {
        if t.ID == id {
            return true
        }
    }
    return false
}
func main() {
    http.HandleFunc("/tasks", tasksHandler)
    http.HandleFunc("/tasks/", taskHandler)
    http.ListenAndServe(":8080", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        json.NewEncoder(w).Encode(tasks)
    case "POST":
        var task Task
        json.NewDecoder(r.Body).Decode(&task)
		if taskExists(task.ID) {
            http.Error(w, "Task with this ID already exists", http.StatusBadRequest)
            return
        }
        tasks = append(tasks, task)
        json.NewEncoder(w).Encode(task)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
    // implementation for handling specific task based on ID
}