package main

// Wrapper for handler functions.
type rootHandler func(http.ResponseWriter, *http.Request) error

func testHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return NewRESTError(nil, 405, "Method not allowed.")
	}

	body, err := ioutil.ReadAll(r.Body) // read request body.
	if err != nil {
		return fmt.Errorf("Request body read error : %v", err)
	}

	// Parse body as json.
	if err := json.Unmarshal(body, &schema); err != nil {
		return NewHttpError(err, 400, "Bad request : invalid JSON.")
	}

	ok, err := loginUser("username", "password")
	if err != nil {
		return fmt.Errorf("loginUser DB error : %v", err)
	}

	if !ok { // Authentication failed.
		return NewHTTPError(nil, 401, "Wrong password or username.")
	}

	w.WriteHeader(200) // Successfully authenticated.
	return nil
}

// Implement the http.Handler interface.
func (fn rootHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r) // Call handler function.
	if err == nil {
		return
	}

	log.Printf("An error occured. %v", err) // log error.

	clientError, ok := err.(ClientError) // Check if it's a ClientError.
	if !ok {
		// If not ClientError, assume it's ServerError
		w.WriteHeader(500)
		return
	}

	body, err := clientError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error accured: %v", err)
		w.WriteHeader(500)
		return
	}

	status, headers := clientError.ResponseHeaders() // Get http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}

func main() {
	// http.Handle accepts any type that implements http.Handler interface,
	// so as long as you pass a type that has ServeHttp method, the http.Handle
	// method will be happy.
	http.Handle("/", rootHandler(testHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
