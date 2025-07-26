# GOOS=linux GOARCH=amd64 go build -o dist/agentd.exe ./cmd
# mv dist/agentd.exe mock/fake-agentd.exe
# GOOS=linux GOARCH=amd64 go build -o dist/agentd.exe ./cmd
# mv dist/agentd.exe /opt/agentd/agentd.exe
# sudo systemctl restart agentd

GOOS=windows GOARCH=amd64 go build -o release-server/release/latest/windows/agentdBeta.exe ./cmd

cd updater/
GOOS=windows GOARCH=amd64 go build -o ../release-server/release/latest/windows/agentdUpdater.exe ./