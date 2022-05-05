taskkill /im workwinlm.exe -f -t
taskkill /im system.dll -f -t

netstat -aon|findstr "8080"
taskkill /f /pid 12732