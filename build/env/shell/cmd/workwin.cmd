taskkill /im workwinlm.exe -f -t
taskkill /im system.dll -f -t

# 端口占用
netstat -aon|findstr "8080"
taskkill /f /pid 12732