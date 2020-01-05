#!/usr/bin/env bash
mkdir -p /data/octocv/
cp -r /vagrant/data/* /data/octocv/
cd /data/octocv/
find ./deployments/local -type f -name "*.sh" -print0 | xargs -0 sudo dos2unix
find ./deployments/local/mongo -type f -name "*.sh" -print0 | xargs -0 sudo dos2unix
find ./deployments/local/nginx -type f -name "*.sh" -print0 | xargs -0 sudo dos2unix
docker-compose up -d mongodb redis