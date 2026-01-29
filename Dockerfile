FROM europe-docker.pkg.dev/kyma-project/prod/external/library/golang:1.23.1-alpine3.20 AS builder

ARG DOCK_PKG_DIR=/go/src/github.com/kyma-project/eventing-publisher-proxy

WORKDIR $DOCK_PKG_DIR
COPY . $DOCK_PKG_DIR

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on GOFIPS140=v1.0.0 go build -o eventing-publisher-proxy ./cmd/eventing-publisher-proxy

FROM gcr.io/distroless/static:nonroot
LABEL source = git@github.com:kyma-project/kyma.git
USER nonroot:nonroot

WORKDIR /
COPY --from=builder /go/src/github.com/kyma-project/eventing-publisher-proxy/eventing-publisher-proxy .

ENV GODEBUG=fips140=only,tlsmlkem=0
ENTRYPOINT ["/eventing-publisher-proxy"]
