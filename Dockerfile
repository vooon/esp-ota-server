FROM golang:alpine AS builder

WORKDIR /build

COPY . .
RUN go build -o espotad ./cmd/espotad

FROM alpine

COPY --from=builder /build/espotad /espotad

ENV EOBIND :8092
#ENV EOBASEURL http://localhost:8092/
ENV EODATADIR /data

VOLUME ["/data"]
EXPOSE 8092

ENTRYPOINT ["/espotad"]
