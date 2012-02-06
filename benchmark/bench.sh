#!/bin/bash

go install hsat

BENCH=`ls benchmark/in-use/*.cnf`
#BENCH=`ls test/*.cnf`
BRANCH="ordered random vsids moms"
DBMS="queue berkmin" # DO NOT put none. swapping then dead
TEX="output.tex"

echo "\\documentclass{article}" > ${TEX}
echo "\\usepackage{geometry}" >> ${TEX}
echo "\\usepackage{morefloats}" >> ${TEX}
echo "\\usepackage{subfig}" >> ${TEX}
echo "" >> ${TEX}
echo "\\title{Base Solver Measurements}" >> ${TEX}
echo "\\author{Tim Kopp}" >> ${TEX}
echo "\\date{`date +%D`}" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{document}" >> ${TEX}
echo "\\maketitle" >> ${TEX}
echo "" >> ${TEX}

for bench in ${BENCH}; do
   ./bin/hsat -q -file=${bench} -b
done

echo "\\end{document}" >> ${TEX}
