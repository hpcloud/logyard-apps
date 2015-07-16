FROM docker-internal.stackato.com/stackato/base-redis:master

ADD . /source
ADD stackon.json /

CMD ["redis-server", "/source/redis-server.conf"]

EXPOSE 6464
VOLUME ["/var/stackato/data/applog_redis"]
