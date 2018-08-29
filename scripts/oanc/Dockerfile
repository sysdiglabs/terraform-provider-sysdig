FROM golang:alpine as builder
WORKDIR /go/src/oanc
COPY . /go/src/oanc
RUN CGO_ENABLED=0 go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates shadow
RUN useradd -m oanc
WORKDIR /
COPY --from=builder /go/src/oanc/oanc .
USER oanc
CMD ["/oanc"]