#!/usr/bin/env bash

baseurl=http://localhost:8082

i=0
while [[ $i -lt 1000 ]] ; do
  i=$(($i+1))
  skill=$(($RANDOM % 14 + 7))
  latency=$(($RANDOM % 7 + 1))
  name="user_$i"
  json="{\"name\":\"$name\", \"skill\":$skill, \"latency\":$latency}"
  echo $json
  res=`curl -s -d "$json" $baseurl/users`
  echo $res
done
