FROM golang:1.18
WORKDIR /app
ENV GOPATH /go
COPY go.mod .
COPY go.sum .
COPY .env /app
RUN go mod download
COPY . .
RUN go build -o ./out/dist .
CMD ./out/dist