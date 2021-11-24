package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	port = "8080"
	db   []Pizza
)

func init() {
	pizza1 := Pizza{
		ID:       1,
		Diameter: 22,
		Price:    500.50,
		Title:    "Pepperoni",
	}

	pizza2 := Pizza{
		ID:       2,
		Diameter: 25,
		Price:    650.23,
		Title:    "BBQ",
	}
	pizza3 := Pizza{
		ID:       3,
		Diameter: 22,
		Price:    450,
		Title:    "Margarita",
	}

	db = append(db, pizza1, pizza2, pizza3)
}

type Pizza struct {
	ID       int     `json:"id"`
	Diameter int     `json:"diameter"`
	Price    float64 `json:"price"`
	Title    string  `json:"title"`
}

// FindPizzaById Helper function for Pizza model
func FindPizzaById(id int) (Pizza, bool) {
	var pizza Pizza
	var found bool
	for _, p := range db {
		if p.ID == id {
			pizza = p
			found = true
			break
		}
	}
	return pizza, found
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func GetAllPizzas(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	log.Println("Get infos about all pizzas in database")
	writer.WriteHeader(200)            // StatusCode for the request
	json.NewEncoder(writer).Encode(db) // Serialization + add to writer
}

func GetPizzaById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request) // {"id" : "12"} []map
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Client tries to use invalid ID param:", err)
		msg := ErrorMessage{Message: "Do not use ID which is no supported for int casting"}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	log.Println("Trying to send to client pizza with id #:", id)
	pizza, ok := FindPizzaById(id)
	if ok {
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(pizza)
	} else {
		msg := ErrorMessage{Message: "Pizza with this id does not exist"}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
	}
}

func main() {
	log.Println("Trying to start REST API pizza!")
	router := mux.NewRouter()
	router.HandleFunc("/pizzas", GetAllPizzas).Methods("GET")
	router.HandleFunc("/pizza/{id}", GetPizzaById).Methods("GET")
	log.Println("Router has been successfully configured! Let's go!")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
