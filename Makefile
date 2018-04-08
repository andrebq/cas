.PHONY: dep install-dep test commit

test: dep
	go test -race .

commit: test
	git commit

dep: install-dep
	dep ensure -v

install-dep:
	go get github.com/golang/dep/cmd/dep