FROM golang:1.18.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.dca ./

RUN go build -o /juanita

CMD [ "/juanita" ]
