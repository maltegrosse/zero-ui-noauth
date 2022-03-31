# build stage
FROM golang:1.16 AS builder
RUN mkdir -p /go/src/p
WORKDIR /go/src/p
COPY . ./
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app .


# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app ./
RUN chmod +x ./app
ENTRYPOINT ["./app"]
EXPOSE 9999