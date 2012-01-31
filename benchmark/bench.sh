#!/bin/bash

go install hsat
FILE=$1

./bin/hsat -cpuprofile=cpu.prof -branch=vsids -dbms=queue -file=${FILE}
