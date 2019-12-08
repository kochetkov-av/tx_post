#!/usr/bin/env bash

#docker run -v $(pwd):/etc/newman -t postman/newman TxPost.postman_collection.json

docker run \
    --rm \
    --network txpost_default \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run TxPost.postman_collection.json \
    --env-var host=app:4000 \
    --folder "Init balance"

docker run \
    --rm \
    --network txpost_default \
    -v $(pwd)/test:/etc/newman \
    -t postman/newman:alpine \
    run TxPost.postman_collection.json \
    --env-var host=app:4000 \
    --folder "Change balance" \
    --iteration-count 100
