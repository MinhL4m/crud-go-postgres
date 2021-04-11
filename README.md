# CRUD Go & Postgres

## Overview

Purpose of this application is for studying: golang + postgres, golang testing, postgres+docker for testing.

Step:

1. Create base scaffold for the whole program. Model + App + base main and main_test
2. Create .env
3. Test connect with db. In this case, I created a docker postgres container.
4. Create all test case
5. Run test case to see if all return 404.

## Dependencies

- gorilla/mux: Package `gorilla/mux` implements a request router and dispatcher for matching incoming requests to their respective handler.

- lib\pq: Go postgres driver for Go's database/sql package

## Run test

`go test -v`

## Structure

Main -> App -> Model -> DB

- Instead of DB return Model, methods in model will access db and get db then directly modify the model.

```go
p := product{ID: id}
p.getProduct(a.DB) // p right now will by modify with all the information returned from db 
```

## Learnt

### Postgres + Docker for testing

- Use postgres for testing: [How to run PostgreSQL in Docker on Mac (for local development)](https://www.saltycrane.com/blog/2019/01/how-run-postgresql-docker-mac-local-development/)

- Summary:

  - Running in Docker allows keeping my database environment isolated from the rest of my system and allows running multiple versions and instances
  - For this application, I used the first option

  - `docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres`: This command will find an image of postgres and install it, if local doens't has it.

    - -it: allocate pseudo-TTY
    - -d: Run container in background and print container ID
    - -p: Publish a container's port(s) to the host

  - To connect to postgres inside docker: [Github](https://stackoverflow.com/questions/37694987/connecting-to-postgresql-in-a-docker-container-from-outside)

### underscore in front of import package

- [What does an underscore in front of an import statement mean?](https://stackoverflow.com/questions/21220077/what-does-an-underscore-in-front-of-an-import-statement-mean)

- Summary:
  - It's for importing a package solely for its side-effects.
  - In our case `pd` is used as driver

### Pointer

Little demonstration:

```go
package main

import "fmt"

func main() {
    i := 42
    fmt.Printf("i: %[1]T %[1]d\n", i)
    p := &i
    fmt.Printf("p: %[1]T %[1]p\n", p)
    j := *p
    fmt.Printf("j: %[1]T %[1]d\n", j)
    q := &p
    fmt.Printf("q: %[1]T %[1]p\n", q)
    k := **q
    fmt.Printf("k: %[1]T %[1]d\n", k)
}

// Output
i: int 42
p: *int 0x10410020
j: int 42
q: **int 0x1040c130
k: int 42
```

- Summary:
  - If datatype of variable or paramater contains `*` -> ask for address -> pass in with `&`
    --> Use for change value.
  - Inside function, if the param has `*`, need to use `*` to dereference.

### Create Scaffolding a Minimal Application

- Before we can write tests, we need to create a minimal application that can be used as the basis for the tests.

- In this application, I used `app.go` to create minimal application. By doing this, `main.go`(dev, prod) and `main_test.go`(test) can both use the `App` to run the application.

- For model, create all functions 1 model need. Example:

```go
type human struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
}

func (h *human) getHuman(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (h *human) updateHuman(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (h *human) deleteHuman(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (h *human) createHuman(db *sql.DB) error {
  return errors.New("Not implemented")
}
```

--> By doing this, We have everything layout. This method also helpful in testing

### Should I define methods on values or pointers?

```go
func (s *MyStruct) pointerMethod() { } // method on pointer
func (s MyStruct)  valueMethod()   { } // method on value
```

- **First, and most important, does the method need to modify the receiver?** If it does, the receiver must be a pointer. (Slices and maps act as references, so their story is a little more subtle, but for instance to change the length of a slice in a method the receiver must still be a pointer.)

- In the examples above, if pointerMethod modifies the fields of s, the caller will see those changes, but valueMethod is called with a copy of the caller's argument (that's the definition of passing a value), so changes it makes will be invisible to the caller.

### Env variables

- [Use Environment Variable in your next Golang Project](https://towardsdatascience.com/use-environment-variable-in-your-next-golang-project-39e17c3aaa66)

- Use `os` + `godotenv` package

```go
// Load the .env file in the current directory
godotenv.Load()

// or

godotenv.Load(".env")

// Work with os

err := godotenv.Load(".env")

os.Getenv(key)
```

### Testing with TestMain

[Why Use TestMain For Testing?](https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc)

- Summary: By using `TestMain`, setup and tear down included in the test.

### Testing with fake request (make request and see how handler handle it):

- [Testing Http handler](https://blog.questionable.services/article/testing-http-handlers-go/)

- using `"net/http/httptest"` package

- Summary:

  1. Create new request with `req, err := http.NewRequest(<request>)`
  2. Create new recorder (`*httptest.ResponseRecorder` type -> record the response so we can check later) with `rr := httptest.NewRecorder()`
  3. Create handler with `handler :=http.HandlerFunc(<function>)`
  4. Run request with ` handler.ServeHTTP(rr, req)`
  5. Check if the response match what we expected by check the `recorder`:

  ```go
  if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body is what we expect.
    expected := `{"alive": true}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
  ```

- In this application, run request step and check http status code will be resuable functions.

### How to use JSON with Go -> Testing

- [How to use JSON with Go [best practices]](https://yourbasic.org/golang/json-example/)
- However, in this application, I also used `[]byte()` to create function (`TestCreateProduct`)

- Summary:
  - Encode (marshal) struct to JSON -> use `json.Marshal(<struct>)` struct to json
  - Decode (unmarshal) JSON to struct or map -> use `json.Unmarshal(<data []byte>, <interface{}>)`. In this application, I converted json to map and used map to check if error exist.

### Unmarshal vs newDecoder.Decode

- It really depends on what your input is. If you look at the implementation of the `Decode` method of `json.Decoder`, **it buffers the entire JSON value in memory before unmarshalling it into a Go value**. So in most cases it won't be any more memory efficient (although this could easily change in a future version of the language).

- So a better rule of thumb is this:

  - Use `json.Decoder` if your data is coming from an io.Reader stream, or you need to decode multiple values from a stream of data. res.Body is stream (need to close()) therefore need to use Decoder stead of Unmarshal.
  - Use`json.Unmarshal` if you already have the JSON data in memory.

For the case of reading from an HTTP request, I'd pick `json.Decoder` since you're obviously reading from a stream. However, in the testing, we already has JSON data in the memeory.

### Query with postgres

- To **get** or **create** -> `db.QueryRow()`: only 1 row or `db.Query()`: multiple rows
- To **delete** or **update** -> `db.Exec()`
- **Scan()**, to get more info go to `model.go`:
  - When get results back from **get** requests, use `Scan()` to modify current variable with the result
  - `db.Query()` will return `sql.rows`:
    - Always `defer rows.Close()`
    - To loop over rows: `for rows.Next()`
    - Inside the for loop, rows variable now will point over each row. To get the value out of row. Use `Scan`

### Request Payload vs Form Data

- [Stackoverflow](https://stackoverflow.com/questions/23118249/whats-the-difference-between-request-payload-vs-form-data-as-seen-in-chrome)

- Summary:
  - Request payload or payload body is sent using `PUT` or `POST`
  ```json
  POST /some-path HTTP/1.1
  Content-Type: application/json
  { "foo" : "bar", "name" : "John" }
  ```
  - Form data is sent using submit form with `POST`
  ```json
  POST /some-path HTTP/1.1
  Content-Type: application/x-www-form-urlencoded
  foo=bar&name=John
  ```
