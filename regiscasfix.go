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

// Registercas struct (Model)
type Registercas struct {
  Ipcas string `json:"IpCas"`
  Lokasi string `json:"Lokasi"`
	DataPortal *DataPortal `json:"DataPortal"`
}

type jointableregis struct{
	ipcas string
	lokasi string
	serialnumber string
	jenisportal string
	tanggalpasang string
}

// DataPortal struct
type DataPortal struct {
  SerialNumber string `json:"SerialNumber"`
	JenisPortal string `json:"JenisPortal"`
	TanggalPasang  string `json:"TanggalPasang"`
}

// Init regiscas var as a slice Registercas struct
var regis []Registercas

// Get all regiscas
func Regiscas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regis)
}

func getRegiscas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
	var stringip string `json:"CAS-ID"`
	var cas_id string

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

  db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	st := "select cas_id from cas where ip_cas = '" +  params["ip"] + "'"
	query := db.QueryRow(st)
	if err != nil { panic (err.Error())}
	stringip = query.Scan(&st)
	json.NewEncoder(w).encode(stringip)
}

func Queryregiscasdb(db *sql.DB){
	var rc jointableregis
	query,err := db.Query("select cas.ip_cas,cas.LOKASI,portal.SERIAL_NUMBER,portal.JENIS_PORTAL,portal.TGL_PASANG from cas INNER JOIN portal ON portal.CAS_ID=cas.CAS_ID")
  if err != nil { panic (err.Error())}
	defer query.Close()

	for query.Next(){
		data := query.Scan(&rc.ipcas,&rc.lokasi,&rc.serialnumber,&rc.jenisportal,&rc.tanggalpasang)
		if data != nil {panic(data.Error())}
    regis = append(regis, Registercas{
        Ipcas:rc.ipcas,
        Lokasi:rc.lokasi,
				DataPortal: &DataPortal{
						SerialNumber: rc.serialnumber,
						JenisPortal: rc.jenisportal,
						TanggalPasang: rc.tanggalpasang }})
	}
}

func Queryinserttodb(db *sql.DB, query string){
	insert, err := db.Prepare(query)
	if err != nil { panic(err.Error()) }
	defer insert.Close()
}


// Main function
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

  Queryregiscasdb(db)

	// Route handles & endpoints
	r.HandleFunc("/registerCas", Regiscas).Methods("GET")
	r.HandleFunc("/getregisterCas/{ip}", getRegiscas).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8004", r))
}
