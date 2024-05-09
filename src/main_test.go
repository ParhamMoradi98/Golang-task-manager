package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"
)

// TestGetTasks tests fetching of tasks
func TestGetTasks(t *testing.T) {
    req, err := http.NewRequest("GET", "/tasks", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(tasksHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Decode the response body
    var tasks []Task
    err = json.Unmarshal(rr.Body.Bytes(), &tasks)
    if err != nil {
        t.Fatal(err)
    }

    expected := []Task{} // initially no tasks
    if !reflect.DeepEqual(tasks, expected) {
        t.Errorf("handler returned unexpected body: got %v want %v", tasks, expected)
    }
}

// TestPostTasks tests adding a new task
func TestPostTasks(t *testing.T) {
    task := Task{ID: "1", Title: "Test Task"}
    jsonData, _ := json.Marshal(task)
    req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(tasksHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Decode the response to check if the task was added correctly
    var returnedTask Task
    err = json.Unmarshal(rr.Body.Bytes(), &returnedTask)
    if err != nil {
        t.Fatal(err)
    }

    if !reflect.DeepEqual(returnedTask, task) {
        t.Errorf("handler returned unexpected body: got %v want %v", returnedTask, task)
    }
}
