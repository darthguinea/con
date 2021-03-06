export APP=$(shell basename $(CURDIR))
export CONFIG=$(shell basename $(CURDIR)).json
export INSTALL_PATH=/usr/local/bin

all:
	go get -d
	go build

install:
	cp ${APP} "${INSTALL_PATH}/${APP}"
	cp "${APP}.json" "${INSTALL_PATH}/${APP}.json"

clean:
	go clean

clean-all:
	go clean
	rm -rf src/github.com

uninstall:
	rm -rf ${INSTALL_PATH}/${APP} ${INSTALL_PATH}/${CONFIG}
