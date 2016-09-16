#!/bin/sh

branch=${1-none}
tag=${2-}

if [ "" != "$tag" ] ; then
  version=$(echo $tag | head -c 2)
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
  docker tag tueftler/doget:latest tueftler/doget:$version
  docker push tueftler/doget:$version
fi