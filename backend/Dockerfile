# Use the official Go image as the base image
FROM golang:1.26

# Set the working directory inside the container
WORKDIR /app

# Copy Go dependency files first
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the complete backend project
COPY . .

# Build the Go application
RUN go build -o food-api .

# Expose the port used by the Go server
EXPOSE 8080

# Start the Go application
CMD ["./food-api"]