# syntax=docker/dockerfile:1
FROM golang:1.23 AS build
WORKDIR /cmd

# Copy the source code
COPY . .

RUN go get -d -v ./...


# Build the Go app
RUN go build -o cmd/api .

RUN go build -o /bin/devices_api ./main.go

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]