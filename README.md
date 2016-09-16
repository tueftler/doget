DoGet
=====
[![Build Status on TravisCI](https://secure.travis-ci.org/tueftler/doget.png)](http://travis-ci.org/tueftler/doget)

Dockerfiles can only inherit from one source. What if you want to want to your favorite application on your own base image? Then you're down to copy&pasting Dockerfile code all over the place. 

**DoGet solves this**. Think of DoGet as "compiler-assisted copy&paste". Here's an example:

```dockerfile
FROM corporate.example.com/images/debian:8        # Our own base image

PROVIDES debian:jessie                            # ...which is basically Debian Jessie

USE github.com/docker-library/php/7.0             # On top of that, use official PHP 7.0

# <Insert app here>
```

DoGet extends Dockerfile syntax with `PROVIDES` and `USE` instructions. The `USE` instruction downloads the Dockerfile from the specified location and includes its contents (as if it had been copy&pasted). The `FROM` instruction in included files is checked for compatibility. The `PROVIDES` instruction serves the purpose of satisfying the compatibility check.

To use the tool, copy the above into a file called `Dockerfile.in` and type:

```sh
$ doget build
```

This will resolve traits, downloading if necessary, and pass on the created Dockerfile to *docker build*.