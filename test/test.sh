#!/bin/bash

`go install hsat`

LIST=`ls test/*.cnf`
for file in ${LIST}; do
   echo -n "${file}: "
   ./bin/hsat -q -file=./${file}
done

