#!/bin/bash

`go install hsat`

LIST=`ls test/*.cnf`
BRANCH_RULES="ordered random vsids"
PASS=true

for file in ${LIST}; do
        for rule in ${BRANCH_RULES} ; do
           echo -n "${file}, ${rule}: "
           RES=`./bin/hsat -q -file=./${file} -branch=${rule}`
           if [ "${RES}" != "SAT" ]; then
              PASS=false
           fi
           echo $RES
        done
done

if $PASS; then
   echo "PASS"
else
   echo "FAIL"
fi
