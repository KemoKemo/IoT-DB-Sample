FROM golang:1.11.1

WORKDIR /go/src/app
COPY . .

ENV GIN_MODE release
ENV DBaddr http://path/to/mongodb
ENV Port :9000

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app", "-dbaddr=$DBaddr", "-addr=$Port"]
