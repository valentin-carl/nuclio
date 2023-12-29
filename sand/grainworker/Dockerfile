FROM alpine:3.19

WORKDIR /app

# install bash and go
RUN apk add bash go

# copy code over and compile
COPY . /app
RUN go build grainworker.go

# leave it running
CMD ./grainworker