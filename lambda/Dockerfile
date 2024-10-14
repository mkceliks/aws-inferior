# Use the official Golang image to build the app
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (for Linux since Lambda runs on Linux)
RUN GOOS=linux GOARCH=amd64 go build -o main

# Use AWS Lambda base image for Go
FROM public.ecr.aws/lambda/go:latest

# Copy the binary to the Lambda container
COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

# Command to run the Lambda function
CMD ["main"]
