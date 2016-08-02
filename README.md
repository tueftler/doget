DoGet
=====

Composes dockerfiles

```dockerfile
FROM debian:jessie

INCLUDE github.com/thekid/gosu

CMD ["/bin/bash"]
```

Running the tool will give you this:

```sh
$ go run doget.go
FROM debian:jessie

# Included from github.com/thekid/gosu

ENV GOSU_VERSION 1.9
RUN set -x \
    && apt-get update && apt-get install -y ...
    && apt-get purge -y --auto-remove ca-certificates wget

CMD ["/bin/bash"]
```