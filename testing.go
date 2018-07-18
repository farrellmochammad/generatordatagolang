package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type idofcas struct{
  Casofid string `json:casofid`
}

var cas []idofcas

func Regiscas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cas)
}

func main() {

	// Init router
	r := mux.NewRouter()
	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"
	//connect to database
	db,err := sql.Open(dbDriver,dbPort + dbName)
  if err != nil{ fmt.Println(err.Error())}
  defer db.Close()

  var ids string
  query,err := db.Query("select cas_id from cas")
  if err != nil {panic(err.Error())}
  defer query.Close()

  for query.Next(){
    data := query.Scan(&ids)
    if data != nil {(panic(data.Error()))}
    cas = append(cas,idofcas{Casofid:ids})
  }

  cas = append(cas[:0], cas[0+1:]...)
	// Route handles & endpoints
	r.HandleFunc("/registerCas", Regiscas).Methods("GET")


	// Start server
	log.Fatal(http.ListenAndServe(":8010", r))
}
