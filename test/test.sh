#!/bin/bash

`go install hsat`

LIST=`ls test/*.cnf`
PASS=true

for file in ${LIST}; do
   echo -n "${file}: "
   RES=`./bin/hsat -q -file=./${file}`
   if [ "${RES}" != "SAT" ]; then
      PASS=false
   fi
   echo $RES
done

if $PASS; then
   echo "PASS"
else
   echo "FAIL"
fi
