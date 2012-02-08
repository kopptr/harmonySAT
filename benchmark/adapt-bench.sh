#!/bin/bash

go install hsat

BENCH=`ls benchmark/in-use/*.cnf`
JSON=$1
rm -f output.txt

for bench in ${BENCH}; do
   echo "./bin/hsat -file=${bench} -a=${JSON}" >> output.txt
   ulimit -t 1200; /usr/bin/time ./bin/hsat -file=${bench} -a=${JSON} >> output.txt 2>> output.txt
   echo "" >> output.txt
done
