#!/bin/bash

if [ "$(grep -Ei 'debian|buntu|mint' /etc/*release)" ]; then
	echo "running apt-get update && apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev"
	sudo apt-get update && \
	sudo apt-get install libsdl2{,-mixer,-image,-ttf,-gfx}-dev
elif [ "$(grep -i 'arch' /etc/issue)" ]; then
	echo "running pacman -Sy && pacman -S sdl2{,_mixer,_image,_ttf,_gfx}"
	sudo pacman -Sy && \
	sudo pacman -S sdl2{,_mixer,_image,_ttf,_gfx}
fi
