FROM golang:1.18

# Set app workdir
WORKDIR /go/src/app

# Copy application sources
COPY . .

# Build app
RUN go build -o app .

# Run app
CMD ["./app"]