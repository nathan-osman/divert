FROM golang:latest

ENV CGO_ENABLED=0

ADD . /src

WORKDIR /src

RUN go generate
RUN go build


FROM scratch

COPY --from=0 /src/divert /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/divert"]
