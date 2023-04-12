FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /dist


#RUN go build -o /main cmd/backend/main.go

CMD ["go", "run", "backend/app/main.go"]
