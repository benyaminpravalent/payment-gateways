FROM golang:1.22-alpine

# Set up environment and install necessary packages
RUN apk add --no-cache git netcat-openbsd gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Default command for the container
CMD ["go", "run", "app/main.go"]
