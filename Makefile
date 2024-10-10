MAGE_BIN := $(shell go env GOPATH)/bin/mage
CLI_NAME := create-project
CONFIG_DIR := global-config
BINARY_DIR := dist

# Detect OS (Windows or Unix-like systems)
ifeq ($(OS),Windows_NT)
    DETECT_OS := windows
else
    DETECT_OS := unix
endif

# Ensure Mage is installed locally if not present
mage:
ifeq ($(DETECT_OS),windows)
	if not exist $(MAGE_BIN).exe go install github.com/magefile/mage@latest
else
	if [ ! -f $(MAGE_BIN) ]; then go install github.com/magefile/mage@latest; fi
endif

# Build the project for distribution
build: mage
ifeq ($(DETECT_OS),windows)
	if not exist "$(BINARY_DIR)\$(CLI_NAME)" mkdir "$(BINARY_DIR)\$(CLI_NAME)"
	set GOOS=windows
	set GOARCH=amd64
	go build -o "$(BINARY_DIR)\$(CLI_NAME)\$(CLI_NAME).exe"
	xcopy /E /I "$(CONFIG_DIR)" "$(BINARY_DIR)\$(CLI_NAME)\$(CONFIG_DIR)"
else
	mkdir -p $(BINARY_DIR)/$(CLI_NAME)
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$(CLI_NAME)/$(CLI_NAME)
	cp -r $(CONFIG_DIR) $(BINARY_DIR)/$(CLI_NAME)/
endif

# Clean up the build artifacts
clean:
	rm -rf $(BINARY_DIR)
