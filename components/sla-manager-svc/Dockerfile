# THIS IS A DOCKERFILE WITH TWO STAGES
# FIRST IMAGE IS CUSTOMIZED TO BUILD THE COMPONENT
# THE SECOND IS THE RUNTIME IMAGE WHICH COPIES THE BINARI FILE
# AND THE NECESSARY TEMPLATES

##########################################################
# Start from golang base image
FROM golang:alpine AS builder

# Set the current working directory inside the container 
WORKDIR /go/src/sla

# Copy folders, main, go mod and sum files / PATH RELATIVE TO INSTALL FOLDER !!!
COPY ./ ./

RUN ls -a

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o SLAMgr .

##########################################################
# Start a new stage from scratch
FROM scratch AS runtime

WORKDIR /opt/sla

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/sla/SLAMgr .

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["/opt/sla/SLAMgr"]