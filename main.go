package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	config "github.com/KevinAntonius/RegistrationCodingTest/config"
	model "github.com/KevinAntonius/RegistrationCodingTest/model"
	null "gopkg.in/guregu/null.v3"
)

func main(){
	model.Initialize()
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginPage)
	router.HandleFunc("/api/login", CheckExist)
	router.HandleFunc("/register", RegisterPage)
	router.HandleFunc("/api/register", SubmitRegister)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	fmt.Println("Listening to port 8080...")
	http.ListenAndServe(":8080", router)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        tmpl := template.Must(template.New("register").ParseFiles("template/register.html"))
        err := tmpl.Execute(w, nil)

        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
        return
    }
    http.Error(w, "Request Method Error", http.StatusBadRequest)
}

func SubmitRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		string_param := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)
		db, err := gorm.Open(config.DB_ADAPTER, string_param)
		if err != nil{
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var result map[string]interface{}
		data := &model.User{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&result)
		if err != nil{
			http.Error(w, "Request Syntax Error", http.StatusBadRequest)
			return
		}
		for key, value := range result {
			if value != nil || value.(string) != "" {
				if key == "mobile"{
					if PhoneValidation(value.(string)){
						data.Mobile = null.StringFrom(value.(string))
					} else{
						http.Error(w, "not-valid mobile", http.StatusConflict)
						return
					}
				} else if key == "fname"{
					data.FirstName = null.StringFrom(value.(string))
				} else if key == "lname"{
					data.LastName = null.StringFrom(value.(string))
				} else if key == "birthdate"{
					t, _ := time.Parse("2006-01-02", value.(string))
					data.BirthDate = &t
				} else if key == "gender"{
					data.Gender = null.StringFrom(value.(string))
				} else if key == "email"{
					if MailValidation(value.(string)){
						data.Email = null.StringFrom(value.(string))
					} else{
						http.Error(w, "not-valid email", http.StatusConflict)
						return
					}
				}
			}
		}
		err = db.Create(&data).Error
		if err != nil {
			error_str := err.Error()
			error_msg := ""
			if strings.Contains(error_str,"null"){
				error_msg = "null "

			} else if strings.Contains(error_str,"duplicate"){
				error_msg = "duplicate "
			} else{
				error_msg = "Internal Server Error"
				http.Error(w, error_msg, http.StatusInternalServerError)
				return
			}
			if strings.Contains(error_str,"mobile"){
				error_msg += "mobile"
			} else if strings.Contains(error_str,"first_name"){
				error_msg += "fname"
			} else if strings.Contains(error_str,"last_name"){
				error_msg += "lname"
			} else if strings.Contains(error_str,"email"){
				error_msg += "email"
			}
			http.Error(w, error_msg, http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.Error(w, "Request Method Error", http.StatusBadRequest)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		tmpl := template.Must(template.New("login").ParseFiles("template/login.html"))
		
		data := map[string]interface{}{
			"name":  r.Form.Get("fname"),
		}
		err := tmpl.Execute(w, data)

        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
        return
    }
	http.Error(w, "Request Method Error", http.StatusBadRequest)
}

func CheckExist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		string_param := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SSLMODE)
		db, err := gorm.Open(config.DB_ADAPTER, string_param)
		if err != nil{
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		r.ParseForm()
		user := &model.User{}
		email := r.Form.Get("email")
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			http.Error(w, "User Not Exist", http.StatusNotFound)
			return
		}
		response := make(map[string]string)
		response["fname"] = user.FirstName.String
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
        return
    }
	http.Error(w, "Request Method Error", http.StatusBadRequest)
}

func PhoneValidation(phone string) bool{
	re := regexp.MustCompile(`^(?:\+62|\(0\d{2,3}\)|0)\s?(?:361|8[17]\s?\d?)?(?:[ -]?\d{3,4}){2,3}$`)
	return re.MatchString(phone)
}

func MailValidation(email string) bool{
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}