# win10教育版
slmgr /ipk NW6C2-QMPVW-D7KKK-3GKT6-VCFB2

slmgr /skms kms.03k.org

slmgr /ato

# windows修改盘符
盘符修改是指更改电脑中分区或设备的驱动器号1。修改盘符的一般步骤是234：
按下组合键“Win+R”，输入“diskmgmt.msc”后按“回车”，打开磁盘管理器。
找到要修改盘符的分区或设备，右键点击，选择“更改驱动器号和路径”。
在弹出的窗口中，点击“更改”，从下拉列表中选择一个未被占用的盘符，点击“确定”。

# windows 端口占用
netstat -aon|findstr "8080"
taskkill /f /pid 12732

# windows 强制关应用
taskkill /im workwinlm.exe -f -t
taskkill /im system.dll -f -t
