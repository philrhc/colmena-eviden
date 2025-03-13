ARG IMAGE_BUILDER=golang:alpine
ARG IMAGE_NAME=alpine:3.17

FROM ${IMAGE_BUILDER} as builder

WORKDIR /go/src/context-awareness-manager

# Only copy require code and delete as much as possible (second stage)
COPY components/context-awareness-manager .

# Need for tests
RUN apk --no-cache add build-base
# Install SQLite3 dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Finally build binaries
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=1 go build -v -a -installsuffix cgo -o main ./cmd ;

# Reduce final image tag as much as possible
FROM ${IMAGE_NAME}

# Install SQLite3 runtime dependencies
RUN apk --no-cache add ca-certificates sqlite docker

# Set the current working directory inside the container 
WORKDIR /root/
USER root

COPY --from=builder /go/src/context-awareness-manager/main ./main

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["./main"]