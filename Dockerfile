FROM golang:1.14-alpine

ARG WORKDIR

RUN mkdir /app
ADD . /app/
WORKDIR ${WORKDIR}
RUN go build -o main .

# Add an appuser, so we don't run as root
RUN adduser -S -D -H -h /app appuser
USER appuser

CMD ["./main"]