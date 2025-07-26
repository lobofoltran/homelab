BINARY_NAME=homelab-gustavo
INSTALL_DIR=/usr/local/bin
COMPLETION_DIR=/etc/bash_completion.d

build:
	go build -o $(BINARY_NAME) main.go

install: build
	sudo mv $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Binário instalado em $(INSTALL_DIR)/$(BINARY_NAME)"

uninstall:
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Binário removido de $(INSTALL_DIR)"

completion:
	sudo $(INSTALL_DIR)/$(BINARY_NAME) completion bash | sudo tee $(COMPLETION_DIR)/$(BINARY_NAME) > /dev/null
	@echo "Autocompletion instalado em $(COMPLETION_DIR)/$(BINARY_NAME)"

deploy: install completion
	@echo "CLI pronta para uso global."
