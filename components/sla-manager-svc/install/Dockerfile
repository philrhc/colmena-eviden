ARG IMAGE_BUILDER=golang:alpine
ARG IMAGE_NAME=alpine:3.17

FROM ${IMAGE_BUILDER} as builder

WORKDIR /go/src/sla-manager

# Only copy require code and delete as much as possible (second stage)
COPY components/sla-manager-svc .

# Neeed for tests
RUN apk --no-cache add build-base

# Finally build binaries
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -v -a -installsuffix cgo -o SLAMgr . ;

# Reduce final image tag as much as possible
FROM ${IMAGE_NAME}

# Set the current working directory inside the container 
WORKDIR /root/

COPY --from=builder /go/src/sla-manager/SLAMgr ./SLAMgr

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["./SLAMgr"]
