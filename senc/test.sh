#!/bin/sh
testf=main.go

go build
key=$(./senc --key-gen)
cat "$testf"  | ./senc -k "$key" -e >.test.enc
cat .test.enc | ./senc -k "$key" -d >.test.dec
diff "$testf" .test.dec
rm .test.enc .test.dec
