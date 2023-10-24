#!/bin/sh
cd $(dirname $(realpath $0)) || return
usage() {
    pwd
    echo "Scripts:"
    echo "$0 tidy"
    echo "    Tidy go modules."
}

tidy() {
    go mod tidy
}

if test -z "$1" -o -n "$(echo "$1" | grep -Ex '\-{0,2}h(elp)?')"; then
usage
exit
fi

case "$1" in
tidy) ;;
*)
usage
exit 1
;;
esac

$@

