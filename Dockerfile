# 1. Use a lightweight Go base image (Alpine Linux)
# This acts as the OS for our container.
FROM golang:1.25-alpine

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy Go module files first (for caching dependencies)
COPY go.mod go.sum ./

# 4. Download dependencies
RUN go mod download

# 5. Copy the source code into the container
COPY . .

# 6. Build the Go application
# This creates a binary executable named "server"
RUN go build -o server .

# 7. Expose the port the app runs on
EXPOSE 8080

# 8. Define the command to start the app
CMD ["./server"]