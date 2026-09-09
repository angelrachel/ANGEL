FROM golang:1.27.1 AS builder
WORKDIR /src
COPY go.mod ./
COPY src/ ./src/
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/angel-server ./src/c2

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/angel-server /angel-server
ENV ANGEL_HOST=0.0.0.0
ENV ANGEL_PORT=8001
EXPOSE 8001
USER nonroot:nonroot
ENTRYPOINT ["/angel-server"]
