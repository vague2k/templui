# Build-Stage
FROM golang:1.24-alpine AS build
WORKDIR /app

# Copy the source code
COPY . .

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ files
RUN templ generate

# Install Node.js and npm for TailwindCSS
RUN apk add --no-cache nodejs npm

# Install TailwindCSS
RUN npm install -g tailwindcss

# Build CSS
RUN tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server/main.go

# Deploy-Stage
FROM alpine:3.20.2
WORKDIR /app

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Set environment variable for runtime
ENV GO_ENV=production

# Copy the binary from the build stage
COPY --from=build /app/main .

# Copy static assets and CSS
COPY --from=build /app/static ./static
COPY --from=build /app/assets ./assets

# Expose the port your application runs on
EXPOSE 8090

# Command to run the application
CMD ["./main"]