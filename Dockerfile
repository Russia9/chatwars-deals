FROM golang:1.18

# Set app workdir
WORKDIR /go/src/app

# Copy application sources
COPY . .

# Build app
RUN go build -o bin .

# Run app
CMD ["./bin"]