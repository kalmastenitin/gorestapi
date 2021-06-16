package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err.Error())
		// make response message
		var response = ErrorResponse{
			ErrorMessage: err.Error(),
			StatusCode:   http.StatusInternalServerError,
		}
		// make return json response
		message, _ := json.Marshal(response)
		w.WriteHeader(response.StatusCode)
		w.Write(message)
		return
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		var response = ErrorResponse{
			ErrorMessage: err.Error(),
			StatusCode:   http.StatusInternalServerError,
		}
		// make return json response
		message, _ := json.Marshal(response)
		w.WriteHeader(response.StatusCode)
		w.Write(message)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// set headers
	w.Header().Set("Content-Type", "application/json")
	var user User
	// decdode body and save it to variable
	_ = json.NewDecoder(r.Body).Decode(&user)
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return
	}

	// insert to db
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err.Error())
		// make response message
		var response = ErrorResponse{
			ErrorMessage: err.Error(),
			StatusCode:   http.StatusInternalServerError,
		}
		// make return json response
		message, _ := json.Marshal(response)
		w.WriteHeader(response.StatusCode)
		w.Write(message)
		return
	}
	// return data on success
	json.NewEncoder(w).Encode(result)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&user)

	update := bson.D{
		{"$set", bson.D{
			{"firstname", user.FirstName},
			{"lastname", user.LastName},
			{"age", user.Age},
			{"email", user.Email},
		}},
	}
	fmt.Println(update)
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
	if err != nil {
		var response = ErrorResponse{
			ErrorMessage: err.Error(),
			StatusCode:   http.StatusInternalServerError,
		}
		// make return json response
		message, _ := json.Marshal(response)
		w.WriteHeader(response.StatusCode)
		w.Write(message)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		var response = ErrorResponse{
			ErrorMessage: err.Error(),
			StatusCode:   http.StatusInternalServerError,
		}
		// make return json response
		message, _ := json.Marshal(response)
		w.WriteHeader(response.StatusCode)
		w.Write(message)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)

}
