package main

import (
    "net/http"
    "encoding/json"
    "strings"
)
type Task struct {
    ID    string `json:"id"`
    Title string `json:"title"`
    IsCompleted bool `json:"isCompleted"`
}

var tasks = []Task{} 

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
    http.HandleFunc("/tasks/", tasksHandler)
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

    case "PUT":
        updateTask(w, r)

    case "DELETE":
        deleteTask(w, r)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}
func updateTask(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/tasks/")
    var updatedTask Task
    if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    for i, task := range tasks {
        if task.ID == id {
            tasks[i].Title = updatedTask.Title
            tasks[i].IsCompleted = updatedTask.IsCompleted
            json.NewEncoder(w).Encode(tasks[i])
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/tasks/")
    for i, task := range tasks {
        if task.ID == id {
            // tasks is updated and element i is deleted
            tasks = append(tasks[:i], tasks[i+1:]...)
            w.WriteHeader(http.StatusOK)
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}
