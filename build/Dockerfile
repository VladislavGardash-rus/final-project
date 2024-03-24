FROM golang:1.20 as build

ENV BIN_FILE /opt/final-project/final-project-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/*

FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="final-project"
LABEL MAINTAINERS="gardashvs@comlogic.ru"

ENV BIN_FILE "/opt/final-project/final-project-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/final-project/config.json
COPY ./configs/config.json ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
