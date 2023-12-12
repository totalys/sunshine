FROM golang:1.21.1 as go-builder

ARG CI_VERSION
ARG CI_COMMIT_SHA

COPY . /go/src/sunshine

RUN true

WORKDIR /go/src/sunshine
RUN ls -l

ENV SUNNY_CI_VERSION=${CI_VERSION}
ENV SUNNY_CI_COMMIT_SHA=${CI_COMMIT_SHA}

RUN echo building sunshine weather app

RUN go build -o /go/bin/sunshine --ldflags "-X main.version=${SUNNY_CI_VERSION} -X main.commitID=${CI_COMMIT_SHA}" cmd/sunshine/main.go

FROM golang:1.21.1

COPY --from=go-builder /go/bin/sunshine                                /go/bin/sunshine
COPY --from=go-builder /go/src/sunshine/configs/config.json            /go/bin/configs/
COPY --from=go-builder /go/src/sunshine/configs/local_city_source.json /go/bin/configs/

WORKDIR /go/bin

CMD ./sunshine