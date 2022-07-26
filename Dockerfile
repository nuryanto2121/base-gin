

# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.17 as builder

# Add Maintainer Info
LABEL maintainer="Nuryanto <nuryantofattih@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app

# Copy go mod , sum files config.json
COPY go.mod go.sum ./

COPY config.ini.example ./config.ini 

COPY pkg/multiLanguage/en.json ./

COPY pkg/multiLanguage/id.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache tzdata

ENV TZ Asia/Jakarta

WORKDIR /usr/src/app


# Build Args
ARG LOG_DIR=/usr/src/app/wwwroot
ARG LANG_DIR=/usr/src/app/pkg/multiLanguage

# Create Log Directorytail
RUN mkdir -p ${LOG_DIR}

RUN mkdir -p ${LANG_DIR}

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/config.ini .
COPY --from=builder /usr/src/app/pkg/multiLanguage/en.json ./pkg/multiLanguage/
COPY --from=builder /usr/src/app/pkg/multiLanguage/id.json ./pkg/multiLanguage/

# Expose port 8080 to the outside world
EXPOSE 8000
EXPOSE 9000

# Declare volumes to mount
VOLUME [${LOG_DIR}]

# Command to run the executable
CMD ["./main"] 
