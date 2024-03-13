FROM alpine

VOLUME /data
WORKDIR /config

ENTRYPOINT ["/usr/bin/ipcmanview", "serve", "--dir=/data"]

ARG TARGETPLATFORM

COPY ./dist/${TARGETPLATFORM}/ipcmanview /usr/bin/ipcmanview
