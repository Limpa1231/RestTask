package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var counter int
var task string

type requestBody struct {
	Message string `json:"message"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var requestBody requestBody

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при разборе JSON: %v", err), http.StatusBadRequest)
			return
		}

		task = requestBody.Message
		fmt.Fprintln(w, "Задача успешно сохранена:", task)
	}
}

func showTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Задачи:", task)
	} else {
		fmt.Fprintln(w, "Поддерживается только метод GET")
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Counter равен", strconv.Itoa(counter))
	} else {
		fmt.Fprintln(w, "Поддерживается только метод GET")
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		counter++
		fmt.Fprintln(w, "Counter увеличен на 1")
	} else {
		fmt.Fprintln(w, "Поддерживается только метод POST")
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", HelloHandler).Methods("Get")
	router.HandleFunc("/get", GetHandler).Methods("Get")
	router.HandleFunc("/showTask", showTaskHandler).Methods("Get")
	router.HandleFunc("/post", PostHandler).Methods("Post")
	router.HandleFunc("/addTask", addTaskHandler).Methods("Post")
	http.ListenAndServe("localhost:8080", router)
}
