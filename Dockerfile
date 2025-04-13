# Build-Stage
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy the source code
COPY . .

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ files
RUN templ generate

# Install build dependencies
RUN apk add --no-cache gcc musl-dev curl

# Get the latest version from GitHub API and save it to version.txt
RUN curl -s https://api.github.com/repos/axzilla/templui/releases/latest | grep tag_name | cut -d '"' -f 4 > version.txt || echo "unknown" > version.txt

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# Deploy-Stage
FROM alpine:3.20.2
WORKDIR /app

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Set environment variable for runtime
ENV GO_ENV=production

# Copy the binary and version file
COPY --from=build /app/main .
COPY --from=build /app/version.txt .

# Expose the port
EXPOSE 8090

# Command to run
CMD ["./main"]
