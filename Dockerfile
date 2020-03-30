FROM golang:alpine as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build -o main .

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
EXPOSE 8090
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]