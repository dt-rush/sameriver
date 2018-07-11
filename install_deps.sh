#!/bin/bash

needed=()

do_pkg_check() {
	for package in "${packages[@]}"; do
		printf "checking for $package ... "
		if [ "$(($1) | grep $package)" ]; then
			echo "installed"
			continue
		else
			echo "not installed"
			needed+=($package)
		fi
	done
}

apt_check() {
	packages=("libsdl2-dev" \
		"libsdl2-mixer-dev" \
		"libsdl2-image-dev" \
		"libsdl2-ttf-dev" \
		"libsdl2-gfx-dev")
	do_pkg_check "dpkg -l"
}

apt_install() {
	sudo apt-get update && \
	sudo apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
}

pacman_check() {
	packages=("sdl2" \
		"sdl_mixer" \
		"sdl_image" \
		"sdl_ttf" \
		"sdl_gfx")
	do_pkg_check "pacman -Qs"
}

pacman_install() {
	sudo pacman -Sy && \
	sudo pacman -S sdl2{,_mixer,_image,_ttf,_gfx}
}

detect_pkgsystem() {
	if [ "$(grep -Ei 'debian|buntu|mint' /etc/*release)" ]; then
		pkg_check="apt_check"
		pkg_install="apt_install"
	elif [ "$(grep -i 'arch' /etc/issue)" ]; then
		pkg_check="pacman_check"
		pkg_install="pacman_install"
	fi
}

detect_pkgsystem
($pkg_check)
if [ ${#needed[@]} -ne 0 ]; then
	echo "######################################################"
	echo -e "need to install:\n\n${needed[@]}\n"
	echo "######################################################"
	($pkg_install)
else
	echo "all needed apt packages installed"
fi
