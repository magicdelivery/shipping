#!/bin/sh

cat $1 | redis-cli -h $REDIS_HOST -p $REDIS_PORT