FROM golang:1.10 as build
RUN go get -u github.com/golang/dep/cmd/dep
ADD main.go /go/src/github.com/thales-e-security/contribstats/main.go
ADD vendor /go/src/github.com/thales-e-security/contribstats/vendor
ADD cmd /go/src/github.com/thales-e-security/contribstats/cmd
ADD pkg /go/src/github.com/thales-e-security/contribstats/pkg
WORKDIR /go/src/github.com/thales-e-security/contribstats
#RUN dep ensure --vendor-only
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w"

### Now we build it
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/thales-e-security/contribstats/contribstats /entrypoint
ADD .sample-contribstats.yml /config/.contribstats.yml
EXPOSE 8080
ENTRYPOINT ["/entrypoint"]
CMD ["--debug"]
