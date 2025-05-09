FROM golang:1.24.3
COPY . /app
WORKDIR /app
RUN go build .
CMD ["go", "run", "."]