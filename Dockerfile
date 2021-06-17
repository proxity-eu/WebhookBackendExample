FROM golang:alpine AS build
RUN apk update && apk add ca-certificates
WORKDIR /src
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -mod=readonly -ldflags="-w -s" -o server

FROM alpine:3.13
RUN apk update && apk add ca-certificates
COPY --from=build /src/server .
ENTRYPOINT ["./server"]