FROM golang:1.22-alpine

WORKDIR /app

COPY . ./

RUN go mod tidy

RUN go build -o /rinha-back ./cmd/main.go

CMD /rinha-back