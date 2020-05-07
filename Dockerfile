FROM golang:1.14-alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

# Add an appuser, so we don't run as root
RUN adduser -S -D -H -h /app appuser
USER appuser

CMD ["./main"]