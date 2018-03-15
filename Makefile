.PHONY: all test start clean 

OUT_PATH=./dist/bin
LOG_PATH=./dist/log
BINARY_NAME=tsdb-mock
OBJECT=$(OUT_PATH)/$(BINARY_NAME)
LOG_FILE=$(LOG_PATH)/stdout.log

DEPENDS=dist

all: $(OBJECT)

$(OBJECT): $(DEPENDS)
	go build -o $(OBJECT) -v

dist:
	mkdir -p $(OUT_PATH) $(LOG_PATH)

clean:
	go clean
	rm -rf ./dist

