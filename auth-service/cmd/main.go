package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Структура для представления пользователя
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Временное "хранилище" пользователей в памяти (заменится на БД)
var users = make(map[string]string)

// Обработчик для главной страницы
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Auth Service is working! Request: %s %s", r.Method, r.URL.Path)
}

// Обработчик для регистрации
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что username и password не пустые
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Проверяем, не существует ли уже пользователь
	if _, exists := users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Сохраняем пользователя в памяти (ВРЕМЕННО!)
	users[user.Username] = user.Password
	log.Printf("User registered: %s", user.Username)

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": user.Username,
	})
}

// Обработчик для логина
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем существование пользователя и пароль
	storedPassword, exists := users[user.Username]
	if !exists || storedPassword != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("User logged in: %s", user.Username)

	// Успешный ответ (пока без JWT)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login successful",
		"username": user.Username,
	})
}

func main() {
	// Регистрируем обработчики
	http.HandleFunc("/", homePage)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	port := "8080"
	log.Printf("Auth Service starting on port %s...", port)
	log.Printf("Endpoints:")
	log.Printf("  POST /register - Register new user")
	log.Printf("  POST /login - Login user")
	log.Printf("  GET  / - Health check")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
