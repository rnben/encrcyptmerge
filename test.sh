#!/bin/bash

out=$(./pg -action decrypt -new '{"a":"aa","b":"bb"}' -fd 'a' 2>&1)
[ "$out" = '{"a":"aa-decrypt","b":"bb"}' ] || echo "failed, out: $out"

out=$(./pg -action encrypt -new '{"a":"aa","b":"bb"}' -old '{"c":"c"}' -out e.log -fd 'a' 2>&1)
[ "$out" = '' ] || echo "failed, out: $out"
out=$(cat e.loga 2>&1)
[ "$out" = '{"a":"aa-encrypt","b":"bb","c":"c"}' ] || echo "failed, out: $out"