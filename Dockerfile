FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/gost-heroku /opt/bin/gost-heroku

CMD ["/opt/bin/gost-heroku"]

