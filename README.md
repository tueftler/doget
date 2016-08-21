DoGet
=====
Composes dockerfiles from traits like the one [here](https://github.com/thekid/gosu).

Setup
-----
Build the tool as follows:

```sh
$ go get gopkg.in/yaml.v2
$ go build github.com/tueftler/doget
```

Usage
-----
Start with this in a file called `Dockerfile.in`:

```dockerfile
FROM debian:jessie

INCLUDE github.com/thekid/gosu

CMD ["/bin/bash"]
```

Running the tool will give you this:

```sh
$ doget
> Running transform using doget\.doget.yml
> Fetching github.com/thekid/gosu: [####################] 0.74kB
Done
```

The resulting `Dockerfile` will now have the trait's contents in place of the *INCLUDE* instruction.

```dockerfile
FROM debian:jessie

# Included from github.com/thekid/gosu
ENV GOSU_VERSION 1.9

RUN set -x \
    && apt-get update && apt-get install -y 
    ...
    && apt-get purge -y --auto-remove ca-certificates wget

CMD ["/bin/bash"]
```

Versioning
----------
Versions can be added to includes just like tags in docker images:

```dockerfile
FROM debian:jessie

INCLUDE github.com/thekid/gosu:v1.0.0

CMD ["/bin/bash"]
```

Including subdirectories
------------------------
The following will include the `Dockerfile` from the subdirectory `7.0` rather than from the repository root.

```dockerfile
FROM debian:jessie

INCLUDE github.com/docker-library/php/7.0

RUN docker-php-ext-install bcmath

CMD /bin/bash
```