# minimalist runtime
FROM --platform=$BUILDPLATFORM alpine:3.17.1
# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="media-normalizer" \
    org.label-schema.description="" \
    org.label-schema.url="" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/MrEhbr/media-normalizer" \
    org.label-schema.vendor="Aleksei Burmistrov" \
    org.label-schema.version=$VERSION \
    org.label-schema.schema-version="1.0" \
    org.label-schema.cmd="docker run -i -t --rm MrEhbr/media-normalizer" \
    org.label-schema.help="docker exec -it $CONTAINER media-normalizer --help"
COPY            media-normalizer /bin/
ENTRYPOINT      ["/bin/media-normalizer"]
#CMD             []
