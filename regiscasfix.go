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

type DataPortal struct {
  SerialNumber string `json:"SerialNumber"`
	JenisPortal string `json:"JenisPortal"`
	TanggalPasang  string `json:"TanggalPasang"`
}

type idofcas struct{
	Casofid string `json:"CAS-ID`
}

type status struct{
	Status string `json:Status`
}

type rpmdata struct{
  StartTime string `json:"StartTime"`
	Durasi string `json:"Durasi"`
	AlarmStatus string `json:"AlarmStatus"`
	ImageData string `json:ImageData`
	NoKontainer string `json:NoKontainer`
	AlarmId string `json:"AlarmId"`
	CasId string `json:"CasId"`
	TanggalBuat string `json:"TanggalBuat"`
	UsernameCas string `json:"UsernameCas"`
}

type statusofportal struct{
	Statusportal string `json:"Status"`
	Portalutama string `json:"PortalUtama"`
	Portalpendukung string `json:"PortalPendukung"`
}

// Init regiscas var as a slice Registercas struct
var regis []Registercas
var cas []idofcas
var datarpm []rpmdata
var stat []status
var portal []statusofportal

// Get all regiscas
func Regiscas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(regis)
}

func getRpm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datarpm)
}

func getRegiscas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
	var casid,stringip,temp,id string

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

  db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	st := "select cas_id,ip_cas from cas"
	query,err := db.Query(st)
	if err != nil {panic(err.Error())}
	defer query.Close()

	for query.Next(){
		data := query.Scan(&casid,&stringip)
		if data != nil {panic(data.Error())}
		if (stringip==params["ip"]){
				temp = casid
		}
	}
	if temp != "" {id = "CAS-" + temp }
	cas = append(cas,idofcas{Casofid:id})
	json.NewEncoder(w).Encode(cas)
	cas = append(cas[:0], cas[1:]...)
}

func regiscasdb(db *sql.DB){
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

func rpmdb(db *sql.DB){
	var start_time,durasi,alarm_status,image_data,no_kontainer,alarm_id,cas_id,tgl_buat,usernamecas string
	query,err := db.Query("select start_time,durasi,alarm_status,image_data,no_kontainer,alarm.alarm_id,cas.cas_id,tgl_buat,username.USER_NAME_CAS from scan_portal  inner join alarm on scan_portal.alarm_id=alarm.alarm_id  inner join cas on scan_portal.cas_id=cas.cas_id inner join username on scan_portal.user_name_cas=username.USER_NAME_CAS")
	if err != nil { panic (err.Error())}
	defer query.Close()

	for query.Next(){
			data := query.Scan(&start_time,&durasi,&alarm_status,&image_data,&no_kontainer,&alarm_id,&cas_id,&tgl_buat,&usernamecas)
			if data != nil { panic (err.Error())}
			datarpm = append(datarpm,rpmdata{
				StartTime : start_time,
				Durasi : durasi,
				AlarmStatus : alarm_status,
				ImageData : image_data,
				NoKontainer : no_kontainer,
				AlarmId : alarm_id,
				CasId : cas_id,
				TanggalBuat : tgl_buat,
				UsernameCas : usernamecas})
	}
}

func insertCas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var rgs Registercas
	var casid string
	_ = json.NewDecoder(r.Body).Decode(&rgs)

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

  db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	stinsertcas := "insert into cas(ip_cas,lokasi) values ('" + rgs.Ipcas + "','" + rgs.Lokasi + "')"
	insertcas,err := db.Query(stinsertcas)
	if err != nil { panic(err.Error()) }
  defer insertcas.Close()

	fkcasid := "select cas_id from cas where ip_cas ='" + rgs.Ipcas + "'"
	querycasid := db.QueryRow(fkcasid).Scan(&casid)

	switch {
	case querycasid == sql.ErrNoRows:
        log.Printf("No user with that ID.")
	case querycasid != nil:
        log.Fatal(err)
	default:
        fmt.Printf("CAS-ID : %s\n", casid)
  }

	stinsertportal := "insert into portal(serial_number,jenis_portal,tgl_pasang,cas_id) values ('" + rgs.DataPortal.SerialNumber + "','" + rgs.DataPortal.JenisPortal + "','" + rgs.DataPortal.TanggalPasang + "','" + casid + "')"
	insertportal,err := db.Query(stinsertportal)
	if err != nil { panic(err.Error()) }
	defer insertportal.Close()

	cas = append(cas,idofcas{Casofid:casid})
	json.NewEncoder(w).Encode(cas)
	cas = append(cas[:0], cas[1:]...)
	regis = append(regis,rgs)
}

func restartcas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var casid,temp string
	id := "0"
  params := mux.Vars(r)

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

	db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	st := "select cas_id from cas"
	query,err := db.Query(st)
	if err != nil {panic(err.Error())}
	defer query.Close()

	for query.Next(){
		data := query.Scan(&casid)
		if data != nil {panic(data.Error())}
		if (casid==params["id"]){
				temp = casid
		}
	}
	if temp != "" {id = "1" }
	stat = append(stat,status{Status:id})
	json.NewEncoder(w).Encode(stat)
	stat = append(stat[:0], stat[1:]...)
}

func cekcas(w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")
	var jenisportal,statusportal,st,portalutama,portalpendukung string
	params := mux.Vars(r)

	portalutama = "false"
	portalpendukung = "false"

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

	db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	st = "select jenis_portal,status_portal from portal where cas_id = " + params["id"]
	query,err := db.Query(st)
	if err != nil {panic(err.Error())}
	defer query.Close()

	cekquery := db.QueryRow(st)
	data := cekquery.Scan(&jenisportal,&statusportal)
	if (data != nil){
		stat = append(stat,status{Status:"0"})
		json.NewEncoder(w).Encode(stat)
		stat = append(stat[:0], stat[1:]...)
	}else{
		for query.Next(){
			data := query.Scan(&jenisportal,&statusportal)
			if data != nil {panic(data.Error())}
			if (jenisportal=="PORTAL UTAMA" && statusportal == "1"){
				portalutama = "true"
			} else if (jenisportal=="PORTAL PENDUKUNG" && statusportal == "1"){
				portalpendukung = "true"
			}
		}
		portal = append(portal,statusofportal{
				Statusportal:"1",
				Portalutama : portalutama,
				Portalpendukung : portalpendukung})
		json.NewEncoder(w).Encode(portal)
		portal = append(portal[:0], portal[1:]...)
	}
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

  regiscasdb(db)
	rpmdb(db)

	// Route handles & endpoints
	r.HandleFunc("/registerCas", Regiscas).Methods("GET")
	r.HandleFunc("/getregisterCas/{ip}", getRegiscas).Methods("GET")
	r.HandleFunc("/getrpmData",getRpm).Methods("GET")
	r.HandleFunc("/insertdatacas",insertCas).Methods("POST")
	r.HandleFunc("/restartcas/casid/{id}",restartcas).Methods("GET")
	r.HandleFunc("/cekstatuscas/casid/{id}",cekcas).Methods("GET")



	// Start server
	log.Fatal(http.ListenAndServe(":8004", r))
}
