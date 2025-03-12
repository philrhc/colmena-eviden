FROM registry.atosresearch.eu:18488/golang:latest-alpine as builder

# Use to skip tests and binary builds to be able to use tools/dev.docker-compose.yml 
ARG DEV

WORKDIR /go/src/colmena/components/sla-manager

# Only copy require code and delete as much as possible (second stage)
#COPY assets/go-common/components/utils ../assets/go-common/components/utils
COPY tools/dockerfile_tools/ .
COPY components/sla-manager-svc .

# Neeed for tests
RUN apk --no-cache add build-base

# Make sure that the code is up to standards and docs are updated
# Run tests and keep report
# RUN if [ -z "$DEV" ]; then \
#     ./go_tools.sh \
#     && ./go_test.sh \
#     && git reset --hard HEAD ; \
#     fi

# Finally build binaries
RUN if [ -z "$DEV" ]; then \
    GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -v -a -installsuffix cgo -o SLAMgr . ;\
    fi

# Reduce final image tag as much as possible
FROM registry.atosresearch.eu:18520/library/alpine:3.7


# Set the current working directory inside the container 
WORKDIR /root/

COPY --from=builder /go/src/colmena/components/sla-manager/SLAMgr ./SLAMgr
#COPY --from=builder /go/src/colmena/components/sla-manager/testReport.txt testReport.txt

# Expose port 8080
EXPOSE 8080

#Command to run the executable
CMD ["./SLAMgr"]
