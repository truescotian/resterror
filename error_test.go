package error_test

func ExampleErrorMessage() {
	if msg := ErrorMessage(err); msg != "" {
		fmt.Printf("ERROR: %s\n", msg)
	}
}
