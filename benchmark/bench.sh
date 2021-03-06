#!/bin/bash

go install hsat

BENCH=`ls benchmark/analysis/*.cnf`
#BENCH=`ls test/*.cnf`
BRANCH="O R V M T"
DBMS="Q B" # DO NOT put none. swapping then dead
BRANCHL="ordered random vsids moms vmtf"
DBMSL="queue berkmin" # DO NOT put none. swapping then dead
TEX="output.tex"

echo "\\documentclass{article}" > ${TEX}
echo "\\usepackage[paperheight=8.5in, paperwidth=15in]{geometry}" >> ${TEX}
echo "\\usepackage{fullpage}" >> ${TEX}
echo "" >> ${TEX}
echo "\\title{Base Solver Measurements}" >> ${TEX}
echo "\\author{Tim Kopp}" >> ${TEX}
echo "\\date{`date +%D`}" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{document}" >> ${TEX}
echo "\\maketitle" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{table}[ht!]" >> ${TEX}
echo "\\centering" >> ${TEX}
echo "\\begin{tabular}{|c||c|c|c|c|c|c|c|c||c|c|c|c|c|c|c|c|c|c|}\\hline" >> ${TEX}

echo -n "FILE & Bin & Tern & Horn & Def & HI & hi & lo & LO " >> ${TEX}
for b in ${BRANCH}; do
   for m in ${DBMS}; do
      echo -n "& \{${b},${m}\} " >> ${TEX}
   done
done
echo "\\\\\\hline\\hline" >> ${TEX}


for bench in ${BENCH}; do
   echo -n "${bench} & " >> ${TEX}
   ./bin/hsat -q -file=${bench} -b
   echo "" >> ${TEX}
   echo "finished ${bench}"
done

echo "\\end{tabular}" >> ${TEX}
echo "\\end{table}" >> ${TEX}
echo "\\end{document}" >> ${TEX}
