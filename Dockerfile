# build binary stage
FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY . /app/
RUN go build -o main main.go

# move binary file to lightweight image
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 8080
CMD ["/app/main"]