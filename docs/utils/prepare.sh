#!/bin/bash

rm demo.html CONTRIBUTING.md .npmignore package.json package-lock.json gulpfile.js

cat ../reveal.js-template/utils/README-md-preffix.md README.md > README-new.md
mv README-new.md README.md
cp -r ../reveal.js-template/utils .
cp ../reveal.js-template/dist/extra.css dist
rm -rf test examples dist/theme/fonts css
