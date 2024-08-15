#!/usr/bin/env bash

# While bash may be not ideal choice as whimsical testing framework (ha-ha)
# let this file serve as a template/reminder of easy way for sending a number
# of requests with differing values using for and curl

baseurl=http://localhost:8082

# every line gives values for separate test, values are colon-separated
# as spaces separate array elements (tests) themselves
tests=(
  'Buddy:14:3'
  'Katty:15:2'
  'Piggy:12:4'
  'Peppi:10:1.5'
  'Fuzzy_Cow:9:2.7'
  'Pedros:17:7'
  'Clown_Gay:13.4:3.5'
  'Enola:8:0.7'
  'Swinoman:20:7'
  'Alco_Lower:6:10'
)


for testline in ${tests[@]} ; do
  readarray -d ':' -t values <<< $testline
  json="{\"name\":\"${values[0]}\", \"skill\":${values[1]}, \"latency\":${values[2]}}"
  echo $json
  res=`curl -s -d "$json" $baseurl/users`
  echo $res
done
