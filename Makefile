# -----------------------------------------------------------------------------
#  MAKEFILE RUNNING COMMAND
# -----------------------------------------------------------------------------
#  Author     : Dwi Fahni Denni (@zeroc0d3)
#  Repository : https://github.com/zeroc0d3/multivpn.git
#  License    : Apache License, version 2
# -----------------------------------------------------------------------------
# Notes:
# use [TAB] instead [SPACE]

PATH_FOLDER=`pwd`

#---------------------------------
# Installation 
#---------------------------------
vpn-install:
	@sudo apt install -y openvpn network-manager-openvpn network-manager-openvpn-gnome
	@make vpn-setup

vpn-setup:
	@sudo mkdir -p /var/log/multivpn
	@sudo touch /var/log/multivpn/multivpn.log
	@sudo chmod -R 777 /var/log/multivpn

#---------------------------------
# Running with 'default.ovpn' key
#---------------------------------
vpn-run:
	@.${PATH_FOLDER}/bin/multivpn default

#---------------------------------
# Build binary for 'multivpn'
#---------------------------------
vpn-build:
	@go build -o ./bin/multivpn ${PATH_FOLDER}/app.go
	@sudo chmod +x ${PATH_FOLDER}/bin/multivpn