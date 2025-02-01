package main

import (
	"encoding/json"
	"firstRest/db"
	"firstRest/orm"
	"fmt"
	"github.com/gorilla/mux"
	"log"
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

func showTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var messages []orm.Message
		db.DB.Find(&messages)
		json.NewEncoder(w).Encode(messages)
	} else {
		fmt.Fprintln(w, "Поддерживается только метод GET")
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		params := mux.Vars(r)
		id := params["id"]

		var message orm.Message
		result := db.DB.First(&message, id)

		if result.Error != nil {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			return
		}

		var updatedMessage orm.Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&updatedMessage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при разборе JSON: %v", err), http.StatusBadRequest)
			return
		}

		message.Task = updatedMessage.Task
		message.IsDone = updatedMessage.IsDone

		result = db.DB.Save(&message)
		if result.Error != nil {
			http.Error(w, "Не удалось обновить запись", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		params := mux.Vars(r)
		id := params["id"]

		log.Printf("Полученный ID: %s\n", id)

		var message orm.Message
		result := db.DB.Unscoped().First(&message, id)

		if result.Error != nil {
			log.Printf("Запись не найдена: %v\n", result.Error)
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			return
		}

		log.Printf("Найденная запись: %+v\n", message)

		result = db.DB.Delete(&message)
		if result.Error != nil {
			log.Printf("Не удалось удалить запись: %v\n", result.Error)
			http.Error(w, "Не удалось удалить запись", http.StatusInternalServerError)
			return
		}

		log.Printf("Запись успешно удалена\n")
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {

	db.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/showTask", showTaskHandler).Methods("Get")
	router.HandleFunc("/addTask", addTaskHandler).Methods("Post")
	router.HandleFunc("/updateTask/{id}", updateTaskHandler).Methods("PUT")
	router.HandleFunc("/deleteTask/{id}", deleteTaskHandler).Methods("DELETE")
	http.ListenAndServe("localhost:8080", router)
}
