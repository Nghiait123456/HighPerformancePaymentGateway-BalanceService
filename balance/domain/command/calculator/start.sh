#!/bin/sh

set -e

if [ -x "$(command -v docker)" ]; then
    echo "Docker was installed"
else
    echo "Update Ubuntu and Install docker"
    sh install_docker_ubuntu.sh
fi

docker-compose up --build