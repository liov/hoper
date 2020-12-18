```bash
sudo apt install vsftpd

sudo passwd jyb

mkdir /home/jyb/ftp

chmod 777 -R /home/jyb/ftp

sudo vim /etc/vsftpd.conf

connect_from_port_21=YES

local_root=/home/jyb/ftp

allow_writeable_chroot=YES

将#chroot_local_user=YES前的注释去掉

pam_service_name=ftp原配置中为vsftpd，ubuntu用户需要更改成ftp

sudo service vsftpd start

sudo service vsftpd restart

```