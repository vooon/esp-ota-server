FROM ghcr.io/goreleaser/goreleaser-cross:v1.25 AS builder

WORKDIR /build

ARG GIT_AUTHOR_NAME="Github CI"
ARG GIT_AUTHOR_EMAIL="build@github"

COPY . .

# hadolint global ignore=DL3059
RUN if [ ! -e ~/.gitconfig ]; then \
      echo "set dummy git user" && \
      git config --global user.email "$GIT_AUTHOR_EMAIL" && \
      git config --global user.name "$GIT_AUTHOR_NAME"; \
    fi

RUN goreleaser build --snapshot --single-target

FROM alpine

LABEL org.opencontainers.image.description="ESP OTA Server"

COPY --from=builder /build/dist/*/bin/esp-ota-server /espotad

ENV EOBIND=:8092
# ENV EOBASEURL=http://localhost:8092/
ENV EODATADIR=/data

VOLUME /data
EXPOSE 8092

CMD ["/espotad"]
