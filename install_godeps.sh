#!/bin/bash

needed=()

packages=(\
"github.com/golang-collections/go-datastructures/bitarray" \
"github.com/veandco/go-sdl2/img" \
"github.com/veandco/go-sdl2/mix" \
"github.com/veandco/go-sdl2/sdl" \
"github.com/veandco/go-sdl2/ttf" \
"github.com/dave/jennifer" \
"go.uber.org/atomic")

do_pkg_check() {
	installed="$(go list '...' 2>/dev/null)"
	if [[ -z ${installed} ]]; then
		go list '...'
		exit 1
	fi
	for package in "${packages[@]}"; do
		printf "checking for %s ... " "${package}"
		if echo "${installed}" | grep -q "${package}"; then
			echo "installed"
			continue
		else
			echo "-------- not installed"
			needed+=($package)
		fi
	done
}

install() {
	go get -v "${needed[@]}"	
}

do_pkg_check
if [ ${#needed[@]} -ne 0 ]; then
	echo "######################################################"
	echo -e "need to install:\n\n" "${needed[@]}" "\n"
	echo "######################################################"
	set -x
	install
else
	echo "all needed go packages installed"
fi
