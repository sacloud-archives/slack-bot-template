FROM golang:1.9 AS builder
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>
LABEL MAINTAINER 'Kazumichi Yamamoto <yamamoto.febc@gmail.com>'

RUN  apt-get update && apt-get -y install \
        bash \
        git  \
        make \
        zip  \
      && apt-get clean \
      && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

ADD . /go/src/github.com/sacloud/slack-bot-template
WORKDIR /go/src/github.com/sacloud/slack-bot-template

RUN ["make", "deps", "clean", "build"]

#----------

FROM alpine:3.6
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>
LABEL MAINTAINER 'Kazumichi Yamamoto <yamamoto.febc@gmail.com>'

RUN set -x && apk add --no-cache --update zip ca-certificates
COPY --from=builder /go/src/github.com/sacloud/slack-bot-template/bin/slack-bot-template /usr/local/bin/
RUN chmod +x /usr/local/bin/slack-bot-template
ENTRYPOINT ["/usr/local/bin/slack-bot-template"]
EXPOSE 3000
