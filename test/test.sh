#!/bin/bash

`go install hsat`

LIST=`ls test/*.cnf`
BRANCH_RULES="ordered random vsids"
CDBMS="none queue"
PASS=true

for file in ${LIST}; do
  for rule in ${BRANCH_RULES} ; do
     for cdbms in ${CDBMS}; do
        echo -n "${file}, ${rule}, ${cdbms}: "
        RES=`./bin/hsat -q -file=./${file} -branch=${rule} -dbms=${cdbms}`
        if [ "${RES}" != "SAT" ]; then
           PASS=false
        fi
        echo $RES
     done
  done
done

if $PASS; then
   echo "PASS"
else
   echo "FAIL"
fi
