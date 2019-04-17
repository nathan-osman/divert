FROM golang:latest

ENV CGO_ENABLED=0

ADD . /src/

RUN go generate
RUN go build


FROM scratch

COPY --from=builder /src/divert /usr/local/bin
