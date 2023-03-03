## Proxy generator

Inside the container you can find the **grpc2http** CLI utility,
it allows to generate a proxy server (ready to run Golang project)
based on the provided **proto** contracts.

To generate code from the proto contracts used:
- proto [compiler](https://github.com/protocolbuffers/protobuf);
- go-grpc and go plugins.

### Install
```bash
make build -C "cmd/grpc2http"
```

### Interface

Help:

```bash
grpc2http -h
```

How to generate a proxy:

```bash
grpc2http --input "/tmp/my-awesome-contracts-dir" --output "/tmp/ready-to-run-proxy"
```

How to run a proxy:

```bash
make run -C "/tmp/ready-to-run-proxy"
```

## How it works?

gGRP code and stubs are generated for each of your contracts.

Within the stubs, the gRPC request is converted into HTTP,
then the HTTP request is sent to the Wiremock API,
and its response is converted back to gRPC.

### Wiremock URL

URL construction rule. Suppose that your gRPC server is called 
```AwesomeService```.  And the method is ```CallSomething```. 
Base URL is ```http://localhost:8080```.
In this case, the URL would be as follows:

```
http://localhost:8080/AwesomeService/CallSomething
```

This URL must be present in mocks.

### Examples

Examples of stubs for each type of gRPC interactions:
- Unary call
    ```Go
    func (p *Service) Unary(ctx context.Context, in *example.Request) (*example.Response, error) {
        const url = "http://localhost:8080/Example/Unary"
    
        requestBody, err := protojson.Marshal(in)
        if err != nil {
            return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("create http request body: %v", err))
        }
        
        request, err := wiremock.DefaultRequest(ctx, url, bytes.NewReader(requestBody))
        if err != nil {
            return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
        }
        
        httpResponseBody, err := wiremock.DoRequestDefault(request)
        if err != nil {
            return nil, err
        }
        
        var protoResponse example.Response
        if err = protojson.Unmarshal(httpResponseBody, &protoResponse); err != nil {
            return nil, status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", err))
        }
        
        return &protoResponse, nil
    }
    ```

- Client-side streaming
    ```Go
    func (p *Service) ClientSideStream(stream example.Example_ClientSideStreamServer) error {
        const url = "http://localhost:8080/Example/ClientSideStream"
    
        unmarshalAndSend := func(responseBody []byte) error {
            var protoResponse example.Response
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
    ```

- Server-side streaming
  ```Go
    func (p *Service) ServerSideStream(in *example.Request, stream example.Example_ServerSideStreamServer) error {
        const url = "http://localhost:8080/Example/ServerSideStream"
        const streamCursor = 1
    
        ctx := stream.Context()
    
        unmarshalAndSend := func(responseBody []byte) error {
            var protoResponse example.Response
            if processErr := protojson.Unmarshal(responseBody, &protoResponse); processErr != nil {
                return status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", processErr))
            }
            if processErr := stream.Send(&protoResponse); processErr != nil {
                return processErr
            }
            return nil
        }
    
        processStream := func(cursor int) error {
            httpRequest, processErr := wiremock.RequestWithCursor(ctx, url, cursor, bytes.NewReader([]byte{}))
            if processErr != nil {
                return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", processErr))
            }
            httpResponseBody, processErr := wiremock.DoRequestDefault(httpRequest)
            if processErr != nil {
                return processErr
            }
            return unmarshalAndSend(httpResponseBody)
        }
    
        defaultRequest, err := wiremock.RequestWithCursor(ctx, url, streamCursor, bytes.NewReader([]byte{}))
        if err != nil {
            return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
        }
    
        httpResponseBody, streamSize, err := wiremock.DoRequestWithStreamSize(defaultRequest)
        if err != nil {
            return err
        }
    
        if err = unmarshalAndSend(httpResponseBody); err != nil {
            return err
        }
    
        for cursor := streamCursor + 1; cursor <= streamSize; cursor++ {
            if err = processStream(cursor); err != nil {
                return err
            }
        }
    
        return nil
    }
  ```

- Bidirectional streaming
    ```Go
    func (p *Service) BidirectionalStream(stream example.Example_BidirectionalStreamServer) error {
        const url = "http://localhost:8080/Example/BidirectionalStream"
    
        ctx := stream.Context()
    
        unmarshalAndSend := func(responseBody []byte) error {
            var protoResponse example.Response
            if processErr := protojson.Unmarshal(responseBody, &protoResponse); processErr != nil {
                return status.Error(http.StatusBadGateway, fmt.Sprintf("marshal json object to proto: %v", processErr))
            }
            if processErr := stream.Send(&protoResponse); processErr != nil {
                return processErr
            }
            return nil
        }
    
        processStream := func(cursor int) error {
            httpRequest, processErr := wiremock.RequestWithCursor(ctx, url, cursor, bytes.NewReader([]byte{}))
            if processErr != nil {
                return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", processErr))
            }
            httpResponseBody, processErr := wiremock.DoRequestDefault(httpRequest)
            if processErr != nil {
                return processErr
            }
            return unmarshalAndSend(httpResponseBody)
        }
    
        streamCursor := 1
    
        request, err := wiremock.RequestWithCursor(ctx, url, streamCursor, bytes.NewReader([]byte{}))
        if err != nil {
            return status.Error(http.StatusBadGateway, fmt.Sprintf("create http request: %v", err))
        }
    
        _, streamSize, err := wiremock.DoRequestWithStreamSize(request)
        if err != nil {
            return err
        }
    
        for {
            req, errReceive := stream.Recv()
            if errReceive != nil && errReceive == io.EOF {
                return nil
            }
            if errReceive != nil {
                return errReceive
            }
            if req == nil {
                continue
            }
            if err = processStream(streamCursor); err != nil {
                return err
            }
            if streamCursor >= streamSize {
                return nil
            }
            streamCursor++
        }
    }
    ```
