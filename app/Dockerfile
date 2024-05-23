# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Install dependencies

# Copy go mod and sum files

COPY . .
# Download all dependencies
RUN go mod tidy

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest


# Generate Swagger docs
RUN /go/bin/swag init

# Build the GoGin app
RUN go build -o main main.go

# Expose port 8080
EXPOSE 8080

ENTRYPOINT ["/bin/bash","-c","./main"]
