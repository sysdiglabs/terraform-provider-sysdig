#!/bin/bash

# set -x # DEBUG

pushd ../website/docs &> /dev/null

cat index.md | sed '1,/---$/d'

printf "\n\n# Data Sources\n\n"

for md in `find ./d -name "*.markdown"`; do
	cat "$md" | sed '1,/---$/d' | sed 's/^#/##/g'
	printf "\n\n---\n\n"
done

printf "\n\n# Resources\n\n"

for md in `find ./r -name "*.md"`; do
	cat "$md" | sed '1,/---$/d' | sed 's/^#/##/g'
	printf "\n\n---\n\n"
done

popd &> /dev/null
