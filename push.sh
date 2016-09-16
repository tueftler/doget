#!/bin/bash

branch=${1-none}
tag=${2-}

if [ "" != "$tag" ] && [ ${tag:1:1} -ge 1 ]; then
  version=${tag:0:2}
elif [ "master" != $branch ] ; then
  echo "Not pushing branch $branch"
  exit 0
else
  version=
fi

echo "Pushing :latest"
docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
docker push tueftler/doget

if [ "" != "$version" ] ; then
  echo "Pushing :$version"
  docker tag tueftler/doget:$version tueftler/doget:latest
  docker push tueftler/doget:$version
fi