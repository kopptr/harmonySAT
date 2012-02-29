#!/bin/bash

BENCH=`ls benchmark/test/*.cnf`
TEX="others-output.tex"

echo "\\documentclass{article}" > ${TEX}
echo "\\usepackage[paperwidth=15in,paperheight=8.5in]{geometry}" >> ${TEX}
echo "" >> ${TEX}
echo "\\title{Adaptive Solver Measurements}" >> ${TEX}
echo "\\author{Tim Kopp}" >> ${TEX}
echo "\\date{`date +%D`}" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{document}" >> ${TEX}
echo "\\maketitle" >> ${TEX}
echo "" >> ${TEX}
echo "\\begin{table}[ht!]" >> ${TEX}
echo "\\centering" >> ${TEX}
echo "\\begin{tabular}{|c||c|c|c|c|}\\hline" >> ${TEX}

echo "FILE & zChaff & Clasp & MiniSAT & MiniSATs\\\\\\hline\\hline" >> ${TEX}
for bench in ${BENCH}; do
echo -n "${bench} & " >> ${TEX}
   /usr/bin/time -f %E -o ${TEX} -a other-solvers/zchaff64/zchaff ${bench}
   echo -n " & " >> ${TEX}
   /usr/bin/time -f %E -o ${TEX} -a clasp ${bench}
   echo -n " & " >> ${TEX}
   /usr/bin/time -f %E -o ${TEX} -a other-solvers/minisat/core/minisat ${bench}
   echo -n " & " >> ${TEX}
   /usr/bin/time -f %E -o ${TEX} -a other-solvers/minisat/simp/minisat ${bench}
   echo "\\\\\\hline" >> ${TEX}
done

echo "\\end{tabular}" >> ${TEX}
echo "\\end{table}" >> ${TEX}
echo "\\end{document}" >> ${TEX}

