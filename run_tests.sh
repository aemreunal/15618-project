#!/bin/sh

# -cpu=$(sysctl -n hw.ncpu)

go test -bench=. -timeout 36000s # -cpu=1,2,4,8

