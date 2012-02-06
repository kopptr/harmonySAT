#!/bin/bash

go install hsat

BENCH=`ls benchmark/in-use/*.cnf`
#BENCH=`ls test/*.cnf`
BRANCH="ordered random vsids moms"
DBMS="queue berkmin" # DO NOT put none. swapping then dead

echo "\\documentclass{article}"
echo "\\usepackage{geometry}"
echo "\\usepackage{morefloats}"
echo ""
echo "\\title{Formula Set Analysis}"
echo "\\author{Tim Kopp}"
echo "\\date{`date +%D`}"
echo ""
echo "\\begin{document}"
echo "\\maketitle"
echo ""

for bench in ${BENCH}; do
   echo "\\begin{table}[ht!]"
   echo "\\centering"
   ./bin/hsat -e -file=${bench}
   echo "\\caption{${bench}}"
   echo "\\label{tab:${bench}}"
   echo "\\end{table}"
   echo ""
done

echo "\\end{document}"
