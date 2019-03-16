package functional_tests

type TestStruct struct {
	Name string
	RequestBody string
	ExpectedStatusCode int
	ExistsFields []string
}