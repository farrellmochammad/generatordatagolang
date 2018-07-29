package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"math"
	"github.com/alok87/goutils/pkg/random"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strconv"
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
	Casofid string `json:"CasId"`
}

type status struct{
	Status int `json:Status`
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
	DataScan *DataScan `json:"DataScan"`
}

type DataScan struct{
	NamaUnsur string `json:"NamaUnsur"`
	CacahGross string `json:"CacahGross"`
	UrutanScan int `json:"Urutan_Scan"`
}

type statusofportal struct{
	Statusportal string `json:"Status"`
	Portalutama string `json:"PortalUtama"`
	Portalpendukung string `json:"PortalPendukung"`
}

type editdata struct{
	Casid string `json:"CasId"`
	Ipcas string `json:"IpCas"`
	Lokasi string `json:"Lokasi"`
	DataPortal *DataPortal `json:"DataPortal"`
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
	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"
	//connect to database
	db,err := sql.Open(dbDriver,dbPort + dbName)
  if err != nil{ fmt.Println(err.Error())}
  defer db.Close()
  regiscasdb(db)
	json.NewEncoder(w).Encode(regis)
}

func generatorRpm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rpm rpmdata
	var casid,alarmid,namaunsur,filename,scanportalid string
	var i float64
	var idfile int
		_ = json.NewDecoder(r.Body).Decode(&rpm)
	Datarpm := [][]float64{}
	i = 1
	floatdata := 2000.0
	for i < 20 {
		temp := i
		if(i<5 || i>15){
			if(math.Mod(i,2)==1 && i<5){
				//randomize := random.RangeInt(1,300,2)
				floatdata += 200
				row1 := []float64{temp,floatdata}
				floatdata += 200
				row2 := []float64{temp+0.5,floatdata}
				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			} else if (math.Mod(i,2)==0 && i<5){
				//randomize := random.RangeInt(1,300,2)
				floatdata -= 200
				row1 := []float64{temp,floatdata}
				floatdata -= 200
				row2 := []float64{temp+0.5,floatdata}
				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			} else if (math.Mod(i,2)==0 && i>15){
				//randomize := random.RangeInt(1,300,2)
				floatdata += 200
				row1 := []float64{temp,floatdata}
				floatdata += 200
				row2 := []float64{temp+0.5,floatdata}
				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			} else if (math.Mod(i,2)==1 && i>15){
				//randomize := random.RangeInt(1,300,2)
				floatdata -= 200
				row1 := []float64{temp,floatdata}
				floatdata -= 200
				row2 := []float64{temp+0.5,floatdata}
				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			}
			if i == 19{
				row3 := []float64{temp+1,2000}
				Datarpm = append(Datarpm,row3)
			}
		} else {
			if (i==15){
				floatdata = 2000
				row1 := []float64{temp,floatdata}
				floatdata += 200
				row2 := []float64{temp+0.5,floatdata}
				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			} else{
				randomize := random.RangeInt(1000,4000,1)
				floatdata += float64(randomize[0])
				randomize = random.RangeInt(1000,2500,1)

				operand := random.RangeInt(1,2,1)
				if operand[0] == 1 {
					floatdata += float64(randomize[0])
				} else {
					floatdata -= float64(randomize[0])
					for (floatdata<2000){
						floatdata += float64(randomize[0])
						floatdata -= float64((random.RangeInt(1000,2500,1))[0])
					}
				}
				row1 := []float64{temp, floatdata}

				operand = random.RangeInt(0,1,1)
				if operand[0] == 1 {
					floatdata += float64(randomize[0])
				} else {
					floatdata -= float64(randomize[0])
					for (floatdata<2000){
						floatdata += float64(randomize[0])
						floatdata -= float64((random.RangeInt(1000,2500,1))[0])
					}
				}
				row2 := []float64{temp+0.5, floatdata}

				Datarpm = append(Datarpm,row1)
				Datarpm = append(Datarpm,row2)
			}
		}
		i += 1
	}

	p,err := plot.New()
	if err != nil {panic(err.Error())}

	p.Title.Text = "RPM Data"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Energy"
	cacahgross := "["

	pts := make(plotter.XYs,39)
	j := 0
  total := 0.0
	biggest := 0.0
	for j <= 38 {
		cacahgross += "[" + strconv.FormatFloat(Datarpm[j][0],'f',-1,64) + " " + strconv.FormatFloat(Datarpm[j][1],'f',-1,64) +  "]"
		pts[j].X = Datarpm[j][0]
		pts[j].Y = Datarpm[j][1]
		total += Datarpm[j][1]
		if (Datarpm[j][1]>biggest){
			biggest = Datarpm[j][1]
		}
		j += 1
		if j != 39{
			cacahgross += ","
		}
	}
	cacahgross += "]"

	graph := plotutil.AddLinePoints(p,"Data",pts)
	if graph != nil {panic(err.Error())}


	if (biggest<10000){
		namaunsur = "Non radiation"
	} else if (biggest>=10000 && biggest<12500){
		namaunsur = "Adamantium"
	} else if (biggest>12500 && biggest<15000){
		namaunsur = "Cobalt"
	} else if (biggest>15000 && biggest<17500){
		namaunsur = "Celsium"
	} else if (biggest>=17500){
		namaunsur = "Potassium"
	}
  rpm.DataScan.NamaUnsur = namaunsur
	rpm.DataScan.CacahGross = cacahgross

	fmt.Println("Data dikirim ")
	fmt.Println("{")
	fmt.Println("'StartTime':",rpm.StartTime)
	fmt.Println("'Durasi':",rpm.Durasi)
	fmt.Println("'AlarmStatus':",rpm.AlarmStatus)
	fmt.Println("'ImageData':",rpm.ImageData)
	fmt.Println("'NoKontainer':",rpm.NoKontainer)
	fmt.Println("'AlarmId':",rpm.AlarmId)
	fmt.Println("'CasId':",rpm.CasId)
	fmt.Println("'TanggalBuat':",rpm.TanggalBuat)
	fmt.Println("'UsernameCas':",rpm.UsernameCas)
	fmt.Println("'DataScan': {")
	fmt.Println("		'NamaUnsur':",rpm.DataScan.NamaUnsur)
	fmt.Println("		'CacahGross':",rpm.DataScan.CacahGross)
	fmt.Println("		'Urutan_Scan':",rpm.DataScan.UrutanScan)
	fmt.Println("		}")
	fmt.Println("}")

	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

  db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()


	fkcasid := "select cas_id from cas where cas_id ='" + rpm.CasId + "'"
	querycasid := db.QueryRow(fkcasid).Scan(&casid)

	fkalarmid := "select alarm_id from alarm where alarm_id ='" + rpm.AlarmId + "'"
	queryalarmid := db.QueryRow(fkalarmid).Scan(&alarmid)

	if (querycasid == sql.ErrNoRows || queryalarmid == sql.ErrNoRows ){
		stat = append(stat,status{Status:0})
		json.NewEncoder(w).Encode(stat)
		stat = append(stat[:0], stat[1:]...)
	} else {
		fkmaxdatascanid := "select max(data_scan_id) from data_scan"
		querymaxdatascanid := db.QueryRow(fkmaxdatascanid).Scan(&idfile)
		if querymaxdatascanid != sql.ErrNoRows{
			filename = "grafikrpm/datascanid"+ strconv.Itoa(idfile+1) +".png"
			if graph := p.Save(4*vg.Inch, 4*vg.Inch, filename); graph != nil {
				panic(graph)
			}
		}


		stinsertscanportal := "insert into scan_portal(start_time,durasi,alarm_status,Image_data,No_kontainer,alarm_id,cas_id,Tgl_buat,user_name_cas) values ('" + rpm.StartTime + "','" + rpm.Durasi + "','" + rpm.AlarmStatus + "','" + rpm.ImageData + "','" + rpm.NoKontainer + "','" + rpm.AlarmId + "','" + rpm.CasId + "','" + rpm.TanggalBuat + "','" + rpm.UsernameCas + "')"
		insertscanportal,err := db.Query(stinsertscanportal)
		if err != nil { panic(err.Error()) }
	  defer insertscanportal.Close()

		fkmaxdatascanportalid := "select max(scan_portal_id) from scan_portal"
		querymaxdatascanportalid := db.QueryRow(fkmaxdatascanportalid).Scan(&scanportalid)
		if querymaxdatascanportalid == sql.ErrNoRows{fmt.Println("Something Wrong")}

		stinsertdatascan := "insert into data_scan(nama_unsur,cacah_gross,urutan_scan,scan_portal_id) values ('" + rpm.DataScan.NamaUnsur + "','" + rpm.DataScan.CacahGross + "','" + strconv.Itoa(rpm.DataScan.UrutanScan) + "','" + scanportalid + "')"
		insertdatascan,err := db.Query(stinsertdatascan)
		if err != nil { panic(err.Error()) }
	  defer insertdatascan.Close()

		stat = append(stat,status{Status:1})
		json.NewEncoder(w).Encode(stat)
		stat = append(stat[:0], stat[1:]...)
	}
}

