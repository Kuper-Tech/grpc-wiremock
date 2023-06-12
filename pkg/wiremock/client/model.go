package client

type Request struct {
	URLPath string `json:"urlPath"`
	Method  string `json:"method"`
}

type Accept struct {
	Contains string `json:"contains"`
}

type Response struct {
	Status  int     `json:"status"`
	Body    string  `json:"body"`
	Headers Headers `json:"headers"`
}

type Headers struct {
	ContentType string `json:"Content-Type"`
}

type Metadata struct {
	Description string `json:"description"`
}

type Mock struct {
	Name string `json:"name"`

	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Metadata Metadata `json:"metadata"`
}

func DefaultMock() Mock {
	return Mock{
		Request:  Request{},
		Response: Response{Headers: Headers{ContentType: "application/json"}},
	}
}

func (m *Mock) WithName(name string) *Mock {
	m.Name = name
	return m
}

func (m *Mock) WithDescription(description string) *Mock {
	m.Metadata.Description = description
	return m
}

func (m *Mock) WithRequestUrlPath(urlPath string) *Mock {
	m.Request.URLPath = urlPath
	return m
}

func (m *Mock) WithRequestMethod(method string) *Mock {
	m.Request.Method = method
	return m
}

func (m *Mock) WithResponseStatusCode(statusCode int) *Mock {
	m.Response.Status = statusCode
	return m
}

func (m *Mock) WithResponseBody(body string) *Mock {
	m.Response.Body = body
	return m
}
