FROM golang:1.11 as build
ADD main.go /src/contribstats/main.go
ADD cmd /src/contribstats/cmd
ADD pkg /src/contribstats/pkg
ADD go.mod /src/contribstats/go.mod
ADD go.sum /src/contribstats/go.sum
WORKDIR /src/contribstats
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w"

### Now we build it
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /src/contribstats/contribstats /entrypoint
ADD .sample-contribstats.yml /config/.contribstats.yml
EXPOSE 8080
ENTRYPOINT ["/entrypoint"]
CMD ["--debug"]
