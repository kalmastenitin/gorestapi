package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	FirstName   string    `json:"firstname" validate:"alpha"`
	LastName    string    `json:"lastname" validate:"alpha"`
	Username    string    `json:"username"`
	Age         uint8     `json:"age"`
	Email       string    `json:"email" validate:"required,email"`
	Status      bool      `json:"status"`
	DateCreated time.Time `json:"datecreated"`
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

var validate *validator.Validate
var collection = connectDB()

func connectDB() *mongo.Collection {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	fmt.Println("database is now ready...")

	collection := client.Database("userinfo").Collection("user")
	return collection

}

func main() {

	validate = validator.New()
	r := mux.NewRouter()
	r.HandleFunc("/api/users", helpers.getUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", helpers.getUser).Methods("GET")
	r.HandleFunc("/api/user/{id}", helpers.updateUser).Methods("PUT")
	r.HandleFunc("/api/user/{id}", helpers.deleteUser).Methods("DELETE")
	r.HandleFunc("/api/user", helpers.createUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
