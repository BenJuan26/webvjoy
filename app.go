package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tajtiattila/vjoy"
)

var joy *vjoy.Device

func main() {
	var err error
	joy, err = vjoy.Acquire(1)
	if err != nil {
		panic("Couldn't find vJoy device 1. Is it configured?")
	}

	r := mux.NewRouter()
	r.HandleFunc("/button/{button}", handleButton)
	r.HandleFunc("/", handleIndex)
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./static/")))

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func handleButton(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	buttonString, ok := vars["button"]
	if !ok {
		http.Error(w, "No button specified", http.StatusBadRequest)
	}

	button, err := strconv.ParseUint(buttonString, 10, 64)
	if err != nil {
		http.Error(w, "Couldn't parse button number", http.StatusBadRequest)
	}

	go pressButton(uint(button))

	fmt.Printf("Button %d pressed\n", button)
}

func pressButton(button uint) {
	joy.Button(button).Set(true)
	if err := joy.Update(); err != nil {
		fmt.Printf("Got error pressing button %d: %s\n", button, err.Error())
		return
	}
	time.Sleep(500 * time.Millisecond)

	joy.Button(button).Set(false)
	if err := joy.Update(); err != nil {
		fmt.Printf("Got error releasing button %d: %s\n", button, err.Error())
	}
}
