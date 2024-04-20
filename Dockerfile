FROM golang:1.20.2 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o app main.go

FROM alpine:3.16.2
WORKDIR /app
COPY --from=builder /app/app app
RUN apk add uuidgen tzdata
EXPOSE 8080
ENTRYPOINT ["/app/app"]
