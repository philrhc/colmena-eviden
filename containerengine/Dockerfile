# Start from golang base image
FROM golang:alpine AS builder

# Set the current working directory inside the container 
WORKDIR /go/src

# Copy go mod and sum files 
COPY ./containerengine ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:3.17

# Add docker group and adduser jrubioc to docker group
# RUN addgroup -S docker && adduser -S jrubioc -G docker -u 1001

# Set the current working directory inside the container 
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/main .

# Set permissions for etsn
# RUN chown -R jrubioc:docker /home/jrubioc
# USER jrubioc

# Expose port 8000
EXPOSE 9000

#Command to run the executable
CMD ["./main"]
