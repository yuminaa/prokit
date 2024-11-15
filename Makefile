.PHONY: build install uninstall clean

BINARY_NAME := prokit

ifeq ($(OS),Windows_NT)
    DETECTED_OS := windows
    BINARY_NAME := $(BINARY_NAME).exe
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        DETECTED_ARCH := amd64
    else
        DETECTED_ARCH := 386
    endif

    ifeq ($(HOME),)
        HOME := $(USERPROFILE)
    endif
    INSTALL_PATH := $(HOME)/AppData/Local/prokit/bin
    CONFIG_PATH := $(HOME)/AppData/Local/prokit/config
    MKDIR := mkdir
    CP := copy
    RM := del /Q
    RMDIR := rmdir /S /Q
    SEP := \\
else
    DETECTED_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
    DETECTED_ARCH := $(shell uname -m)
    ifeq ($(DETECTED_ARCH),x86_64)
        DETECTED_ARCH := amd64
    endif

    ifeq ($(DETECTED_OS),darwin)
        INSTALL_PATH := $(HOME)/Library/Application Support/prokit/bin
        CONFIG_PATH := $(HOME)/Library/Application Support/prokit/config
    else
        INSTALL_PATH := $(HOME)/.local/bin
        CONFIG_PATH := $(HOME)/.local/share/prokit/config
    endif
    MKDIR := mkdir -p
    CP := cp
    RM := rm -f
    RMDIR := rm -rf
    SEP := /
endif

build:
	@echo "Building $(BINARY_NAME) for $(DETECTED_OS)/$(DETECTED_ARCH)..."
	@go build -o bin/$(BINARY_NAME) ./cmd

install: build
	@echo "Installing $(BINARY_NAME)..."
	@$(MKDIR) "$(INSTALL_PATH)"
	@$(MKDIR) "$(CONFIG_PATH)"
	@$(CP) bin/$(BINARY_NAME) "$(INSTALL_PATH)/$(BINARY_NAME)"
	@$(CP) internal/config/*.json "$(CONFIG_PATH)/"
ifeq ($(OS),Windows_NT)
	@echo "Adding $(INSTALL_PATH) to User PATH if not exists..."
	@powershell -Command "[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';$(INSTALL_PATH)', 'User')"
	@echo "Installation complete! Please restart your terminal and try running: prokit -help"
else
	@echo "Adding $(INSTALL_PATH) to PATH if not exists..."
	@grep -q "export PATH=\"$(INSTALL_PATH):\$$PATH\"" ~/.bashrc || echo "export PATH=\"$(INSTALL_PATH):\$$PATH\"" >> ~/.bashrc
	@grep -q "export PATH=\"$(INSTALL_PATH):\$$PATH\"" ~/.zshrc 2>/dev/null || echo "export PATH=\"$(INSTALL_PATH):\$$PATH\"" >> ~/.zshrc 2>/dev/null || true
	@echo "Installation complete! Please restart your terminal or run: source ~/.bashrc"
endif

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@$(RM) "$(INSTALL_PATH)/$(BINARY_NAME)"
	@$(RMDIR) "$(CONFIG_PATH)"
ifeq ($(OS),Windows_NT)
	@echo "Please manually remove $(INSTALL_PATH) from your PATH if desired"
else
	@echo "Please manually remove $(INSTALL_PATH) from your PATH in ~/.bashrc and ~/.zshrc if desired"
endif
	@echo "Uninstallation complete!"

clean:
	@echo "Cleaning..."
ifeq ($(OS),Windows_NT)
	@if exist bin $(RMDIR) bin
else
	@$(RMDIR) bin 2>/dev/null || true
endif
	@go clean