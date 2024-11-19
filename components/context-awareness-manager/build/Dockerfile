# THIS IS A DOCKERFILE WITH TWO STAGES
# FIRST IMAGE IS CUSTOMIZED TO BUILD THE COMPONENT
# THE SECOND IS THE RUNTIME IMAGE WHICH COPIES THE BINARI FILE
# AND THE NECESSARY TEMPLATES

# Start from golang base image
FROM golang:alpine AS builder

# Install SQLite3 dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the current working directory inside the container 
WORKDIR /go/src

# Copy go mod and sum files 
COPY components/context-awareness-manager ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/context-awareness-manager

# Start a new stage from scratch
FROM alpine:3.17

# Install SQLite3 runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

RUN addgroup -S etsn && adduser -S etsn -G etsn -u 1001
USER etsn

# Set the current working directory inside the container 
WORKDIR /home/etsn/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/main .

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["./main"]