package functional_tests

type TestStruct struct {
	Name               string
	RequestBody        string
	Headers            map[string]string
	ExpectedStatusCode int
	ExistsFields       []string
}
