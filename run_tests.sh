#!/bin/sh

# -cpu=$(sysctl -n hw.ncpu)

go test -bench=. -timeout 3600s # -cpu=1,2,4,8

