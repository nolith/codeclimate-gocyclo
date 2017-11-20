FROM golang:1.8-alpine AS build

RUN apk add -U git jq
RUN go get github.com/fzipp/gocyclo
RUN cd $GOPATH/src/github.com/fzipp/gocyclo && git describe --tag --always > /gocyclo_version

WORKDIR /go/src/gitlab.com/nolith/codeclimate-gocyclo
COPY . .

RUN go build

RUN cat engine.json | jq ".version = \"$(cat /gocyclo_version)-$(git describe --tag --always)\"" > /engine.json

FROM alpine:3.6

LABEL maintainer="Alessio Caiazza" 

RUN adduser -u 9000 -D app
COPY --from=build /go/bin/gocyclo /usr/bin/gocyclo
COPY --from=build /go/src/gitlab.com/nolith/codeclimate-gocyclo/codeclimate-gocyclo /usr/bin/codeclimate-gocyclo
COPY --from=build /engine.json /engine.json

WORKDIR /code
VOLUME /code

USER app

CMD ["/usr/bin/codeclimate-gocyclo"]
