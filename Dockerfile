FROM golang:alpine as builder
RUN mkdir -p /go/src/temperature-backend
WORKDIR /go/src/temperature-backend
ADD . ./
RUN go build -o app

FROM alpine
COPY --from=builder /go/src/temperature-backend/app /app
EXPOSE 8080
RUN ls -lah
CMD ["/app"]
