FROM golang:1.9 AS build

WORKDIR /go/src/github.com/target/gelvedere
ENV CGO_ENABLED=0 GOOS=linux
RUN mkdir vendor client cmd model version
COPY vendor ./vendor/
COPY client ./client/
COPY cmd ./cmd/
COPY model ./model/
COPY version ./version/
RUN go build github.com/target/gelvedere/cmd/gelvedere

FROM scratch

COPY --from=build /go/src/github.com/target/gelvedere/gelvedere /bin/
ENTRYPOINT ["/bin/gelvedere"]
