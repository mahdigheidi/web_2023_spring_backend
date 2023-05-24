FROM golang:1.20

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
COPY pb/ pb/
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build 
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

CMD ["/main"]