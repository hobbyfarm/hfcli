FROM registry.suse.com/bci/golang:latest as build

WORKDIR /usr/src/hfcli

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

FROM registry.suse.com/bci/bci-base:latest

COPY --from=build /usr/src/hfcli/hfcli /

ENTRYPOINT ["/hfcli"]
