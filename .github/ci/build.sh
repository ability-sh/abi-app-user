#!/bin/sh

DIR=$(dirname $0)

echo $DIR

docker run --rm -v $DIR/../..:/workdir -w "/workdir" docker.io/nginx/unit:1.27.0-go1.18 /workdir/.github/workflows/ci/docker_build.sh
