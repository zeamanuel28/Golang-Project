# --- Stage 1: Builder ---
# This stage compiles your Go application.
FROM golang:1.24-alpine AS builder 

# Set the working directory inside the container for the build process
WORKDIR /app

# Copy the go.mod and go.sum files first. This allows Docker to cache the
# go mod download step if your dependencies haven't changed.
COPY go.mod go.sum ./

# Download all necessary Go modules.
RUN go mod download

# Copy the rest of your application's source code.
# Assuming your main entry point is in the current directory (./)
COPY . .

# Build your Go application.
# CGO_ENABLED=0: Disables CGO, crucial for creating static binaries that run
#                on minimal base images like Alpine without external C dependencies.
# GOOS=linux GOARCH=amd64: Ensures the binary is built for a Linux environment on AMD64 architecture.
#                          This is important because Docker containers run Linux.
# -ldflags="-s -w": Strips debugging information and symbol tables, reducing binary size.
# -o /app/gocheck: Specifies the output path and name for the compiled executable.
#                  We'll name the binary after your project: 'gocheck'.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/gocheck .

# --- Stage 2: Final Image ---
# This stage creates a minimal image containing only the compiled application binary.
# Using a specific Alpine version (e.g., 3.22) is recommended for stability and security.
FROM alpine:3.22

# Set the working directory for the final running application
WORKDIR /app

# Copy the compiled binary from the 'builder' stage to the final, minimal image.
COPY --from=builder /app/gocheck .

# Expose the port your Gin application listens on.
# This serves as documentation within the image metadata.
# Adjust this if your Gin app listens on a different port (e.g., 3000, 8081).
EXPOSE 8080

# Define the command to run when the container starts.
# This now runs your 'gocheck' binary.
CMD ["./gocheck"]