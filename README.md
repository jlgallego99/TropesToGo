# Trabajo de Fin de Máster: Scraping de tvtropes.org 

### Autor: Jose Luis Gallego Peña
### Tutor: Juan Julián Merelo Guervós
___

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![License: CC BY-NC-SA 3.0](https://img.shields.io/badge/License-CC%20BY--NC--SA%203.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-sa/4.0/)

All [HTML files](/tropestogo/service/scraper/resources/) and all data on the generated datasets belong to [TvTropes](https://tvtropes.org/) and are under a Creative Commons Attribution-NonCommercial-ShareAlike 3.0 Unported License.

The slides are available on GitHub pages, made with reveal.js using a [template](https://github.com/JJ/reveal.js-template)

---

La documentación de este proyecto está realizada con `LaTeX`, por lo
tanto para generar el archivo PDF necesitaremos instalar `TeXLive` en
nuestra distribución.

Una vez instalada, tan solo deberemos situarnos en el directorio `doc` y ejecutar:

`
$ pdflatex proyecto.tex
`

Seguido por

    bibtex proyecto
    
y de nuevo

    pdflatex proyecto.tex

O directamente

    make
    
(que habrá que editar si el nombre del archivo del proyecto cambia)