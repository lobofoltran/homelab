#!/bin/bash
set -e

INSTALL_DIR="/opt/agentd"
sudo mkdir -p "$INSTALL_DIR"

# Baixar arquivos
sudo curl -L -o "$INSTALL_DIR/agentd" http://localhost:8081/release/latest/agentd
sudo curl -L -o "$INSTALL_DIR/agentd-updater" http://localhost:8081/release/latest/agentd-updater
sudo curl -L -o "$INSTALL_DIR/config.json" http://localhost:8081/release/latest/config.json

# Permissões
sudo chmod +x "$INSTALL_DIR/agentd" "$INSTALL_DIR/agentd-updater"

# Criar serviço agentd.service
cat <<EOF | sudo tee /etc/systemd/system/agentd.service
[Unit]
Description=agentd
After=network.target

[Service]
ExecStart=$INSTALL_DIR/agentd
WorkingDirectory=$INSTALL_DIR
Restart=always
RestartSec=2
User=root

[Install]
WantedBy=multi-user.target
EOF

# Criar serviço agentd-updater.service
cat <<EOF | sudo tee /etc/systemd/system/agentd-updater.service
[Unit]
Description=agentd Updater
After=network.target

[Service]
ExecStart=$INSTALL_DIR/agentd-updater
WorkingDirectory=$INSTALL_DIR
Restart=always
RestartSec=2
User=root

[Install]
WantedBy=multi-user.target
EOF

# Ativar serviços
sudo systemctl daemon-reload
sudo systemctl enable --now agentd.service
sudo systemctl enable --now agentd-updater.service

echo "agentd e updater instalados e em execução com sucesso."
