FROM golang:1.8

WORKDIR /go/src/mightbot
COPY . .

RUN go install    # "go install -v ./..."
EXPOSE 6600

CMD ["run"] # ["mightbot"]