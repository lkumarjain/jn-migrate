#Environment variables
BINARY_NAME = jn-migrate

ifeq ($(OS), Windows_NT)
TARGET = $(BINARY_NAME).exe
else
TARGET = $(BINARY_NAME)
endif

TARGETPATH=$(GOPATH)/bin/$(TARGET)


#Targets
all: dependencies linux windows darwin

build: dependencies test linux windows darwin

linux:
	echo 'Building linux binary'

windows: 
	echo 'Building windows binary'

darwin:
	echo 'Building darwin binary'

dependencies: 
	glide install

quality:
	gometalinter --vendor --exclude=mocks.go --exclude=_test.go --skip=testApps --disable=gotype --disable=dupl --enable=unused --enable=misspell --enable=unparam --deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

test:
	go test -coverprofile ./store/csv/cover.out -covermode=count ./store/csv

mockgen: store-mockgen

store-mockgen:
	mkdir -p store/mock
	mockgen -package mock github.com/lkumarjain/jn-migrate/store Reader > store/mock/mocks.go
