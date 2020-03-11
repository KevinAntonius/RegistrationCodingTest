package main

import (
	"bytes"
	"strconv"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func TestPhoneValidation(t *testing.T) {
	var mobiles = [14]string{"+62 361 222777", "+62 813-444-5555", "+62 812-3333-3333","+62 811 391 2103",
		"+62 361-2277777","(0361) 227337","+62 8113912103","08134455555","0361-2277777","+62 812 3333 3333",
		"+62 877 80803550","081339222111","081 339 222 111","+62 811338429196"}
	for _, mobile := range mobiles{
		result := PhoneValidation(mobile)
		t.Log("Phone : "+mobile+" -> "+strconv.FormatBool(result))
		if !result {
			t.Error("Wrong")
		}
	}
}

func TestMailValidation(t *testing.T) {
	var emails = [12]string{"simple@example.com","very.common@example.com","disposable.style.email.with+symbol@example.com",
		"other.email-with-hyphen@example.com","fully-qualified-domain@example.com","user.name+tag+sorting@example.com", "x@example.com",
		"example-indeed@strange-example.com","admin@mailserver1","example@s.example", "mailhost!username@example.org","user%example.com@example.org"}
	for _, email := range emails{
		result := MailValidation(email)
		t.Log("Email : "+email+" -> "+strconv.FormatBool(result))
		if !result {
			t.Error("Wrong")
		}
	}
}

func TestRegisterPageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(RegisterPage).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func TestSubmitRegisterHandler(t *testing.T) {
	param := make(map[string]string)
	param["mobile"] = "080989999"
	param["fname"] = "Satu"
	param["lname"] = "Dua"
	param["birthdate"] = "2020-03-11"
	param["email"] = "example@example.com"

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(param)
	r, err := http.NewRequest("POST", "/api/register", buf)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	http.HandlerFunc(SubmitRegister).ServeHTTP(w, r)

	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func TestLoginPageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	q := req.URL.Query()
    q.Add("name", "Gordon")
    req.URL.RawQuery = q.Encode()


	rr := httptest.NewRecorder()

	http.HandlerFunc(LoginPage).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func TestCheckExistHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/login", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
	q := r.URL.Query()
    q.Add("email", "example@example.com")
    r.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	http.HandlerFunc(CheckExist).ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
	
	var result map[string]interface{}
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&result)
	if err != nil{
		t.Errorf("Json format is not correct: %v", err)
	}

	if result["fname"].(string) != "Satu"{
		t.Errorf("Wrong value: %s", result["fname"].(string))
	}
}