package mock

// Mock is a basic type, which contains a description of a stub.
// Basis for generating a Wiremock stub.
type Mock struct {
	Name               string
	Description        string
	RequestUrlPath     string
	RequestMethod      string
	ResponseBody       string
	ResponseStatusCode int
}
