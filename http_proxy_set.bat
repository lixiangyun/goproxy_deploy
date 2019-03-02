@echo off 
echo 开始设置IE代理上网 
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 1 /f 
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer /d "http=127.0.0.1:8080;https=127.0.0.1:8080" /f
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyOverride /t REG_SZ /d "<-loopback>" /f
echo 代理设置完成按任意键关闭 
pause>nul