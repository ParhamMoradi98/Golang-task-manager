package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "reflect"
    "testing"
    "strings"
)
func resetTasks() {
    tasks = []Task{}  // Reset the global task list
}

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

// This test ensures that only unique IDs are posted
func TestPostTaskWithNonUniqueID(t *testing.T) {
    resetTasks()  

    // Add an initial task 
    tasks = append(tasks, Task{ID: "1", Title: "Initial Task"})

    // Add another task with the same ID, which should fail
    duplicateTask := Task{ID: "1", Title: "Duplicate Task"}
    jsonData, _ := json.Marshal(duplicateTask)
    req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        t.Fatal(err)  
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(tasksHandler)

    handler.ServeHTTP(rr, req)

    // Check that the API correctly responded with a 400 Bad Request
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
    }

    expectedErrorMessage := "Task with this ID already exists"
    actualMessage := strings.TrimSpace(rr.Body.String())  
    if actualMessage != expectedErrorMessage {
        t.Errorf("handler returned unexpected body: got '%v' want '%v'", actualMessage, expectedErrorMessage)
    }
}



func TestUpdateTask(t *testing.T) {
    resetTasks()

    // Create a new task via POST
    task := Task{ID: "1", Title: "Test Task", IsCompleted: false}
    jsonData, _ := json.Marshal(task)
    postReq, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
    if err != nil {
        t.Fatal(err)
    }
    postReq.Header.Set("Content-Type", "application/json")

    postRR := httptest.NewRecorder()
    handler := http.HandlerFunc(tasksHandler)

    handler.ServeHTTP(postRR, postReq)

    if status := postRR.Code; status != http.StatusOK {
        t.Errorf("POST handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Define the changes to the task
    updates := Task{Title: "New Title", IsCompleted: true}
    updateData, _ := json.Marshal(updates)
    putReq, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(updateData))
    if err != nil {
        t.Fatal(err)
    }
    putReq.Header.Set("Content-Type", "application/json")

    putRR := httptest.NewRecorder()

    handler.ServeHTTP(putRR, putReq)

    if status := putRR.Code; status != http.StatusOK {
        t.Errorf("PUT handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var updatedTask Task
    if err := json.Unmarshal(putRR.Body.Bytes(), &updatedTask); err != nil {
        t.Fatal(err)
    }
    // Was task updated correctly ?

    if updatedTask.Title != "New Title" || !updatedTask.IsCompleted {
        t.Errorf("handler did not update task correctly: got %+v, want Title=%s, IsCompleted=%t",
            updatedTask, "New Title", true)
    }
}
func TestDeleteTask(t *testing.T) {
    resetTasks()
   
    tasks = append(tasks, Task{ID: "1", Title: "Task to Delete", IsCompleted: false})

    // DELETE request
    req, err := http.NewRequest("DELETE", "/tasks/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(tasksHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Ensure the task is actually removed
    if len(tasks) != 0 {
        t.Errorf("Task was not deleted, tasks count = %d", len(tasks))
    }
}
