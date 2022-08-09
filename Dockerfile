FROM registry.ci.openshift.org/openshift/release:golang-1.18 as builder

WORKDIR /workspace

COPY . .

RUN go vet ./...

RUN go build -a -o tenant-manager cmd/main.go

FROM registry.access.redhat.com/ubi9-minimal

COPY --from=builder /workspace/tenant-manager /usr/bin/

EXPOSE 8080

ENTRYPOINT ["/usr/bin/tenant-manager"]