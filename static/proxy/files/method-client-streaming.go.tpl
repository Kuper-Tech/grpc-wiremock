package {{ .PackageHeader }}

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"grpc-proxy/pkg/wiremock"

	{{ range .GoPackages }}
	"{{ . }}"
	{{- end }}
)

func (p *Service) {{ .Method }}(stream {{ .MethodPackage }}.{{ .Service }}_{{ .Method }}Server) error {
	const url = "{{ .URL }}"

	unmarshalAndSend := func(responseBody []byte) error {
		var protoResponse {{ .MethodOutPackage }}.{{ .MethodOutName }}
		if processErr := protojson.Unmarshal(responseBody, &protoResponse); processErr != nil {
			return status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", processErr))
		}
		if processErr := stream.SendAndClose(&protoResponse); processErr != nil {
			return processErr
		}
		return nil
	}

	defaultRequest, err := wiremock.DefaultRequest(stream.Context(), url, bytes.NewReader([]byte{}))
	if err != nil {
		return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
	}

	httpResponseBody, streamSize, err := wiremock.DoRequestWithStreamSize(defaultRequest)
	if err != nil {
		return err
	}

	streamCursor := 1

	for {
		req, errReceive := stream.Recv()
		if errReceive != nil && errReceive == io.EOF {
			return unmarshalAndSend(httpResponseBody)
		}
		if errReceive != nil {
			return errReceive
		}
		if streamCursor >= streamSize {
			return unmarshalAndSend(httpResponseBody)
		}
		if req == nil {
			continue
		}
		streamCursor++
	}
}
