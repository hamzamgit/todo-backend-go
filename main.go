package main

import (
	"fmt"
	"os"
	"time"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// struct for every Task
type Task struct {
	gorm.Model
	createdAt time.Time
	Task      string
	Completed bool
}

var db *gorm.DB
var err error

// this will makw API routes
func RouteHandler() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/tasks/list", GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", GetTask).Methods("GET")
	router.HandleFunc("/task/gen", GenerateTask).Methods("POST")
	router.HandleFunc("/task/del/{id}", DeleteTask).Methods("DELETE")
	return router
}

func main() {
	// loading environmental variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbNmae := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbNmae, password, dbPort)

	// Opening connection to database
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	// Close connection to database when the main function finishes
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// delete(Task.)
	db.AutoMigrate(&Task{})

	// making API routes
	router := RouteHandler()

	options := cors.Options{
		// 		AllowedOrigins: []string{},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		Debug:          false,
	}
	handler := cors.New(options).Handler(router)
	log.Fatal(http.ListenAndServe(":8081", handler))
}

// Create Task
func GenerateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)
	createdTask := db.Create(&newTask)
	err = createdTask.Error
	if err != nil {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			return
		}
	} else {
		json.NewEncoder(w).Encode(&newTask)
	}
	fmt.Println(time.Now(), "Create Task")
}

// Get an Array of tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	db.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
	w.WriteHeader(200)
	fmt.Println(time.Now(), "Get tasks")
}

// Get a Specific task based on ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	params := mux.Vars(r)
	db.First(&task, params["id"])
	json.NewEncoder(w).Encode(&task)
	fmt.Printf("%s Get task %s \n", time.Now(), params["id"])
}

// Delete a specific task based on ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	params := mux.Vars(r)
	db.First(&task, params["id"])
	db.Delete(&task)
	var tasks []Task
	db.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
	fmt.Printf("%s Del task %s \n", time.Now(), params["id"])
}

