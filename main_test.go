package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/MinhL4m/crud-go-postgres"
	"github.com/joho/godotenv"
)

// This is in another package
// need to import "github.com/MinhL4m/crud-go-postgres" to work
var a main.App

// TestMain include setup and tear down steps
func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		os.Exit(1)
	}

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

/**----------------Util Functions-----------*/
// Check if table exist if not create one.
// This will be used in set up step
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// Delete all rows added during test and rest id to 1
// This will be used in tear down step
func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

// Execute request and return ResponseRecorder
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

// Check if http status code match the expected
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Add mock products
func addProducts(count int) {
	if count < 1{
		count = 1
	}
	for i := 0; i < count; i++ {
        a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
    }
}

/**----------------Test Cases-----------*/

/**-----Test Case1: Empty Table--------*/
func TestEmptyTable(t *testing.T) {
	// Delete all row before test
	clearTable()
	req, _ := http.NewRequest("GET", "/products", nil)
	res := executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)

	// Check res body
	if body := res.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

/**-----Test Case2: Fetch Non-existent Product--------*/
func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	res := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, res.Code)

	// Create map to store body
	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

/**-----Test Case3: Create Product--------*/
func TestCreateProduct(t *testing.T) {
	clearTable()

	jsonStr := []byte(`{"name":"test product", "price": 11.22}`)

	//NewBuffer is intended to prepare a Buffer to read existing data.
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, res.Code)

	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

/**-----Test Case4: Fetch Product--------*/
func TestGetProduct(t *testing.T){
	clearTable()
	addProducts(1)

	req,_:= http.NewRequest("GET","/product/1",nil)

	res:=executeRequest(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

/**-----Test Case5: Update Product--------*/
func TestUpdateProduct(t *testing.T) {

    clearTable()
    addProducts(1)

	// Get current product
    req, _ := http.NewRequest("GET", "/product/1", nil)
    res := executeRequest(req)
    var originalProduct map[string]interface{}
    json.Unmarshal(res.Body.Bytes(), &originalProduct)

	// Put update product
    var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
    req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    res = executeRequest(req)

    checkResponseCode(t, http.StatusOK, res.Code)

    var m map[string]interface{}
    json.Unmarshal(res.Body.Bytes(), &m)

	// check if id match
    if m["id"] != originalProduct["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
    }

	// check if product get updated
    if m["name"] == originalProduct["name"] {
        t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
    }

    if m["price"] == originalProduct["price"] {
        t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
    }
}

/**-----Test Case5: Delete Product--------*/
func TestDeleteProduct(t *testing.T) {
    clearTable()
    addProducts(1)

    req, _ := http.NewRequest("GET", "/product/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/product/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/product/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}