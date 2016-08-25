#!/bin/sh

branch=${1-none}
tag=${2-}

if [ "master" != $branch ] ; then
  echo "Not pushing branch $branch"
  exit 0
fi

echo "Pushing :latest"
docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD" ;
docker push tueftler/doget ;

if [ "" != "$tag" ] && [ ${tag:1:1} -ge 1 ]; then
  tag=${tag:0:2};

  echo "Pushing :$tag"
  docker tag tueftler/doget:$tag tueftler/doget:latest;
  docker push tueftler/doget:$tag;
fi