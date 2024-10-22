# use official Golang image
FROM golang:1.23.2 AS builder

# Set name image
LABEL name="storage-api"

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code of the project and configuration files into the container
COPY . .

# Set the working directory to /app
WORKDIR /app/src

# Compile the project. This assumes your main.go file is in the root of the /app directory.
# The '-ldflags="-w -s"' flag is used to reduce the size of the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o ../storage-api .

# Execution stage using scratch
FROM scratch AS runner

# Set the working directory to /app
WORKDIR /app

# Import the CA certificates from the build stage to allow HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the .env file and the compiled binary from the build stage to the execution stage
COPY --from=builder /app/storage-api .

# Command to run the binary
CMD ["./storage-api"]