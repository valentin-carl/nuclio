FROM alpine:3.19

WORKDIR /app

# install bash and go
RUN apk add bash go

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy code over and compile
COPY . /app
RUN go build grainworker.go

# leave it running
CMD ./grainworker