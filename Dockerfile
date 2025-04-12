FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]
# Use the following command to build the Docker image
# docker build -t my-go-app .
# Use the following command to run the Docker container
# docker run -p 8080:8080 my-go-app
# Use the following command to run the Docker container in detached mode