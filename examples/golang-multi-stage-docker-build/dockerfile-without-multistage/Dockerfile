###########################################
# BASE IMAGE
###########################################
# Using Ubuntu as the primary operating system for the build and runtime
FROM ubuntu:latest AS build

# Update package lists and install the Go compiler
# -y flag automatically confirms installation to prevent build hangs
RUN apt-get update && apt-get install -y golang-go

# Set environment variable to handle Go projects without a go.mod file
ENV GO111MODULE=off

# Copy all files from the local directory into the container's root
COPY . .

# Compile the Go application
# CGO_ENABLED=0 ensures the binary is statically linked for better portability
RUN CGO_ENABLED=0 go build -o /app .

###########################################
# EXECUTION
###########################################
# Defines the command that runs when the container starts
ENTRYPOINT ["/app"]
