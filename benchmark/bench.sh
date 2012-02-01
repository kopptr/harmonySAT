#!/bin/bash

go install hsat

BENCH=`ls benchmark/*.cnf`
#BENCH=`ls test/*.cnf`
BRANCH="ordered random vsids moms"
DBMS="queue berkmin" # DO NOT put none. swapping then dead

echo "\\documentclass{article}"
echo "\\usepackage{geometry}"
echo "\\usepackage{morefloats}"
echo "\\usepackage{subfig}"
echo ""
echo "\\title{Base Solver Measurements}"
echo "\\author{Tim Kopp}"
echo "\\date{`date +%D`}"
echo ""
echo "\\begin{document}"
echo "\\maketitle"
echo ""

for bench in ${BENCH}; do
   echo "\\begin{table}[ht!]"
   echo "\\centering"
   echo "\\subfloat[][]{"
   ./bin/hsat -a -file=${bench}
   echo "}"
   echo "\\subfloat[][]{"
   echo "\\begin{tabular}{|c|c||c|}\\hline"
   echo "Branch & DBMS & Time\\\\\\hline\\hline"
   for rule in ${BRANCH}; do
      for dbms in ${DBMS}; do
         echo -n "${rule} & ${dbms} & "
         RES=`ulimit -t 1200; /usr/bin/time -f %E ./bin/hsat -q -file=${bench} -branch=${rule} -dbms=${dbms} 2>&1 >/dev/null`
         echo -n ${RES}
         echo "\\\\\\hline"
      done
   done
   echo "\\end{tabular}"
   echo "}"
   echo "\\caption{${bench}}"
   echo "\\label{tab:${bench}}"
   echo "\\end{table}"
   echo ""
done

echo "\\end{document}"
