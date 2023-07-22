#!/bin/bash

out=$(./tests $*)
[ $? -eq 1 ] && echo "err: $out" || echo "out: $out"