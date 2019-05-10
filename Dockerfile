FROM golang:1.11-alpine

WORKDIR /go/src/go_web_file_explorer

COPY . .

RUN go install -v ./...

CMD ["go_web_file_explorer"]
