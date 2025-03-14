#!/usr/bin/env bash

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."
gofmt_files=$(find . -name '*.go' ! -name '*pb.go' | grep -v vendor | xargs gofmt -l)
if [[ -n ${gofmt_files} ]]; then
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0
