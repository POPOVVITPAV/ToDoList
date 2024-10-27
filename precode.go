package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications, omitempty"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// эндпоинт для получения всех задач
func getTasks(rw http.ResponseWriter, _ *http.Request) {
	// сериализуем данные из мапы tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// в заголовок записываем тип контента в формате JSON
	rw.Header().Set("Content-Type", "application/json")
	//  статус OK
	rw.WriteHeader(http.StatusOK)
	//записываем сериализованные данные в тело ответа в формате JSON
	_, _ = rw.Write(resp)
}

// эндпоинт для отправки задачи на сервер
func postTask(rw http.ResponseWriter, req *http.Request) {

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body) // читаем поток данных
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//десериализация данных
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := tasks[task.ID]; ok {
		http.Error(rw, "Задача присутствует в списке", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

}

// эндпоинт для получения задачи по ID
func getTaskID(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(rw, "Задача не найдена", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func dellTaskID(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	taskId := chi.URLParam(req, "id")

	if _, ok := tasks[taskId]; !ok {
		http.Error(rw, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, taskId)
	rw.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// Обработчик для получения всех задач
	r.Get("/tasks", getTasks)
	//Обработчик для отправки задачи на сервер
	r.Post("/tasks", postTask)
	//Обработчик для получения задачи по ID
	r.Get("/task/{id}", getTaskID)
	//Обработчик удаления задачи по ID
	r.Delete("/task/{id}", dellTaskID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
