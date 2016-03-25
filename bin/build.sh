#!/bin/bash

if [[ "$(which gox)X" == "X" ]]; then
  echo "Please install gox. https://github.com/mitchellh/gox#readme"
  exit 1
fi


rm -f diego-enabler*

#gox -os linux -os windows -arch 386 --output="diego-enabler_{{.OS}}_{{.Arch}}"
#gox -os darwin -os linux -os windows -arch amd64 --output="diego-enabler_{{.OS}}_{{.Arch}}"
gox -os darwin -arch amd64 --output="diego-enabler_{{.OS}}_{{.Arch}}"

rm -rf out
mkdir -p out
mv diego-enabler* out/
