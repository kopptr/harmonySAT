#!/bin/bash

go install hsat

#BENCH=`ls benchmark/in-use/*.cnf`
BENCH=`ls test/*.cnf`
BRANCH="O R V M"
DBMS="Q B" # DO NOT put none. swapping then dead
BRANCHL="ordered random vsids moms"
DBMSL="queue berkmin" # DO NOT put none. swapping then dead
TEX="adapt-output.tex"
JSON=$1

echo "\\documentclass{article}" > ${TEX}
echo "\\usepackage[landscape]{geometry}" >> ${TEX}
echo "" >> ${TEX}
echo "\\title{Adaptive Solver Measurements}" >> ${TEX}
echo "\\author{Tim Kopp}" >> ${TEX}
echo "\\date{`date +%D`}" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{document}" >> ${TEX}
echo "\\maketitle" >> ${TEX}
echo "" >> ${TEX}
echo "Analysis used: ${JSON}" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{table}[ht!]" >> ${TEX}
echo "\\centering" >> ${TEX}
echo "\\begin{tabular}{|c||c|c|c|c||c|c|c|c|c|c|c|c||c|c|}\\hline" >> ${TEX}

echo -n "FILE & Bin & Tern & Horn & Def " >> ${TEX}
for b in ${BRANCH}; do
   for m in ${DBMS}; do
      echo -n "& \{${b},${m}\} " >> ${TEX}
   done
done
echo "& A & \\#\\\\\\hline\\hline" >> ${TEX}


for bench in ${BENCH}; do
   echo -n "${bench} & " >> ${TEX}
   ./bin/hsat -file=${bench} -a=${JSON} -b >> ${TEX}
   echo "" >> ${TEX}
done

echo "\\end{tabular}" >> ${TEX}
echo "\\end{table}" >> ${TEX}
echo "\\end{document}" >> ${TEX}
