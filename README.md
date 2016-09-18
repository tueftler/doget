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
$ doget build -t [tag] .
```

This will resolve traits, downloading if necessary, and pass on the created Dockerfile to *docker build*.

## Using traits

Traits are just regular *Dockerfile*s stored somewhere on GitHub or BitBucket. You can reference them by using a `domain/vendor/repo[/dir][:version]` syntax. 

**Example:** `github.com/docker-library/php/7.0` will use the Dockerfile in [docker-library/php >> 7.0](https://github.com/docker-library/php/tree/master/7.0)

By default, this will check out the master branch. To reference a version, you can either use commit SHAs, branch names or tags and append them, e.g. `github.com/thekid/traits/xp:v1.0.0`.

## Authoring traits

As said, traits are nothing special. However, if you're creating Dockerfiles specifically designed for reuse, here are some things to keep in mind:

* Always add a *FROM* instruction to express what your Dockerfile extends from.
* If your traits provides an official base image, use *PROVIDES* and add its name.
* You can use *USE* to declare transitive dependencies. If you do so, you should reference a specific version, otherwise you risk problems at a later point.
* Think twice about adding an *ENTRYPOINT* or *CMD*, people will typically want to do this themselves.
* Test it using a continuous integration system like Travis CI

## Caching

DoGet caches downloaded traits inside the working directory. Their contents are stored zipped in a file called `doget_modules.zip`. To force a fresh download, simply remove this file.

You can check the file in to your SCM - this way, you can create repeatable builds even if the remote locations should not be reachable at build time.