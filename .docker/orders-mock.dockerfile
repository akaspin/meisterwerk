FROM --platform=$BUILDPLATFORM golang:bookworm AS build

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

USER root
WORKDIR /usr/src

COPY . ./
RUN --mount=type=ssh \
    mkdir /dist && go build -trimpath -o /dist/meisterwerk ./mock/orders


FROM --platform=$TARGETPLATFORM debian:12.6-slim

COPY --from=ghcr.io/tarampampam/curl:8.0.1 /bin/curl /curl
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /

COPY --from=build /dist/* /usr/bin/

CMD ["meisterwerk"]
