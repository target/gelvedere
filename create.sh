#!/bin/bash

name="gelvedere"

go get github.com/mitchellh/gox
gox -osarch="linux/amd64 darwin/amd64" github.com/target/gelvedere/cmd/gelvedere

for n in $(ls *$name*)
do
  mv $n $name
  platform=$(echo $n | awk -F[=_] '{print $2}')
  arch=$(echo $n | awk -F[=_] '{print $3}')
  tar czf ${name}-${platform}-${arch}.tgz ${name}
  rm -f ${name}
done