// applic.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

// App.
type App struct {
	Router *chi.Mux
	DB     *sql.DB
}

// Initialize create a database connection
func (a *App) Initialize() {
	fmt.Println(">< Initializing... ><")
	connectionString := fmt.Sprintf("postgres://postgres:password@postgres/postgres?sslmode=disable") //DB configuration lies in dbConfig file

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(">< DB initialized ><")
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = a.DB.Ping(); err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	a.Router = chi.NewRouter() //Creating and initializing Router
	a.initializeRoutes()
	fmt.Println(">< Routes initialized ><")

}

// Run simply starts the application.
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

//Setting paths
func (a *App) initializeRoutes() {
	a.Router.Get("/employees", a.getEmployees)
	a.Router.Post("/employee", a.createEmployee)
	a.Router.Get("/employee/{uid:[0-9]+}", a.getEmployee)
	a.Router.Put("/employee/{uid:[0-9]+}", a.updateEmployee)
	a.Router.Delete("/employee/{uid:[0-9]+}", a.deleteEmployee)
}

func (a *App) getEmployee(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	p := employee{Uid: uid}
	if err := p.getEmployee(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "employee not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	fmt.Println(">< Reading employee`s data  is successful ><")

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getEmployees(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	employees, err := getEmployees(a.DB, start, count)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	/*employee1 := employee{}
	for _, value := range employees {
		employee1 = value

		err := json.NewDecoder(r.Body).Decode(&employee1)
		if err != nil{
			panic(err)
		}

		employeeJson, err := json.Marshal(employee1)
		if err != nil{
			panic(err)
		}

		//Set Content-Type header so that clients will know how to read response
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		//Write json response back to response
		w.Write(employeeJson)
	}*/

	fmt.Println(">< Reading employees` data is successful ><")

	respondWithJSON(w, http.StatusOK, employees)

}

func (a *App) createEmployee(w http.ResponseWriter, r *http.Request) {
	var p employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(">< Employee`s profile is created ><")
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var p employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.Uid = id

	if err := p.updateEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println(">< Employee data ( uid - " + string(id) + " ) are updated. ><")
	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "uid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	p := employee{Uid: id}
	if err := p.deleteEmployee(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println(">< Employee data ( uid - " + string(id) + " ) are deleted. ><")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
