FROM golang:1.14.2 as build
ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/sigs.k8s.io/pingmesh-agent
COPY go.mod .
COPY go.sum .
COPY vendor vendor
RUN go mod download

COPY pkg pkg
COPY cmd cmd

ARG GOARCH
ARG LDFLAGS
RUN go build -mod=readonly -ldflags "$LDFLAGS" -o /pingmesh-agent $PWD/cmd/pingmesh-agent

FROM golang:1.14.2

COPY --from=build pingmesh-agent /

USER 65534

ENTRYPOINT ["/pingmesh-agent"]