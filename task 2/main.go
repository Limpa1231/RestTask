package main

import (
	"encoding/json"
	"firstRest/db"
	"firstRest/orm"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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
		// Создаем новую запись в базе данных
		message := orm.Message{
			Task:   task,
			IsDone: false,
		}

		result := db.DB.Create(&message)
		if result.Error != nil {
			fmt.Println("Ошибка при сохранении записи в базу данных:", result.Error)
			http.Error(w, "Не удалось сохранить задачу", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

}

func showTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var messages []orm.Message
		db.DB.Find(&messages)
		json.NewEncoder(w).Encode(messages)
	} else {
		fmt.Fprintln(w, "Поддерживается только метод GET")
	}
}

func main() {

	db.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/showTasks", showTasksHandler).Methods("Get")
	router.HandleFunc("api/addTask", addTaskHandler).Methods("Post")
	http.ListenAndServe("localhost:8080", router)
}
