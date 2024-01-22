# install curl, git, ...
apt-get update
apt-get install -y curl git jq

sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable

sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install golang

# vscode-go dependencies 
echo "Getting dependencies for the vscode-go plugin "
# via: https://github.com/microsoft/vscode-go/blob/master/.travis.yml
go get -u -v github.com/acroca/go-symbols
go get -u -v github.com/cweill/gotests/...
go get -u -v github.com/davidrjenni/reftools/cmd/fillstruct
go get -u -v github.com/haya14busa/goplay/cmd/goplay
go get -u -v github.com/mdempsky/gocode
go get -u -v github.com/ramya-rao-a/go-outline
go get -u -v github.com/rogpeppe/godef
go get -u -v github.com/sqs/goreturns
go get -u -v github.com/uudashr/gopkgs/cmd/gopkgs
go get -u -v github.com/zmb3/gogetdoc
go get -u -v golang.org/x/lint/golint
go get -u -v golang.org/x/tools/cmd/gorename
