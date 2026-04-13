###########################################
# STAGE 1: BUILD STAGE
###########################################
# Using Ubuntu as the heavy base image to provide the Go compiler environment
FROM ubuntu AS build

# Update the package manager and install the Go compiler
# -y is used to skip manual confirmation during the build process
RUN apt-get update && apt-get install -y golang-go

# Disable Go Modules to support older project structures or local builds
ENV GO111MODULE=off

# Copy the entire source code from the current local directory to the container
COPY . .

# Compile the application into a statically linked binary named 'app'
# CGO_ENABLED=0 ensures the binary has no dependencies on external C libraries
RUN CGO_ENABLED=0 go build -o /app .

############################################
# STAGE 2: FINAL RUNTIME (THE MAGIC)
############################################
# Use 'scratch'—the smallest possible base image (zero bytes)
FROM scratch

# Precisely copy only the compiled binary from the 'build' stage
# This keeps the final image size extremely small and secure
COPY --from=build /app /app

# Define the binary as the main process to execute upon container start
ENTRYPOINT ["/app"]
