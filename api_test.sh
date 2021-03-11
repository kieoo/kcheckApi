#!/bin/sh

url="192.168.56.1:7001/total_check/Junit_xml"
postF=""

for chart in `ls -l |grep "^d" |awk '{print $9}'`;
do
  for f in `find $chart -name "*.*"`;
  do
    postF="$postF -F \"files=@$f;fileName=$f\""
  done
curl_cmd="curl $url$postF"
eval "$curl_cmd 2>&1" | tee report_$chart.xml
done