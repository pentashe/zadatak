package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"zadatak/saxpy"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

func isAuth(w http.ResponseWriter, r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	if ok && verifyUserPass(user, pass) {
		return true
	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
}

func returnJMBAG(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "JMBAG: ", c.JMBAG)
	}
}

func returnSum(w http.ResponseWriter, r *http.Request) {
	if isAuth(w, r) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			queryParams := r.URL.Query()
			aParam := queryParams.Get("a")
			bParam := queryParams.Get("b")

			a, err := strconv.Atoi(aParam)
			if err != nil {
				http.Error(w, "Parameter missing or not a number", http.StatusBadRequest)
				return
			}
			b, err := strconv.Atoi(bParam)
			if err != nil {
				http.Error(w, "Parameter missing or not a number", http.StatusBadRequest)
				return
			}

			response := OpResponse{a, b, a + b}
			_ = json.NewEncoder(w).Encode(response)
		}
	}
}

func returnMultiply(w http.ResponseWriter, r *http.Request) {
	if isAuth(w, r) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			queryParams := r.URL.Query()
			aParam := queryParams.Get("a")
			bParam := queryParams.Get("b")

			a, err := strconv.Atoi(aParam)
			if err != nil {
				http.Error(w, "Parameter missing or not a number", http.StatusBadRequest)
				return
			}
			b, err := strconv.Atoi(bParam)
			if err != nil {
				http.Error(w, "Parameter missing or not a number", http.StatusBadRequest)
				return
			}

			response := OpResponse{a, b, a * b}
			_ = json.NewEncoder(w).Encode(response)
		}
	}
}

func fetch(w http.ResponseWriter, r *http.Request) {
	if isAuth(w, r) {
		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")

			type MyJSON struct {
				URL string
			}
			var j MyJSON
			_ = json.NewDecoder(r.Body).Decode(&j)

			rs, _ := http.Get(j.URL)
			_ = json.NewEncoder(w).Encode(rs.Header)
			_ = rs.Body.Close()
		}
	}
}

func writeDataAdela(w http.ResponseWriter, r *http.Request) {
	if isAuth(w, r) {
		if r.Method == http.MethodPost {
			_ = json.NewDecoder(r.Body).Decode(&adela)
			createFileAdela("student1.txt")
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/plain")

			data, err := os.ReadFile("student1.txt")
			if err != nil {
				http.Error(w, "File not yet created", http.StatusNotFound)
			}
			_, _ = w.Write(data)
		}
	}
}

func writeDataIvo(w http.ResponseWriter, r *http.Request) {
	if isAuth(w, r) {
		if r.Method == http.MethodPost {
			_ = json.NewDecoder(r.Body).Decode(&ivo)
			createFileIvo("student2.txt")
		}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/plain")

			data, err := os.ReadFile("student2.txt")
			if err != nil {
				http.Error(w, "File not yet created", http.StatusNotFound)
			}
			_, _ = w.Write(data)
		}
	}
}

func main() {
	yamlFile, _ := os.ReadFile("config.yaml")
	_ = yaml.Unmarshal(yamlFile, &c)
	loadUserPass()

	mux := http.NewServeMux()
	mux.HandleFunc("/jmbag", returnJMBAG)
	mux.HandleFunc("/sum", returnSum)
	mux.HandleFunc("/multiply", returnMultiply)
	mux.HandleFunc("/fetch", fetch)
	mux.HandleFunc("/0246096698", writeDataAdela)
	mux.HandleFunc("/0036522500", writeDataIvo)

	go saxpy.Init()
	mux.HandleFunc("/saxpy", saxpy.Handler)

	srv := &http.Server{
		Addr:              c.HTTP.Address + ":" + c.HTTP.Port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5}
	log.Fatal(srv.ListenAndServe())
}

var usersPasswords map[string][]byte

func loadUserPass() {
	usersPasswords = make(map[string][]byte)
	for _, u := range c.Users {
		usersPasswords[u.Name] = []byte(u.Password)
	}
}
func verifyUserPass(name, password string) bool {
	wantPass, hasUser := usersPasswords[name]
	if !hasUser {
		return false
	}
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}

var c = &Config{}

type Config struct {
	JMBAG string `yaml:"jmbag"`
	HTTP  struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"http"`
	Users []struct {
		Name     string `yaml:"name"`
		JMBAG    string `yaml:"jmbag"`
		Password string `yaml:"password"`
	} `yaml:"users"`
}
type OpResponse struct {
	A      int `json:"a"`
	B      int `json:"b"`
	Result int `json:"result"`
}
type Student struct {
	Ime     string
	Prezime string
	JMBAG   string
}

var adela Student
var ivo Student

func createFileAdela(filename string) {
	file, _ := os.Create(filename)
	defer file.Close()

	_, _ = file.WriteString(adela.Ime + " " + adela.Prezime + " " + adela.JMBAG)
}
func createFileIvo(filename string) {
	file, _ := os.Create(filename)
	defer file.Close()

	_, _ = file.WriteString(ivo.Ime + " " + ivo.Prezime + " " + ivo.JMBAG)
}