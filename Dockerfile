FROM alpine
ADD gomysql /
USER nobody
ENTRYPOINT ["/gomysql"]
