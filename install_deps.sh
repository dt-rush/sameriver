#!/bin/bash

needed=()

do_pkg_check() {
	for package in "${packages[@]}"; do
		printf "checking for %s ... " "${package}"
		if ($1) | grep -q "${package}"; then
			echo "installed"
			continue
		else
			echo "-------- not installed"
			needed+=($package)
		fi
	done
}

apt_check() {
	packages=(\
		"libsdl2-dev" \
		"libsdl2-mixer-dev" \
		"libsdl2-image-dev" \
		"libsdl2-ttf-dev" \
		"libsdl2-gfx-dev")
	do_pkg_check "dpkg -l"
}

apt_install() {
	sudo apt-get update && \
	sudo apt-get install "${needed[@]}"
}

pacman_check() {
	packages=(\
		"sdl2" \
		"sdl2_mixer" \
		"sdl2_image" \
		"sdl2_ttf" \
		"sdl2_gfx")
	do_pkg_check "pacman -Qs"
}

pacman_install() {
	sudo pacman -Sy && \
	sudo pacman -S "${needed[@]}"
}

detect_pkgsystem() {
	if grep -q -Ei 'debian|buntu|mint' /etc/*release; then
		pkg_check="apt_check"
		pkg_install="apt_install"
	elif grep -q -i 'arch' /etc/issue; then
		pkg_check="pacman_check"
		pkg_install="pacman_install"
	fi
}

detect_pkgsystem
eval $pkg_check
if [ ${#needed[@]} -ne 0 ]; then
	echo "######################################################"
	echo -e "need to install:\n\n" "${needed[@]}" "\n"
	echo "######################################################"
	set -x
	echo "\$pkg_install is $pkg_install"
	eval $pkg_install
else
	echo "all needed apt packages installed"
fi
