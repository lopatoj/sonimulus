FROM ubuntu:latest
WORKDIR /app
COPY .. .
RUN apt-get update && apt-get upgrade
RUN go build -o server cmd/server/main.go
EXPOSE 3000
