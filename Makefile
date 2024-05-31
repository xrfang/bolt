BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git log -n1 --pretty=format:%h)
REVS=$(shell git log --oneline|wc -l)
native: export GOARCH=amd64
native: export GOOS=linux
native: release
amd64: export GOARCH=amd64
amd64: export GOOS=linux
amd64: release
arm: export GOOS=linux
arm: export GOARCH=arm
arm: export GOARM=6
arm: setver comprel
win: export GOOS=windows
win: export GOARCH=amd64
win: release
debug: setver compdbg
release: setver comprel
setver:
	cp verinfo.tpl version.go
	sed -i 's/{_BRANCH}/$(BRANCH)/' version.go
	sed -i 's/{_G_HASH}/$(HASH)/' version.go
	sed -i 's/{_G_REVS}/$(REVS)/' version.go
	sed -i 's/{_G_OS}/$(GOOS)/' version.go
	sed -i 's/{_G_HW}/$(GOARCH)/' version.go
comprel:
	CGO_ENABLED=0 go build -ldflags="-s -w" .
compdbg:
	go build -race -gcflags=all=-d=checkptr=0 .
clean:
	rm -fr bolt* version.go
