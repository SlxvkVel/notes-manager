package main

import (
	"fmt"
	"log"
	"net/http"
)

// Обработчик для главной страницы
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Service is working!")
}

// Обработчик для регистрации
func register(w http.ResponseWriter, r *http.Request) {
	// Пока просто проверяем, что метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Register endpoint hit!") // Пока просто ответ для проверки
}

func main() {
	// Регистрируем обработчики для путей
	http.HandleFunc("/", homePage)
	http.HandleFunc("/register", register)

	// Запускаем сервер на порту 8080
	log.Println("Auth Service starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
