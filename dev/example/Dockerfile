## Build
FROM golang:1.19-buster AS build

WORKDIR /sources

COPY . .

RUN go mod download
RUN go build -o /app cmd/*.go

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app /app

ENTRYPOINT ["/app"]