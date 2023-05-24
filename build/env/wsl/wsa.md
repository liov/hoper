git clone https://github.com/LSPosed/MagiskOnWSALocal

sudo scripts/run.sh 缺依赖装依赖
sudo scripts/build.sh 下载文件可能会失败，多试几次
`output` folder Right-click `Install.ps1` and select Run with PowerShell

安装https://github.com/makazeu/WsaToolbox


安装 https://github.com/NVISOsecurity/MagiskTrustUserCerts
安装 https://github.com/Fuzion24/JustTrustMe
安装 https://github.com/LSPosed/LSPosed

安装 证书

adb connect 127.0.0.1:58526

代理 adb shell "settings put global http_proxy `ip route list match 0 table all scope global | cut -F3`:8888"
删除代理 adb shell settings put global http_proxy :0