# Runtime image only: GoReleaser builds the static binary and places it under
# $TARGETPLATFORM/ in the build context. distroless/static has no shell or package manager.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/nana /usr/local/bin/nana

USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/nana"]
