# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Chidiebere Okpala <okpalacollins4@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable. We define the Docker image's behavior at runtime
CMD ["./main"]

#This is a very simple dockerfile for our app. If you want to see a more complicated dockerfile for your go app where
#were you run bash files look at https://github.com/techschool/simplebank/blob/master/Dockerfile