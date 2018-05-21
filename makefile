#Environment variables
BINARY_NAME = jn-migrate

ifeq ($(OS), Windows_NT)
TARGET = $(BINARY_NAME).exe
else
TARGET = $(BINARY_NAME)
endif

TARGETPATH=$(GOPATH)/bin/$(TARGET)

#Targets
default: dependencies test linux windows darwin quality

build: dependencies test linux windows darwin

linux:
	go build -o $(TARGETPATH)

windows: 
	GOOS=windows GOARCH=386 go build -o $(TARGETPATH).exe

darwin:
	echo 'Building darwin binary'
#	GOOS=darwin GOARCH=arm64 go build -o $(TARGETPATH)

dependencies: 
	glide install

quality:
	gometalinter --vendor --exclude=mocks.go --exclude=_test.go --skip=testApps --disable=gotype --disable=dupl --enable=unused --enable=misspell --enable=unparam --deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

test:
	go test -coverprofile ./store/cover.out -covermode=count ./store
	go test -coverprofile ./store/csv/cover.out -covermode=count ./store/csv
	go test -coverprofile ./store/sql/cover.out -covermode=count ./store/sql

#Mocks
mockgen: store-mockgen

store-mockgen:
	mkdir -p store/mock
	mockgen -package mock github.com/lkumarjain/jn-migrate/store Reader > store/mock/mocks.go


#Prepration-Scripts
install: 
	go get -v github.com/Masterminds/glide
	cd $GOPATH/src/github.com/Masterminds/glide && git checkout 3e13fd16ed5b0808ba0fb2e4bd98eb325ccde0a1 && go install && cd -

before-script:
	go get github.com/alecthomas/gometalinter
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/fzipp/gocyclo
	go get github.com/tsenart/deadcode
	go get github.com/opennota/check/cmd/aligncheck
	go get github.com/opennota/check/cmd/structcheck
	go get github.com/opennota/check/cmd/varcheck
	go get github.com/mdempsky/maligned
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mibk/dupl
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/gordonklaus/ineffassign
	go get -u mvdan.cc/interfacer
	go get github.com/mdempsky/unconvert
	go get github.com/jgautheron/goconst/cmd/goconst
	go get github.com/GoASTScanner/gas/cmd/gas/...
	go get -u github.com/client9/misspell/cmd/misspell
	go get honnef.co/go/tools/cmd/unused