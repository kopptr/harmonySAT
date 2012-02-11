#!/bin/bash

go install hsat

DIR=$1
BENCH=`ls ${DIR}`
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
echo "\\begin{table}[ht!]"
echo "\\centering"
echo "\\begin{tabular}{|c||c|c|c|c|}\\hline"
echo "File & Bin & Tern & Horn & Def\\\\\\hline\\hline"

for bench in ${BENCH}; do
   echo -n "${bench} & "
   ./bin/hsat -e -file=${DIR}/${bench}
   echo "\\\\\\hline"
done

echo "\\end{tabular}"
echo "\\end{table}"
echo ""

echo "\\end{document}"
