#!/bin/bash

go install hsat

BENCH=`ls benchmark/in-use/*.cnf`
JSON=$1
rm -f output.txt

for bench in ${BENCH}; do
   echo "./bin/hsat -file=${bench} -a=${JSON}" >> output.txt
   ./bin/hsat -file=${bench} -a=${JSON} >> output.txt
   echo "" >> output.txt
done
