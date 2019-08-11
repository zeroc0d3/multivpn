# -----------------------------------------------------------------------------
#  MAKEFILE RUNNING COMMAND
# -----------------------------------------------------------------------------
#  Author     : Dwi Fahni Denni (@zeroc0d3)
#  Repository : https://github.com/zeroc0d3/multivpn.git
#  License    : Apache License, version 2
# -----------------------------------------------------------------------------
# Notes:
# use [TAB] instead [SPACE]

VERSION=1.0.1
PATH_FOLDER=`pwd`
DEFAULT_KEYS_VPN='/opt/multivpn/keys/default.ovpn'

#---------------------------------
# Cleanup all binary 
#---------------------------------
clean:
	@rm -rf ./bin
	@rm -rf ./build

#---------------------------------
# Installation 
#---------------------------------
install:
	@sudo apt install -y openvpn network-manager-openvpn network-manager-openvpn-gnome
	@make setup

setup:
	@sudo mkdir -p /var/log/multivpn
	@sudo touch /var/log/multivpn/multivpn.log
	@sudo chmod -R 777 /var/log/multivpn
	@sudo mkdir -p /opt/multivpn/config
	@sudo mkdir -p /opt/multivpn/keys
	@sudo chmod -R 777 /opt/multivpn

develop:
	@make setup
	@cp ./keys/auth.txt /opt/multivpn/keys/auth.txt
	@cp ./src/config/*.yaml /opt/multivpn/config
	@sudo chmod -R 777 /opt/multivpn

#---------------------------------
# Running with 'default.ovpn' key
#---------------------------------
run:
	@.${PATH_FOLDER}/bin/multivpn ${DEFAULT_KEYS_VPN}

#---------------------------------
# Build binary for 'multivpn'
#---------------------------------
build:
	@$(GOPATH)/bin/goxc \
	  -bc="darwin,amd64" \
	  -pv=$(VERSION) \
	  -d=build \
	  -build-ldflags "-X main.VERSION=$(VERSION)"
	@go build -o ./bin/multivpn ${PATH_FOLDER}/main.go
	@sudo chmod +x ${PATH_FOLDER}/bin/multivpn

version:
	@echo $(VERSION)