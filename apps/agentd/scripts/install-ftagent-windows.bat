@echo off

sc stop agentdUpdater
sc stop agentdBeta

sc delete agentdUpdater
sc delete agentdBeta

if not exist "C:\ProgramData\agentd\" mkdir "C:\ProgramData\agentd\"

REM Baixar arquivos
curl -L -o "C:\ProgramData\agentd\config.json" http://localhost:8081/release/latest/config.json
curl -L -o "C:\ProgramData\agentd\agentdUpdater.exe" http://localhost:8081/release/latest/windows/agentdUpdater.exe
curl -L -o "C:\ProgramData\agentd\agentdBeta.exe" http://localhost:8081/release/latest/windows/agentdBeta.exe

C:\ProgramData\agentd\agentdUpdater.exe install
C:\ProgramData\agentd\agentdBeta.exe install

C:\ProgramData\agentd\agentdUpdater.exe start
C:\ProgramData\agentd\agentdBeta.exe start
