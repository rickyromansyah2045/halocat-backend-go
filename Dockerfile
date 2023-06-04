FROM golang:1.18
WORKDIR /app
ENV GOPATH /go
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./out/dist .
CMD ./out/dist