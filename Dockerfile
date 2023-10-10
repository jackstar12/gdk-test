FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=1 go build -o /valid valid/main.go
RUN CGO_ENABLED=1 go build -o /invalid invalid/main.go

# Run
CMD ["/valid"]
