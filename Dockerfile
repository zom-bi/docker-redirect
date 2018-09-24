FROM golang:1

WORKDIR /go/src/github.com/zom-bi/docker-redirect
ADD . .
RUN \
    go get -tags="dev" -v github.com/zom-bi/docker-redirect && \
    go build -tags="netgo" -o redirect

FROM scratch

ENV \
    REDIRECT_CODE=302 \
    REDIRECT_LOG=false \
    REDIRECT_BEHIND_PROXY=true \
    REDIRECT_BEHIND_CLOUDFLARE=false \
    REDIRECT_URI="http://example.com"

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/zom-bi/docker-redirect/redirect /

ENTRYPOINT [ "/redirect" ]