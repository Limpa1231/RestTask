package main

import (
	"encoding/json"
	"firstRest/db"
	"firstRest/orm"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var task string

type requestBody struct {
	Message string `json:"message"`
}

func addTasksHandler(w http.ResponseWriter, r *http.Request) {
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

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Кодируем данные в JSON и отправляем клиенту
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

		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")

		// Кодируем данные в JSON и отправляем клиенту
		json.NewEncoder(w).Encode(messages)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {

	db.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/showTasks", showTasksHandler).Methods("Get")
	router.HandleFunc("api/addTasks", addTasksHandler).Methods("Post")
	http.ListenAndServe("localhost:8080", router)
}
