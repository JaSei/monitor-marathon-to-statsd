FROM scratch

COPY monitor-marathon-to-statsd /

ENTRYPOINT ["/monitor-marathon-to-statsd"]
