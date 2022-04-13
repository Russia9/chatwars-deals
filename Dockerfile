# Build container
FROM golang:1.18-bullseye AS build

# Set build workdir
WORKDIR /app

# Copy app sources
COPY . .

# Build app
RUN go build -o bin .

# ---
# Production container
FROM debian:bullseye-slim

# Set app workdir
WORKDIR /app

# Copy binary
COPY --from=build /app/bin .

# Run app
CMD ["./bin"]
