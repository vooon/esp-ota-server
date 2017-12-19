FROM ubuntu:xenial

ENV EOBIND :8092
#ENV EOBASEURL http://localhost:8092/
ENV EODATADIR /data

VOLUME ["/data"]
EXPOSE 8092

COPY app /app
ENTRYPOINT ["/app"]
