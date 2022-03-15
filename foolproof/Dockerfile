FROM golang:1.17.1 as builder

WORKDIR /app

# Fetch the go mod deps first so docker build can cache these
# making the builds faster.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN CGO_ENABLED=0 go build -o /api cmd/api/main.go

# Start a new stage from scratch:
FROM alpine:3.11.3

COPY --from=builder /api /api

CMD /api
