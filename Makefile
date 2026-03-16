.PHONY: build install clean

BIN_DIR := bin
BIN := $(BIN_DIR)/docs7
LOCAL_BIN := /usr/local/bin/docs7

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) .

install: build
	@if echo "$$PATH" | tr ':' '\n' | grep -qx "$$(cd $(BIN_DIR) && pwd)"; then \
		echo "$(BIN_DIR)/ is in PATH, no symlink needed"; \
	else \
		ln -sf "$$(cd $(BIN_DIR) && pwd)/docs7" $(LOCAL_BIN); \
		echo "linked → $(LOCAL_BIN)"; \
	fi

clean:
	rm -rf $(BIN_DIR)
