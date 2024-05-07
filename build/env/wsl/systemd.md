# wsl systemd
-- pre
```bash
[boot]
command="/usr/libexec/wsl-systemd"
```
-- now
wsl.conf 的配置设置
wsl.conf 文件基于每个分发配置设置。 (有关 WSL 2 分发版的全局配置，请参阅 .wslconfig) 。

wsl.conf 文件支持四个部分：automount、network和interopuser。 (在.ini文件约定之后建模，密钥将在节下声明，如 .gitconfig files.) 有关存储 wsl.conf 文件的位置的信息，请参阅 wsl.conf 。

systemd 支持
默认情况下，许多 Linux 分发版运行“systemd” (，包括 Ubuntu) 和 WSL 最近添加了对此系统/服务管理器的支持，以便 WSL 更类似于在裸机计算机上使用你喜欢的 Linux 分发版。 需要版本 0.67.6+ 的 WSL 才能启用系统化。 使用命令 wsl --version检查 WSL 版本。 如果需要更新，可以在 Microsoft Store 中获取最新版本的 WSL。 在 博客公告中了解详细信息。

若要启用 systemd，请使用sudo管理员权限在文本编辑器中打开文件wsl.conf，并将以下行添加到/etc/wsl.conf：

```bash
[boot]
systemd=true
```
然后，需要使用 PowerShell 关闭 WSL 分发 wsl.exe --shutdown 版来重启 WSL 实例。 分发重启后，系统应运行。 可以使用以下命令进行确认： systemctl list-unit-files --type=service这将显示服务的状态。