FROM golang:1 as builder

RUN go install github.com/korylprince/fileenv@v1.1.0

FROM ubuntu:latest

ARG GO_PROJECT_NAME
ENV GO_PROJECT_NAME=${GO_PROJECT_NAME}

RUN apt-get update && apt-get install -y \
    unixodbc \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/fileenv /
COPY docker-entrypoint.sh /
COPY ${GO_PROJECT_NAME} /

# container expects pgoe27.so, libpgicu27.so, and PGODBC.LIC are located in /usr/local/lib/
RUN echo "[Progress]" > /etc/odbcinst.ini && echo "Driver=/usr/local/lib/pgoe27.so" >> /etc/odbcinst.ini

CMD ["/fileenv", "/docker-entrypoint.sh"]
