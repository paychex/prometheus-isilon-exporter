# B"H
FROM golang as builder 

WORKDIR /build

COPY . .

RUN make build

FROM alpine:3.13.5 as runtime

COPY --from=builder /build/bin/prometheus-isilon-exporter /build/entrypoint.sh /app/

# Add glibc to be able to run the isilon exporter
RUN apk --no-cache add ca-certificates wget \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk \
    && apk add glibc-2.28-r0.apk

RUN chmod +x /app/entrypoint.sh

CMD [ "/app/entrypoint.sh" ]