func getRpm(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"
	//connect to database
	db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil{ fmt.Println(err.Error())}
	defer db.Close()
	rpmdb(db)
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
	var start_time,durasi,alarm_status,image_data,no_kontainer,alarm_id,cas_id,tgl_buat,usernamecas,namaunsur,cacahgross string
	var urutanscan int
	query,err := db.Query("select start_time,durasi,alarm_status,image_data,no_kontainer,alarm_id,cas_id,tgl_buat,username.USER_NAME_CAS,data_scan.nama_unsur,data_scan.cacah_gross,data_scan.urutan_scan from scan_portal inner join username on scan_portal.user_name_cas=username.USER_NAME_CAS inner join data_scan on scan_portal.scan_portal_id = data_scan.scan_portal_id")
	if err != nil { panic (err.Error())}
	defer query.Close()

	for query.Next(){
			data := query.Scan(&start_time,&durasi,&alarm_status,&image_data,&no_kontainer,&alarm_id,&cas_id,&tgl_buat,&usernamecas,&namaunsur,&cacahgross,&urutanscan)
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
				UsernameCas : usernamecas,
				DataScan: &DataScan{
						NamaUnsur: namaunsur,
						CacahGross: cacahgross,
						UrutanScan: urutanscan}})
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

func editCas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var edit editdata
	var temp,count int
	var casid string
	_ = json.NewDecoder(r.Body).Decode(&edit)


	dbDriver := "mysql"
	dbPort := "root@tcp(127.0.0.1:3306)/"
	dbName := "cas_db"

  db,err := sql.Open(dbDriver,dbPort + dbName)
	if err != nil { panic(err.Error()) }
	defer db.Close()

	for i, rune := range edit.Casid {
			if ( string(rune) == "-"){
				temp = i
			}
	}

  id := edit.Casid[temp+1:]

	sqlUpdate := "UPDATE cas SET ip_cas = '" + edit.Ipcas  + "', Lokasi = '" + edit.Lokasi + "' WHERE cas_id = " + id
	_, err = db.Exec(sqlUpdate)
	if err != nil {
  	panic(err)
	}

	sqlUpdate = "UPDATE portal SET serial_number = '" + edit.DataPortal.SerialNumber + "', jenis_portal ='" + edit.DataPortal.JenisPortal + "', tgl_pasang = '" + edit.DataPortal.TanggalPasang + "' WHERE cas_id = " + id
	_, err = db.Exec(sqlUpdate)
	if err != nil {
  	panic(err)
	}

	fkcasid := "select cas_id from cas where cas_id ='" + id + "'"
	querycasid := db.QueryRow(fkcasid).Scan(&casid)

	if (querycasid == sql.ErrNoRows){
		count = 0
	} else {
		count = 1
	}
	stat = append(stat,status{Status:count})
	json.NewEncoder(w).Encode(stat)
	stat = append(stat[:0], stat[1:]...)
}

func restartcas(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var casid,temp string
	id := 0
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
	if temp != "" {id = 1 }
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
		stat = append(stat,status{Status:0})
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

	// Route handles & endpoints
	r.HandleFunc("/registerCas", Regiscas).Methods("GET")
	r.HandleFunc("/getregisterCas/{ip}", getRegiscas).Methods("GET")
	r.HandleFunc("/generatorRpmData",generatorRpm).Methods("POST")
	r.HandleFunc("/getRpmData",getRpm).Methods("GET")
	r.HandleFunc("/insertdatacas",insertCas).Methods("POST")
	r.HandleFunc("/editdatacas",editCas).Methods("POST")
	r.HandleFunc("/restartcas/casid/{id}",restartcas).Methods("GET")
	r.HandleFunc("/cekstatuscas/casid/{id}",cekcas).Methods("GET")


	// Start server
	log.Fatal(http.ListenAndServe(":8004", r))
}
