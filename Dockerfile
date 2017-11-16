FROM golang:1.8-alpine AS build

RUN apk add -U git
RUN go get github.com/fzipp/gocyclo

WORKDIR /go/src/gitlab.com/nolith/codeclimate-gocyclo
COPY . .

RUN go build

FROM alpine:3.6

LABEL maintainer="Alessio Caiazza" 

RUN adduser -u 9000 -D app
COPY --from=build /go/bin/gocyclo /usr/bin/gocyclo
COPY --from=build /go/src/gitlab.com/nolith/codeclimate-gocyclo/codeclimate-gocyclo /usr/bin/codeclimate-gocyclo

WORKDIR /code
VOLUME /code

USER app

CMD ["/usr/bin/codeclimate-gocyclo"]
