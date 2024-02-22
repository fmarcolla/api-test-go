package main

import (
	"log"
    "net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"gopkg.in/validator.v2"
)

type Address struct {
	City 	string `json:"city,omitempty" validate:"nonzero"`
	State 	string `json:"state,omitempty" validate:"nonzero"`
}

type Person struct {
	Id		int 	`json:"id,omitempty" validate:"nonzero"`
	Name 	string 	`json:"name,omitempty" validate:"nonzero"`
	Age 	int 	`json:"age,omitempty" validate:"nonzero, min=0"`
	Address Address `json:"address,omitempty" validate:"nonzero"`
}

var person []Person

func GetAllPerson(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(person)
}

func GetPersonById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, item := range person {

		if item.Id == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Person{})
}

func CreatePerson(w http.ResponseWriter, r *http.Request){
	var personNew Person
	err := json.NewDecoder(r.Body).Decode(&personNew)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	if errs := validator.Validate(personNew); errs != nil {
		http.Error(w, "Campos Inválidos", http.StatusBadRequest)
        return
	}

	person = append(person, personNew)

	json.NewEncoder(w).Encode(person)
}

func DeletePerson(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for index, item := range person {

		if item.Id == id {
			person = append(person[:index], person[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(person)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/contato", GetAllPerson).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPersonById).Methods("GET")
	router.HandleFunc("/contato", CreatePerson).Methods("POST")
	router.HandleFunc("/contato/{id}", DeletePerson).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}