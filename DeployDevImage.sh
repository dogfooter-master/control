#!/bin/sh
docker service rm dogfooter_dogfooter_control 2>/dev/null
docker stack deploy -c docker-stack.develop.yml dogfooter
