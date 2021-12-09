FROM golang:1-alpine as builder

ARG CREDENTIALS
ARG VERSION

RUN apk add --no-cache git ca-certificates unixodbc-dev build-base

RUN echo "$CREDENTIALS" > /root/.git-credentials && git config --global credential.helper store

RUN git clone --branch "v11.7.3" --single-branch --depth 1 \
    https://git.bullardisd.net/administrator/skyward-odbc.git /odbc && \
    rm /odbc/PGODBC.LIC

RUN go install github.com/korylprince/fileenv@v1.1.0
RUN go install "github.com/korylprince/bisd-device-checkout-server@$VERSION"


FROM alpine:3.15

RUN apk add --no-cache ca-certificates unixodbc libstdc++

COPY --from=builder /odbc /usr/local/lib/
COPY --from=builder /go/bin/fileenv /
COPY --from=builder /go/bin/bisd-device-checkout-server /
COPY setenv.sh /

RUN echo "[Progress]" > /etc/odbcinst.ini && echo "Driver=/usr/local/lib/pgoe27.so" >> /etc/odbcinst.ini

CMD ["/fileenv", "sh", "/setenv.sh", "/bisd-device-checkout-server"]
