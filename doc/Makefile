DOC := proyecto.tex

all: doc

doc: $(DOC:.tex=.pdf)

all: proyecto.pdf

%.pdf: %.tex FORCE
	pdflatex --shell-escape $< && bibtex $* && pdflatex --shell-escape $< && pdflatex --shell-escape $<

FORCE: 


