package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	err := a.Initialize(DBUser, DBPass, "test")
	if err != nil {
		log.Fatal("error occurred while initializing test database")
	}
	createTable()
	m.Run()

}

func createTable() {
	createTable := `CREATE TABLE IF NOT EXISTS products(
    		id INT AUTO_INCREMENT PRIMARY KEY,
    		name VARCHAR(50) NOT NULL,
    		price FLOAT NOT NULL,
    		quantity INT NOT NULL
    		);`
	_, err := a.DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
}

func addProduct(name string, price float64, quantity int) {
	_, err := a.DB.Exec("INSERT INTO products(name, price, quantity) VALUES(?, ?, ?)", name, price, quantity)
	if err != nil {
		log.Fatal(err)
	}

}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("Keyboard", 100, 5000)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func sendRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, req)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	payload := []byte(`{"name":"Chair","price":100,"quantity":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "Chair" {
		t.Errorf("Expected product name to be 'Chair'. Got '%v'", m["name"])
	}
	if m["quantity"] != 100.0 {
		t.Errorf("Expected product quantity. Got '%v'", m["quantity"])
	}
	//err := json.Unmarshal(response.Body.Bytes(), &m)
	//if err != nil {
	//	t.Errorf("Error occurred while unmarshalling response body")
	//}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("Keyboard", 100, 5000)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	//req, _ = http.NewRequest("GET", "/product/1", nil)
	//response = sendRequest(req)
	//checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("Keyboard", 100, 5000)
	payload := []byte(`{"name":"Keyboard","price":200,"quantity":5000}`)
	req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["price"] != 200.0 {
		t.Errorf("Expected product price to be 200. Got '%v'", m["price"])
	}
	if m["quantity"] != 5000.0 {
		t.Errorf("Expected product quantity to be 5000. Got '%v'", m["quantity"])
	}
}
