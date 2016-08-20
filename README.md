DoGet
=====

Composes dockerfiles from traits like the one [here](https://github.com/thekid/gosu):

```dockerfile
FROM debian:jessie

INCLUDE github.com/thekid/gosu

CMD ["/bin/bash"]
```

Running the tool will give you this:

```sh
$ go build github.com/tueftler/doget

$ doget -in Dockerfile.in
> Fetching github.com/thekid/gosu: [####################] 0.74kB
Done

FROM debian:jessie

# Included from github.com/thekid/gosu
ENV GOSU_VERSION 1.9

RUN set -x \
    && apt-get update && apt-get install -y 
    ...
    && apt-get purge -y --auto-remove ca-certificates wget

CMD ["/bin/bash"]
```