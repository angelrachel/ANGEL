FROM golang:1.27.1 AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/fixture-http .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/fixture-http /fixture-http
ENV FIXTURE_HTTP_ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/fixture-http"]
