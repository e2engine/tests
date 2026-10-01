FROM scratch

ARG SERVICE

COPY docker-bin/${SERVICE} /usr/bin/service

ENTRYPOINT ["/usr/bin/service"]