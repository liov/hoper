https://www.runoob.com/linux/linux-command-manual.html
[toc]
# Linux命令大全完整版
1、文件管理
cat	chattr	chgrp	chmod chown	cksum	cmp	diff    diffstat	file	find	git
gitview	indent	cut	ln less	locate	lsattr	mattrib mc	mdel	mdir	mktemp
more	mmove	mread	mren mtools	mtoolstest	mv	od paste	patch	rcp	rm
slocate	split	tee	tmpwatch touch	umask	which	cp whereis	mcopy	mshowfat	rhmask
scp	awk	read	updatedb
2、文档编辑
col	colrm	comm	csplit ed	egrep	ex	fgrep fmt	fold	grep	ispell
jed	joe	join	look mtype	pico	rgrep	sed sort	spell	tr	expr
uniq	wc	let	 
3、文件传输
lprm	lpr	lpq	lpd bye	ftp	uuto	uupick uucp	uucico	tftp	ncftp
ftpshut	ftpwho	ftpcount	 
4、磁盘管理
cd	df	dirs	du edquota	eject	mcd mdeltree    mdu mkdir	mlabel	mmd
mrd	mzip	pwd	quota mount	mmount	rmdir	rmt stat	tree	umount	ls
quotacheck	quotaoff	lndir	repquota    quotaon	 	 	 
5、磁盘维护
badblocks	cfdisk	dd	e2fsck  ext2ed  fsck	fsck.minix	fsconf
fdformat	hdparm	mformat	mkbootdisk  mkdosfs	mke2fs	mkfs.ext2	mkfs.msdos
mkinitrd	mkisofs	mkswap	mpartition  swapon	symlinks	sync	mbadblocks
mkfs.minix	fsck.ext2	fdisk	losetup mkfs	sfdisk	swapoff	 
6、网络通讯
apachectl	arpwatch	dip	getty   mingetty	uux	telnet	uulog
uustat	ppp-off	netconfig	nc  httpd	ifconfig	minicom	mesg
dnsconf	wall	netstat	ping    pppstats	samba	setserial	talk
traceroute	tty	newaliases	uuname  netconf	write	statserial	efax
pppsetup	tcpdump	ytalk	cu  smbd	testparm	smbclient	shapecfg
7、系统管理
adduser	chfn	useradd	date    exit	finger	fwhios	sleep
suspend	groupdel	groupmod	halt    kill	last	lastb	login
logname	logout	ps	nice    procinfo	top	pstree	reboot
rlogin	rsh	sliplogin	screen  shutdown	rwho	sudo	gitps
swatch	tload	logrotate	uname   chsh	userconf	userdel	usermod
vlock	who	whoami	whois   newgrp	renice	su	skill
w	id	groupadd	free
8、系统设置
reset	clear	alias	dircolors   aumix	bind	chroot	clock
crontab	declare	depmod	dmesg   enable	eval	export	pwunconv
grpconv	rpm	insmod	kbdconfig   lilo	liloconfig	lsmod	minfo
set	modprobe	ntsysv	mouseconfig passwd	pwconv	rdate	resize
rmmod	grpunconv	modinfo	time    setup	sndconfig	setenv	setconsole
timeconfig	ulimit	unset	chkconfig   apmd	hwclock	mkkickstart	fbset
unalias	SVGATextMode	gpasswd	 
9、备份压缩
ar	bunzip2	bzip2	bzip2recover    gunzip	unarj	compress	cpio
dump	uuencode	gzexe	gzip    lha	restore	tar	uudecode
unzip	zip	zipinfo	 
10、设备管理
setleds	loadkeys	rdev	dumpkeys    MAKEDEV	 	 	 
其他命令
Linux bc 命令 Linux tail 命令   Linux xargs 命令

 
## linux系统管理命令
### adduser：新增用户帐号。
语　　法：adduser
补充说明：在Slackware中，adduser指令是个script程序，利用交谈的方式取得输入的用户帐号资料，然后再交由真正建立帐号的useradd指令建立新用户，如此可方便管理员建立用户帐号。在Red Hat Linux中，adduser指令则是useradd指令的符号连接，两者实际上是同一个指令。
### chfn(change finger information)：改变finger指令显示的信息。
语　　法：chfn [-f <真实姓名>][-h <家中电话>][-o <办公地址>][-p <办公电话>][-uv][帐号名称]
补充说明：chfn指令可用来更改执行finger指令时所显示的信息，这些信息都存放在/etc目录里的asswd文件里。若不指定任何参数，则chfn指令会进入问答式界面。
参　　数：
　-f<真实姓名>或--full-name<真实姓名> 　设置真实姓名。 
　-h<家中电话>或--home-phone<家中电话> 　设置家中的电话号码。 
　-o<办公地址>或--office<办公地址> 　设置办公室的地址。 
　-p<办公电话>或--office-phone<办公电话> 　设置办公室的电话号码。 
　-u或--help 　在线帮助。 
　-v或-version 　显示版本信息。
### chsh(change shell)：更换登入系统时使用的shell。
语　　法：chsh [-luv][-s <shell 名称>][用户名称]
补充说明：每位用户在登入系统时，都会拥有预设的shell环境，这个指令可更改其预设值。若不指定任何参数与用户名称，则chsh会以应答的方式进行设置。
参　　数：
　-s<shell 名称>或--shell<shell 名称> 　更改系统预设的shell环境。 
　-l或--list-shells 　列出目前系统可用的shell清单。 
　-u或--help 　在线帮助。 
　-v或-version 　显示版本信息。
### date：显示或设置系统时间与日期。
语    法：
date [-d <字符串>][-u][+%H%I%K%l%M%P%r%s%S%T%X%Z%a%A%b%B%c%d%D%j%m%U%w%x%y%Y%n%t] 或date [-s <字符串>][-u][MMDDhhmmCCYYss] 或 date [--help][--version]
补充说明：第一种语法可用来显示系统日期或时间，以%为开头的参数为格式参数，可指定日期或时间的显示格式。第二种语法可用来设置系统日期与时间。只有管理员才有设置日期与时间的权限。若不加任何参数，data会显示目前的日期与时间。
参　　数：
　%H 　小时(以00-23来表示)。 
　%I 　小时(以01-12来表示)。 
　%K 　小时(以0-23来表示)。 
　%l 　小时(以0-12来表示)。 
　%M 　分钟(以00-59来表示)。 
　%P 　AM或PM。 
　%r 　时间(含时分秒，小时以12小时AM/PM来表示)。 
　%s 　总秒数。起算时间为1970-01-01 00:00:00 UTC。 
　%S 　秒(以本地的惯用法来表示)。 
　%T 　时间(含时分秒，小时以24小时制来表示)。 
　%X 　时间(以本地的惯用法来表示)。 
　%Z 　市区。 
　%a 　星期的缩写。 
　%A 　星期的完整名称。 
　%b 　月份英文名的缩写。 
　%B 　月份的完整英文名称。 
　%c 　日期与时间。只输入date指令也会显示同样的结果。 
　%d 　日期(以01-31来表示)。 
　%D 　日期(含年月日)。 
　%j 　该年中的第几天。 
　%m 　月份(以01-12来表示)。 
　%U 　该年中的周数。 
　%w 　该周的天数，0代表周日，1代表周一，异词类推。 
　%x 　日期(以本地的惯用法来表示)。 
　%y 　年份(以00-99来表示)。 
　%Y 　年份(以四位数来表示)。 
　%n 　在显示时，插入新的一行。 
　%t 　在显示时，插入tab。 
　MM 　月份(必要)。 
　DD 　日期(必要)。 
　hh 　小时(必要)。 
　mm 　分钟(必要)。 
　CC 　年份的前两位数(选择性)。 
　YY 　年份的后两位数(选择性)。 
　ss 　秒(选择性)。 
　-d<字符串> 　显示字符串所指的日期与时间。字符串前后必须加上双引号。 
　-s<字符串> 　根据字符串来设置日期与时间。字符串前后必须加上双引号。 
　-u 　显示GMT。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### exit：退出目前的shell。
语　　法：exit [状态值]
补充说明：执行exit可使shell以指定的状态值退出。若不设置状态值参数，则shell以预设值退出。状态值0代表执行成功，其他值代表执行失败。exit也可用在script，离开正在执行的script，回到shell。
### finger：查找并显示用户信息。
语　　法：finger [-lmsp][帐号名称...]
补充说明：finger指令会去查找，并显示指定帐号的用户相关信息，包括本地与远端主机的用户皆可，帐号名称没有大小写的差别。单独执行finger指令，它会显示本地主机现在所有的用户的登陆信息，包括帐号名称，真实姓名，登入终端机，闲置时间，登入时间以及地址和电话。
参　　数：
  -l 　列出该用户的帐号名称，真实姓名，用户专属目录，登入所用的Shell，登入时间，转信地址，电子邮件状态，还有计划文件和方案文件内容。 
  -m 　排除查找用户的真实姓名。
  -s 　列出该用户的帐号名称，真实姓名，登入终端机，闲置时间，登入时间以及地址和电话。 
  -p 　列出该用户的帐号名称，真实姓名，用户专属目录，登入所用的Shell，登入时间，转信地址，电子邮件状态，但不显示该用户的计划文件和方案文件内容。
### free：显示内存状态。
语　　法： free [-bkmotV][-s <间隔秒数>]
补充说明：free指令会显示内存的使用情况，包括实体内存，虚拟的交换文件内存，共享内存区段，以及系统核心使用的缓冲区等。
参　　数：
　-b 　以Byte为单位显示内存使用情况。 
　-k 　以KB为单位显示内存使用情况。 
　-m 　以MB为单位显示内存使用情况。 
　-o 　不显示缓冲区调节列。 
　-s<间隔秒数> 　持续观察内存使用状况。 
　-t 　显示内存总和列。 
　-V 　显示版本信息。
### fwhois：查找并显示用户信息。
语　　法：fwhios [帐号名称]
补充说明：本指令的功能有点类似finger指令，它会去查找并显示指定帐号的用户相关信息。不同之处在于fwhois指令是到Network Solutions的WHOIS数据库去查找，该帐号名称必须有在上面注册才能寻获，且名称没有大小写的差别
### gitps(gnu interactive tools process status)：报告程序状况。
语　　法：gitps [acefgjlnrsSTuvwxX][p <程序识别码>][t <终端机编号>][U <帐号名称>]
补充说明：gitps是用来报告并管理程序执行的指令，基本上它就是通过ps指令来报告，管理程序，也能通过gitps指令随时中断，删除不必要的程序。因为gitps指令会去执行ps指令，所以其参数和ps指令相当类似。
参　　数：
　a 　显示　现行终端机下的所有程序，包括其他用户的程序。 
　c 　列出程序时，显示每个程序真正的指令名称，而不包含路径，参数或是常驻服务的标示． 
　e 　列出程序时，显示每个程序所使用的环境变量。 
　f 　用ASCII字符显示树状结构，表达程序间的相互关系。 
　g 　显示现行终端机下的所有程序，包括群组领导者的程序。 
　j 　采用工作控制的格式来显示程序状况。 
　l 　采用纤细的格式来显示程序状况。 
　n 　以数字来表示USER和WCHAN栏位。 
　p<程序识别码> 　指定程序识别码，并列出该程序的状况。 
　r 　只列出现行终端机正在执行中的程序。 
　s 　采用程序信号的格式显示程序状况。 
　S 　列出程序时，包括已中断的子程序信息。 
　t<终端机机标号> 　指定终端机编号，并列出属于该终端机的程序的状况。 
　T 　显示现行终端机下的所有程序。 
　u 　以用户为主的格式来显示程序状况。 
　U<帐号名称> 　列出属于该用户的程序的状况。 
　v 　采用虚拟内存的格式显示程序状况。 
　w 　采用宽阔的格式来显示程序状况。 
　x 　显示所有程序，不以终端机来区分。 
　X 　采用旧试的Linux i386登陆格式显示程序状况。
### groupdel(group delete)：删除群组。
语　　法：groupdel [群组名称]
补充说明：需要从系统上删除群组时，可用groupdel指令来完成这项工作。倘若该群组中仍包括某些用户，则必须先删除这些用户后，方能删除群组。
### groupmod(group modify)：更改群组识别码或名称。
语　　法：groupmod [-g <群组识别码> <-o>][-n <新群组名称>][群组名称]
补充说明：需要更改群组的识别码或名称时，可用groupmod指令来完成这项工作。
参　　数：
　-g <群组识别码> 　设置欲使用的群组识别码。 
　-o 　重复使用群组识别码。 
　-n <新群组名称> 　设置欲使用的群组名称。
### halt：关闭系统。
语　　法：halt [-dfinpw]
补充说明：halt会先检测系统的runlevel。若runlevel为0或6，则关闭系统，否则即调用shutdown来关闭系统。
参　　数：
  -d 　不要在wtmp中记录。 
  -f 　不论目前的runlevel为何，不调用shutdown即强制关闭系统。 
  -i 　在halt之前，关闭全部的网络界面。 
  -n 　halt前，不用先执行sync。 
  -p 　halt之后，执行poweroff。 
  -w 　仅在wtmp中记录，而不实际结束系统。
### id：显示用户的ID，以及所属群组的ID。
语　　法：id [-gGnru][--help][--version][用户名称]
补充说明：id会显示用户以及所属群组的实际与有效ID。若两个ID相同，则仅显示实际ID。若仅指定用户名称，则显示目前用户的ID。
参　　数：
　-g或--group 　显示用户所属群组的ID。 
　-G或--groups 　显示用户所属附加群组的ID。 
　-n或--name 　显示用户，所属群组或附加群组的名称。 
　-r或--real 　显示实际ID。 
　-u或--user 　显示用户ID。 
　-help 　显示帮助。 
　-version 　显示版本信息。
### kill：删除执行中的程序或工作。
语　　法：kill [-s <信息名称或编号>][程序]　或　kill [-l <信息编号>]
补充说明：kill可将指定的信息送至程序。预设的信息为SIGTERM(15)，可将指定程序终止。若仍无法终止该程序，可使用SIGKILL(9)信息尝试强制删除程序。程序或工作的编号可利用ps指令或jobs指令查看。 
参　　数：
　-l <信息编号> 　若不加<信息编号>选项，则-l参数会列出全部的信息名称。 
　-s <信息名称或编号> 　指定要送出的信息。 
　[程序] 　[程序]可以是程序的PID或是PGID，也可以是工作编号。
### last：列出目前与过去登入系统的用户相关信息。
语　　法：last [-adRx][-f <记录文件>][-n <显示列数>][帐号名称...][终端机编号...]
补充说明：单独执行last指令，它会读取位于/var/log目录下，名称为wtmp的文件，并把该给文件的内容记录的登入系统的用户名单全部显示出来。
参　　数：
　-a 　把从何处登入系统的主机名称或IP地址，显示在最后一行。 
　-d 　将IP地址转换成主机名称。 
　-f <记录文件> 　指定记录文件。 
　-n <显示列数>或-<显示列数> 　设置列出名单的显示列数。 
　-R 　不显示登入系统的主机名称或IP地址。 
　-x 　显示系统关机，重新开机，以及执行等级的改变等信息。
### lastb：列出登入系统失败的用户相关信息。
语　　法：lastb [-adRx][-f <记录文件>][-n <显示列数>][帐号名称...][终端机编号...]
补充说明：单独执行lastb指令，它会读取位于/var/log目录下，名称为btmp的文件，并把该文件内容记录的登入失败的用户名单，全部显示出来。
参　　数：
　-a 　把从何处登入系统的主机名称或IP地址显示在最后一行。 
　-d 　将IP地址转换成主机名称。 
　-f<记录文件> 　指定记录文件。 
　-n<显示列数>或-<显示列数> 　设置列出名单的显示列数。 
　-R 　不显示登入系统的主机名称或IP地址。 
　-x 　显示系统关机，重新开机，以及执行等级的改变等信息。
### login：登入系统。
语　　法：login
补充说明：login指令让用户登入系统，您亦可通过它的功能随时更换登入身份。在Slackware发行版中 ，您可在指令后面附加欲登入的用户名称，它会直接询问密码，等待用户输入。当/etc目录里含名称为nologin的文件时，系统只root帐号登入系统，其他用户一律不准登入。
### logname：显示用户名称。
语　　法：logname [--help][--version]
补充说明：执行logname指令，它会显示目前用户的名称。
参　　数：
　--help 　在线帮助。 
　--vesion 　显示版本信息。
### logout：退出系统。
语　　法：logout
补充说明：logout指令让用户退出系统，其功能和login指令相互对应。
### logrotate：管理记录文件。
语　　法：logrotate [-?dfv][-s <状态文件>][--usage][配置文件]
补充说明：使用logrotate指令，可让你轻松管理系统所产生的记录文件。它提供自动替换，压缩，删除和邮寄记录文件，每个记录文件都可被设置成每日，每周或每月处理，也能在文件太大时立即处理。您必须自行编辑，指定配置文件，预设的配置文件存放在/etc目录下，文件名称为logrotate.conf。
参　　数：
　-?或--help 　在线帮助。 
　-d或--debug 　详细显示指令执行过程，便于排错或了解程序执行的情况。 
　-f或--force 　强行启动记录文件维护操作，纵使logrotate指令认为没有需要亦然。 
　-s<状态文件>或--state=<状态文件> 　使用指定的状态文件。 
　-v或--version 　显示指令执行过程。 
　-usage 　显示指令基本用法。
### newgrp：登入另一个群组。
语　　法：newgrp [群组名称]
补充说明：newgrp指令类似login指令，当它是以相同的帐号，另一个群组名称，再次登入系统。欲使用newgrp指令切换群组，您必须是该群组的用户，否则将无法登入指定的群组。单一用户要同时隶属多个群组，需利用交替用户的设置。若不指定群组名称，则newgrp指令会登入该用户名称的预设群组。
### nice：设置优先权。
语　　法：nice [-n <优先等级>][--help][--version][执行指令]
补充说明：nice指令可以改变程序执行的优先权等级。
参　　数：</b 
　-n<优先等级>或-<优先等级>或--adjustment=<优先等级> 　设置欲执行的指令的优先权等级。等级的范围从-20－19，其中-20最高，19最低，只有系统管理者可以设置负数的等级。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### procinfo(process information)：显鞠低匙刺?BR>
语　　法：procinfo [-abdDfhimsSv][-F <输出文件>][-n <间隔秒数>]
补充说明：procinfo指令从/proc目录里读取相关数据，将数据妥善整理过后输出到标准输出设备。
参　　数：
　-a 　显示所有信息。 
　-b 　显示磁盘设备的区块数目，而非存取数目。 
　-d 　显示系统信息每秒间的变化差额，而非总和的数值。本参数必须配合"-f"参数使用 
　-D 　此参数效果和指定"-d"参数类似，但内存和交换文件的信息为总和数值。 
　-f 　进入全画面的互动式操作界面。 
　-F<输出文件> 　把信息状态输出到文件保存起来，而非预设的标准输出设备。 
　-h 　在线帮助。 
　-i 　显示完整的IRP列表。 
　-m 　显示系统模块和外围设备等相关信息。 
　-n间隔秒数> 　设置全画面互动模式的信息更新速度，单位以秒计算。 
　-s 　显示系统的内存，磁盘空间，IRP和DMA等信息，此为预设值。 
　-S 　搭配参数"-d"或"-D"使用时，每秒都会更新信息，不论是否有使用参数"-n"。 
　-v 　显示版本信息。
### ps(process status)：报告程序状况。
语　　法：ps [-aAcdefHjlmNVwy][acefghLnrsSTuvxX][-C <指令名称>][-g <群组名称>][-G <群组识别码>][-p <程序识别码>][p <程序识别码>][-s <阶段作业>][-t <终端机编号>][t <终端机编号>][-u <用户识别码>][-U <用户识别码>][U <用户名称>][-<程序识别码>][--cols <每列字符数>][--columns <每列字符数>][--cumulative][--deselect][--forest][--headers][--help][--info][--lines <显示列数>][--no-headers][--group <群组名称>][-Group <群组识别码>][--pid <程序识别码>][--rows <显示列数>][--sid <阶段作业>][--tty <终端机编号>][--user <用户名称>][--User <用户识别码>][--version][--width <每列字符数>]
补充说明：ps是用来报告程序执行状况的指令，您可以搭配kill指令随时中断，删除不必要的程序。
参　　数：
　-a 　显示所有终端机下执行的程序，除了阶段作业领导者之外。 
　a 　显示现行终端机下的所有程序，包括其他用户的程序。 
　-A 　显示所有程序。 
　-c 　显示CLS和PRI栏位。 
　c 　列出程序时，显示每个程序真正的指令名称，而不包含路径，参数或常驻服务的标示。 
　-C<指令名称> 　指定执行指令的名称，并列出该指令的程序的状况。 
　-d 　显示所有程序，但不包括阶段作业领导者的程序。 
　-e 　此参数的效果和指定"A"参数相同。 
　e 　列出程序时，显示每个程序所使用的环境变量。 
　-f 　显示UID,PPIP,C与STIME栏位。 
　f 　用ASCII字符显示树状结构，表达程序间的相互关系。 
　-g<群组名称> 　此参数的效果和指定"-G"参数相同，当亦能使用阶段作业领导者的名称来指定。 
　g 　显示现行终端机下的所有程序，包括群组领导者的程序。 
　-G<群组识别码> 　列出属于该群组的程序的状况，也可使用群组名称来指定。 
　h 　不显示标题列。 
　-H 　显示树状结构，表示程序间的相互关系。 
　-j或j 　采用工作控制的格式显示程序状况。 
　-l或l 　采用详细的格式来显示程序状况。 
　L 　列出栏位的相关信息。 
　-m或m 　显示所有的执行绪。 
　n 　以数字来表示USER和WCHAN栏位。 
　-N 　显示所有的程序，除了执行ps指令终端机下的程序之外。 
　-p<程序识别码> 　指定程序识别码，并列出该程序的状况。 
　p<程序识别码> 　此参数的效果和指定"-p"参数相同，只在列表格式方面稍有差异。 
　r 　只列出现行终端机正在执行中的程序。 
　-s<阶段作业> 　指定阶段作业的程序识别码，并列出隶属该阶段作业的程序的状况。 
　s 　采用程序信号的格式显示程序状况。 
　S 　列出程序时，包括已中断的子程序资料。 
　-t<终端机编号> 　指定终端机编号，并列出属于该终端机的程序的状况。 
　t<终端机编号> 　此参数的效果和指定"-t"参数相同，只在列表格式方面稍有差异。 
　-T 　显示现行终端机下的所有程序。 
　-u<用户识别码> 　此参数的效果和指定"-U"参数相同。 
　u 　以用户为主的格式来显示程序状况。 
　-U<用户识别码> 　列出属于该用户的程序的状况，也可使用用户名称来指定。 
　U<用户名称> 　列出属于该用户的程序的状况。 
　v 　采用虚拟内存的格式显示程序状况。 
　-V或V 　显示版本信息。 
　-w或w 　采用宽阔的格式来显示程序状况。　 
　x 　显示所有程序，不以终端机来区分。 
　X 　采用旧式的Linux i386登陆格式显示程序状况。 
　-y 　配合参数"-l"使用时，不显示F(flag)栏位，并以RSS栏位取代ADDR栏位　。 
　-<程序识别码> 　此参数的效果和指定"p"参数相同。 
　--cols<每列字符数> 　设置每列的最大字符数。 
　--columns<每列字符数> 　此参数的效果和指定"--cols"参数相同。 
　--cumulative 　此参数的效果和指定"S"参数相同。 
　--deselect 　此参数的效果和指定"-N"参数相同。 
　--forest 　此参数的效果和指定"f"参数相同。 
　--headers 　重复显示标题列。 
　--help 　在线帮助。 
　--info 　显示排错信息。 
　--lines<显示列数> 　设置显示画面的列数。 
　--no-headers 　此参数的效果和指定"h"参数相同，只在列表格式方面稍有差异。 
　--group<群组名称> 　此参数的效果和指定"-G"参数相同。 
　--Group<群组识别码> 　此参数的效果和指定"-G"参数相同。 
　--pid<程序识别码> 　此参数的效果和指定"-p"参数相同。 
　--rows<显示列数> 　此参数的效果和指定"--lines"参数相同。 
　--sid<阶段作业> 　此参数的效果和指定"-s"参数相同。 
　--tty<终端机编号> 　此参数的效果和指定"-t"参数相同。 
　--user<用户名称> 　此参数的效果和指定"-U"参数相同。 
　--User<用户识别码> 　此参数的效果和指定"-U"参数相同。 
　--version 　此参数的效果和指定"-V"参数相同。 
　--widty<每列字符数> 　此参数的效果和指定"-cols"参数相同。
### pstree(process status tree)：以树状图显示程序。
语　　法：pstree [-acGhlnpuUV][-H <程序识别码>][<程序识别码>/<用户名称>]
补充说明：pstree指令用ASCII字符显示树状结构，清楚地表达程序间的相互关系。如果不指定程序识别码或用户名称，则会把系统启动时的第一个程序视为基层，并显示之后的所有程序。若指定用户名称，便会以隶属该用户的第一个程序当作基层，然后显示该用户的所有程序。
参　　数：
　-a 　显示每个程序的完整指令，包含路径，参数或是常驻服务的标示。 
　-c 　不使用精简标示法。 
　-G 　使用VT100终端机的列绘图字符。 
　-h 　列出树状图时，特别标明现在执行的程序。 
　-H<程序识别码> 　此参数的效果和指定"-h"参数类似，但特别标明指定的程序。 
　-l 　采用长列格式显示树状图。 
　-n 　用程序识别码排序。预设是以程序名称来排序。 
　-p 　显示程序识别码。 
　-u 　显示用户名称。 
　-U 　使用UTF-8列绘图字符。 
　-V 　显示版本信息。
### reboot：重新开机。
语　　法：dreboot [-dfinw]
补充说明：执行reboot指令可让系统停止运作，并重新开机。
参　　数：
  -d 　重新开机时不把数据写入记录文件/var/tmp/wtmp。本参数具有"-n"参数的效果。 
  -f 　强制重新开机，不调用shutdown指令的功能。 
  -i 　在重开机之前，先关闭所有网络界面。 
  -n 　重开机之前不检查是否有未结束的程序。 
  -w 　仅做测试，并不真的将系统重新开机，只会把重开机的数据写入/var/log目录下的wtmp记录文件。
### renice：调整优先权。
语　　法：renice [优先等级][-g <程序群组名称>...][-p <程序识别码>...][-u <用户名称>...]
补充说明：renice指令可重新调整程序执行的优先权等级。预设是以程序识别码指定程序调整其优先权，您亦可以指定程序群组或用户名称调整优先权等级，并修改所有隶属于该程序群组或用户的程序的优先权。等级范围从-20--19，只有系统管理者可以改变其他用户程序的优先权，也仅有系统管理者可以设置负数等级。
参　　数：
　-g <程序群组名称> 　使用程序群组名称，修改所有隶属于该程序群组的程序的优先权。 
　-p <程序识别码> 　改变该程序的优先权等级，此参数为预设值。 
　-u <用户名称> 　指定用户名称，修改所有隶属于该用户的程序的优先权。
### rlogin(remote login)：远端登入。
语　　法：rlogin [-8EL][-e <脱离字符>][-l <用户名称>][主机名称或IP地址]
补充说明：执行rlogin指令开启终端机阶段操作，并登入远端主机。
参　　数：
　-8 　允许输入8位字符数据。 
　-e脱离字符> 　设置脱离字符。 
　-E 　滤除脱离字符。 
　-l用户名称> 　指定要登入远端主机的用户名称。 
　-L 　使用litout模式进行远端登入阶段操作。
### rsh(remote shell)：远端登入的Shell。
语　　法：rsh [-dn][-l <用户名称>][主机名称或IP地址][执行指令]
补充说明：rsh提供用户环境，也就是Shell，以便指令能够在指定的远端主机上执行。
参　　数：
　-d 　使用Socket层级的排错功能。 
　-l<用户名称> 　指定要登入远端主机的用户名称。 
　-n 　把输入的指令号向代号为/dev/null的特殊外围设备。
### rwho：查看系统用户。
语　　法：rwho [-a]
补充说明：rwho指令的效果类似who指令，但它会显示局域网里所有主机的用户。主机必须提供rwhod常驻服务的功能，方可使用rwho指令。
参　　数：
  -a 　列出所有的用户，包括闲置时间超过1个小时以上的用户。
### screen：多重视窗管理程序。
语　　法：screen [-AmRvx -ls -wipe][-d <作业名称>][-h <行数>][-r <作业名称>][-s <shell>][-S <作业名称>]
补充说明：screen为多重视窗管理程序。此处所谓的视窗，是指一个全屏幕的文字模式画面。通常只有在使用telnet登入主机或是使用老式的终端机时，才有可能用到screen程序。
参　　数：
　-A 　将所有的视窗都调整为目前终端机的大小。 
　-d<作业名称> 　将指定的screen作业离线。 
　-h<行数> 　指定视窗的缓冲区行数。 
　-m 　即使目前已在作业中的screen作业，仍强制建立新的screen作业。 
　-r<作业名称> 　恢复离线的screen作业。 
　-R 　先试图恢复离线的作业。若找不到离线的作业，即建立新的screen作业。 
　-s<shell> 　指定建立新视窗时，所要执行的shell。 
　-S<作业名称> 　指定screen作业的名称。 
　-v 　显示版本信息。 
　-x 　恢复之前离线的screen作业。 
　-ls或--list 　显示目前所有的screen作业。 
　-wipe 　检查目前所有的screen作业，并删除已经无法使用的screen作业。
### shutdown：系统关机指令。
语　　法：shutdown [-efFhknr][-t 秒数][时间][警告信息]
补充说明：shutdown指令可以关闭所有程序，并依用户的需要，进行重新开机或关机的动作。
参　　数：
　-c 　当执行"shutdown -h 11:50"指令时，只要按+键就可以中断关机的指令。 
　-f 　重新启动时不执行fsck。 
　-F 　重新启动时执行fsck。 
　-h 　将系统关机。 
　-k 　只是送出信息给所有用户，但不会实际关机。 
　-n 　不调用init程序进行关机，而由shutdown自己进行。 
　-r 　shutdown之后重新启动。 
　-t<秒数> 　送出警告信息和删除信息之间要延迟多少秒。 
　[时间] 　设置多久时间后执行shutdown指令。 
　[警告信息] 　要传送给所有登入用户的信息
### sliplogin：将SLIP接口加入标准输入。
语　　法：sliplogin [用户名称]
补充说明：sliplogin可将SLIP接口加入标准输入，把一般终端机的连线变成SLIP连线。通常可用来建立SLIP服务器，让远端电脑以SLIP连线到服务器。sliplogin活去检查/etc/slip/slip.hosts文件中是否有相同的用户名称。通过检查后，sliplogin会调用执行shell script来设置IP地址，子网掩码等网络界面环境。此shell script通常是/etc/slip/slip.login。
### su(super user)：变更用户身份。
语　　法：su [-flmp][--help][--version][-][-c <指令>][-s <shell>][用户帐号]
补充说明：su可让用户暂时变更登入的身份。变更时须输入所要变更的用户帐号与密码。
参　　数：
　-c<指令>或--command=<指令> 　执行完指定的指令后，即恢复原来的身份。 
　-f或--fast 　适用于csh与tsch，使shell不用去读取启动文件。 
　-.-l或--login 　改变身份时，也同时变更工作目录，以及HOME,SHELL,USER,LOGNAME。此外，也会变更PATH变量。  
　-m,-p或--preserve-environment 　变更身份时，不要变更环境变量。 
　-s<shell>或--shell=<shell> 　指定要执行的shell。 
　--help 　显示帮助。 
　--version 　显示版本信息。 
　[用户帐号] 　指定要变更的用户。若不指定此参数，则预设变更为root。
### sudo：以其他身份来执行指令。
语　　法：sudo [-bhHpV][-s <shell>][-u <用户>][指令] 或 sudo [-klv]
补充说明：sudo可让用户以其他的身份来执行指定的指令，预设的身份为root。在/etc/sudoers中设置了可执行sudo指令的用户。若其未经授权的用户企图使用sudo，则会发出警告的邮件给管理员。用户使用sudo时，必须先输入密码，之后有5分钟的有效期限，超过期限则必须重新输入密码。  
参　　数：
　-b 　在后台执行指令。 
　-h 　显示帮助。 
　-H 　将HOME环境变量设为新身份的HOME环境变量。 
　-k 　结束密码的有效期限，也就是下次再执行sudo时便需要输入密码。 
　-l 　列出目前用户可执行与无法执行的指令。 
　-p 　改变询问密码的提示符号。 
　-s<shell> 　执行指定的shell。 
　-u<用户> 　以指定的用户作为新的身份。若不加上此参数，则预设以root作为新的身份。 
　-v 　延长密码有效期限5分钟。 
　-V 　显示版本信息。
### suspend：暂停执行shell。
语　　法：suspend [-f]
补充说明：suspend为shell内建指令，可暂停目前正在执行的shell。若要恢复，则必须使用SIGCONT信息。
参　　数：
  -f 　若目前执行的shell为登入的shell，则suspend预设无法暂停此shell。若要强迫暂停登入的shell，则必须使用-f参数。
### swatch(simple watcher)：系统监控程序。
语　　法：swatch [-A <分隔字符>][-c <设置文件>][-f <记录文件>][-I <分隔字符>][-P <分隔字符>][-r <时间>][-t <记录文件>]
补充说明：swatch可用来监控系统记录文件，并在发现特定的事件时，执行指定的动作。swatch所监控的事件以及对应事件的动作都存放在swatch的配置文件中。预设的配置文件为拥护根目录下的.swatchrc。然而在Red Hat Linux的预设用户根目录下并没有.swatchrc配置文件，您可将/usr/doc/swatch-2.2/config_files/swatchrc.personal文件复制到用户根目录下的.swatchrc，然后修改.swatchrc所要监控的事件及执行的动作。
参　　数：
　-A<分隔字符> 　预设配置文件中，动作的分隔字符，预设为逗号。 
　-c设置文件> 　指定配置文件，而不使用预设的配置文件。 
　-f记录文件> 　检查指定的记录文件，检查完毕后不会继续监控该记录文件。 
　-I分隔字符> 　指定输入记录的分隔字符，预设为换行字符。 
　-P分隔字符> 　指定配置文件中，事件的分隔字符，预设为逗号。 
　-r时间> 　在指定的时间重新启动。 
　-t<记录文件> 　检查指定的记录文件，并且会监控加入记录文件中的后继记录。
### tload：显示系统负载状况。
语　　法：tload [-V][-d <间隔秒数>][-s <刻度大小>][终端机编号]
补充说明：tload指令使用ASCII字符简单地以文字模式显示系统负载状态。假设不给予终端机编号，则会在执行tload指令的终端机显示负载情形。
参　　数：
　-d<间隔秒数> 　设置tload检测系统负载的间隔时间，单位以秒计算。 
　-s<刻度大小> 　设置图表的垂直刻度大小，单位以列计算。 
　-V 　显示版本信息。
### top：显示，管理执行中的程序。
语　　法：top [bciqsS][d <间隔秒数>][n <执行次数>]
补充说明：执行top指令可显示目前正在系统中执行的程序，并通过它所提供的互动式界面，用热键加以管理。
参　　数：
　b 　使用批处理模式。 
　c 　列出程序时，显示每个程序的完整指令，包括指令名称，路径和参数等相关信息。 
　d<间隔秒数> 　设置top监控程序执行状况的间隔时间，单位以秒计算。 
　i 　执行top指令时，忽略闲置或是已成为Zombie的程序。 
　n<执行次数> 　设置监控信息的更新次数。 
　q 　持续监控程序执行的状况。 
　s 　使用保密模式，消除互动模式下的潜在危机。 
　S 　使用累计模式，其效果类似ps指令的"-S"参数。
### uname：显示系统信息。
语　　法：uname [-amnrsv][--help][--version]
补充说明：uname可显示电脑以及操作系统的相关信息。
参　　数：
　-a或--all 　显示全部的信息。 
　-m或--machine 　显示电脑类型。 
　-n或-nodename 　显示在网络上的主机名称。 
　-r或--release 　显示操作系统的发行编号。 
　-s或--sysname 　显示操作系统名称。 
　-v 　显示操作系统的版本。 
　--help 　显示帮助。 
　--version 　显示版本信息。
### useradd：建立用户帐号。
语　　法：useradd [-mMnr][-c <备注>][-d <登入目录>][-e <有效期限>][-f <缓冲天数>][-g <群组>][-G <群组>][-s <shell>][-u <uid>][用户帐号] 或 useradd -D [-b][-e <有效期限>][-f <缓冲天数>][-g <群组>][-G <群组>][-s <shell>]
补充说明：useradd可用来建立用户帐号。帐号建好之后，再用passwd设定帐号的密码．而可用userdel删除帐号。使用useradd指令所建立的帐号，实际上是保存在/etc/passwd文本文件中。
参　　数：
　-c<备注> 　加上备注文字。备注文字会保存在passwd的备注栏位中。　 
　-d<登入目录> 　指定用户登入时的启始目录。 
　-D 　变更预设值． 
　-e<有效期限> 　指定帐号的有效期限。 
　-f<缓冲天数> 　指定在密码过期后多少天即关闭该帐号。 
　-g<群组> 　指定用户所属的群组。 
　-G<群组> 　指定用户所属的附加群组。 
　-m 　自动建立用户的登入目录。 
　-M 　不要自动建立用户的登入目录。 
　-n 　取消建立以用户名称为名的群组． 
　-r 　建立系统帐号。 
　-s<shell>　 　指定用户登入后所使用的shell。 
　-u<uid> 　指定用户ID。
### userconf：用户帐号设置程序。
语　　法：userconf [--addgroup <群组>][--adduser <用户ID><群组><用户名称><shell>][--delgroup <群组>][--deluser <用户ID>][--help]
补充说明：userconf实际上为linuxconf的符号连接，提供图形界面的操作方式，供管理员建立与管理各类帐号。若不加任何参数，即进入图形界面。 
参　　数：
　--addgroup<群组> 　新增群组。 
　--adduser<用户ID><群组><用户名称><shell> 　新增用户帐号。 
　--delgroup<群组> 　删除群组。 
　--deluser<用户ID> 　删除用户帐号。 
　--help 　显示帮助。
### userdel：删除用户帐号。
语　　法：userdel [-r][用户帐号]
补充说明：userdel可删除用户帐号与相关的文件。若不加参数，则仅删除用户帐号，而不删除相关文件。
参　　数：
  -f 　删除用户登入目录以及目录中所有文件。
### usermod：修改用户帐号。
语　　法：usermod [-LU][-c <备注>][-d <登入目录>][-e <有效期限>][-f <缓冲天数>][-g <群组>][-G <群组>][-l <帐号名称>][-s <shell>][-u <uid>][用户帐号]
补充说明：usermod可用来修改用户帐号的各项设定。
参　　数：
　-c<备注> 　修改用户帐号的备注文字。 
　-d登入目录> 　修改用户登入时的目录。 
　-e<有效期限> 　修改帐号的有效期限。 
　-f<缓冲天数> 　修改在密码过期后多少天即关闭该帐号。 
　-g<群组> 　修改用户所属的群组。 
　-G<群组> 　修改用户所属的附加群组。 
　-l<帐号名称> 　修改用户帐号名称。 
　-L 　锁定用户密码，使密码无效。 
　-s<shell> 　修改用户登入后所使用的shell。 
　-u<uid> 　修改用户ID。 
　-U 　解除密码锁定。
### vlock(virtual console lock)：锁住虚拟终端。
语　　法：vlock [-achv]
补充说明：执行vlock指令可锁住虚拟终端，避免他人使用。
参　　数：
　-a或--all 　锁住所有的终端阶段作业，如果您在全屏幕的终端中使用本参数，则会将用键盘
　切换终端机的功能一并关闭。 
　-c或--current 　锁住目前的终端阶段作业，此为预设值。 
　-h或--help 　在线帮助。 
　-v或--version 　显示版本信息。
### w：显示目前登入系统的用户信息。
语　　法：w [-fhlsuV][用户名称]
补充说明：执行这项指令可得知目前登入系统的用户有那些人，以及他们正在执行的程序。单独执行w
指令会显示所有的用户，您也可指定用户名称，仅显示某位用户的相关信息。
参　　数：
  -f 　开启或关闭显示用户从何处登入系统。 
  -h 　不显示各栏位的标题信息列。 
  -l 　使用详细格式列表，此为预设值。 
  -s 　使用简洁格式列表，不显示用户登入时间，终端机阶段作业和程序所耗费的CPU时间。 
  -u 　忽略执行程序的名称，以及该程序耗费CPU时间的信息。 
  -V 　显示版本信息。
### who：显示目前登入系统的用户信息。
语　　法：who [-Himqsw][--help][--version][am i][记录文件]
补充说明：执行这项指令可得知目前有那些用户登入系统，单独执行who指令会列出登入帐号，使用的 终端机，登入时间以及从何处登入或正在使用哪个X显示器。
参　　数：
　-H或--heading 　显示各栏位的标题信息列。 
　-i或-u或--idle 　显示闲置时间，若该用户在前一分钟之内有进行任何动作，将标示成"."号，如果该用户已超过24小时没有任何动作，则标示出"old"字符串。 
　-m 　此参数的效果和指定"am i"字符串相同。 
　-q或--count 　只显示登入系统的帐号名称和总人数。 
　-s 　此参数将忽略不予处理，仅负责解决who指令其他版本的兼容性问题。 
　-w或-T或--mesg或--message或--writable 　显示用户的信息状态栏。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### whoami：先似乎用户名称。
语　　法：whoami [--help][--version]
补充说明：显示自身的用户名称，本指令相当于执行"id -un"指令。
参　　数：
　--help 　在线帮助。 
　--version 　显示版本信息。
### whois：查找并显示用户信息。
语　　法：whois [帐号名称]
补充说明：whois指令会去查找并显示指定帐号的用户相关信息，因为它是到Network Solutions的WHOIS数据库去查找，所以该帐号名称必须在上面注册方能寻获，且名称没有大小写的差别。

 
## linux系统设置命令
### alias：设置指令的别名。
语　　法：alias[别名]=[指令名称]
补充说明：用户可利用alias，自定指令的别名。若仅输入alias，则可列出目前所有的别名设置。　alias的效力仅及于该次登入的操作。若要每次登入是即自动设好别名，可在.profile或.cshrc中设定指令的别名。
参　　数：若不加任何参数，则列出目前所有的别名设置。
### apmd(advanced power management BIOS daemon)：进阶电源管理服务程序。
语　　法：apmd [-u v V W][-p <百分比变化量>][-w <百分比值>]
补充说明：apmd负责BIOS进阶电源管理(APM)相关的记录，警告与管理工作。
参　　数：
  -p<百分比变化量>或--percentage<百分比变化量>  当电力变化的幅度超出设置的百分比变化量，即记录事件百分比变化量的预设值为5，若设置值超过100，则关闭此功能。 
  -u或--utc 将BIOS时钟设为UTC，以便从悬待模式恢复时，将-u参数传送至clock或hwclock程序。 
  -v或--verbose 记录所有的APM事件。 
  -V或--version 显示版本信息。 
  -w<百分比值>或--warn<百分比值> 当电池不在充电状态时，且电池电量低于设置的百分比值，则在syslog(2)的ALERT层记录警告信息。百分比值的预设置为10，若设置为0，则关闭此功能。  
  -W或--wall 发出警告信息给所有人。
### aumix(audio mixer)：设置音效装置。
语　　法：aumix [-123bcilmoprstvwWx][(+/-)强度][PqR][-dfhILqS]
补充说明：设置各项音效装置的信号强度以及指定播放与录音的装置。
参　　数：
  [-123bcilmoprstvwWx]为频道参数，用来指定装置的频道；[PqR]可用来指定播放或录音装置；[-dfhILqS] 则为指令参数。若不加任何参数，aumix会显示简单的图形界面供调整设置频道参数：
  -1 　输入信号线１。 
  -2 　输入信号线２。 
  -3 　输入信号线３。 
  -b 　低音。 
  -c 　CD。 
  -i 　输入信号强度。 
  -m 　麦克风。 
  -o 　输出信号强度。  
  -p 　PC喇叭。 
  -r 　录音。 
  -s 　合成器。 
  -t 　高音。 
  -v 　主音量。 
  -w 　PCM。 
  -W 　PCM2。 
  -x 　混音器。 
  (+/-)强度 　出现(+/-)时，代表在原有的强度上加减指定值。若未使用(+/-)，则直接将强度设为指定值。 　指定音效装置 
  P 　指定播放装置。 
  q 　显示频道设置。 
  R 　指定录音装置。 
　指令参数 
  -d 　指定音效装置的名称。 
  -f 　指定存储或载入设置的文件。 
  -h 　在使用时显示信息。 
  -I 　以图形界面方式来执行aumix。 
  -L 　从$HOME/.aumixrc或/etc/aumixrc载入设置。 
  -q 　显示所有频道的设置值。 
  -S 　将设置值保存至/HOME/.aumixrc。
### bind：显示或设置键盘按键与其相关的功能。
语　　法：bind [-dlv][-f <按键配置文件>][-m <按键配置>][-q <功能>]
补充说明：您可以利用bind命令了解有哪些按键组合与其功能，也可以自行指定要用哪些按键组合。
参　　数：
　-d 　显示按键配置的内容。 
  -f<按键配置文件> 　载入指定的按键配置文件。 
　-l 　列出所有的功能。 
　-m<按键配置> 　指定按键配置。 
　-q<功能> 　显示指定功能的按键。 
　-v 　列出目前的按键配置与其功能。
### chkconfig(check config)：检查，设置系统的各种服务。
语　　法：chkconfig [--add][--del][--list][系统服务] 或 chkconfig [--level <等级代号>][系统服务][on/off/reset]
补充说明：这是Red Hat公司遵循GPL规则所开发的程序，它可查询操作系统在每一个执行等级中会执行哪些系统服务，其中包括各类常驻服务。
参　　数：
　--add 　增加所指定的系统服务，让chkconfig指令得以管理它，并同时在系统启动的叙述文件内增加相关数据。 
　--del 　删除所指定的系统服务，不再由chkconfig指令管理，并同时在系统启动的叙述文件内删除相关数据。 
　--level <等级代号> 指定读系统服务要在哪一个执行等级中开启或关毕
### chroot(change root)：改变根目录。
语　　法：chroot [--help][--version][目的目录][执行指令...]
补充说明：把根目录换成指定的目的目录。
参　　数：
　--help 　在线帮助。
　--version 　显示版本信息。
### clock：调整 RTC 时间。
语　　法：
clock [--adjust][--debug][--directisa][--getepoch][--hctosys][--set --date="<日期时间>"][--setepoch --epoch=< >][--show][--systohc][--test][--utc][--version]
补充说明：RTC 是电脑内建的硬件时间，执行这项指令可以显示现在时刻，调整硬件时钟的时间，将系统时间设成与硬件时钟之时间一致，或是把系统时间回存到硬件时钟。
参　　数：
　--adjust 　第一次使用"--set"或"--systohc"参数设置硬件时钟，会在/etc目录下产生一个名称为adjtime的文件。当再次使用这两个参数调整硬件时钟，此文件便会记录两次调整间之差异，日后执行clock指令加上"--adjust"参数时，程序会自动根 据记录文件的数值差异，计算出平均值，自动调整硬件时钟的时间。 
　--debug 　详细显示指令执行过程，便于排错或了解程序执行的情形。 
　--directisa 　告诉clock指令不要通过/dev/rtc设备文件，直接对硬件时钟进行存取。这个参数适用于仅有ISA总线结构的老式电脑。 
　--getepoch 　把系统核心内的硬件时钟新时代数值，呈现到标准输出设备。 
　--hctosys 　Hardware Clock to System Time，把系统时间设成和硬件时钟一致。由于这个动作将会造成系统全面更新文件的存取时间，所以最好在系统启动时就执行它。 
　--set--date 　设置硬件时钟的日期和时间。 
　--setepoch--epoch=<年份>  　设置系统核心之硬件时钟的新时代数值，年份以四位树字表示。 
　--show 　读取硬件时钟的时间，并将其呈现至标准输出设备。 
　--systohc 　System Time to Hardware Clock，将系统时间存回硬件时钟内。 
　--test 　仅作测试，并不真的将时间写入硬件时钟或系统时间。 
　--utc 　把硬件时钟上的时间时为CUT，有时也称为UTC或UCT。 
　--version 　显示版本信息。
### crontab：设置计时器。
语　　法：crontab [-u <用户名称>][配置文件] 或 crontab [-u <用户名称>][-elr]
补充说明：cron是一个常驻服务，它提供计时器的功能，让用户在特定的时间得以执行预设的指令或程序。只要用户会编辑计时器的配置文件，就可以使用计时器的功能。其配置文件格式如下：
  Minute Hour Day Month DayOFWeek Command 
 参　　数： 
　-e 　编辑该用户的计时器设置。 
　-l 　列出该用户的计时器设置。 
　-r 　删除该用户的计时器设置。 
　-u<用户名称> 　指定要设定计时器的用户名称。
### declare：声明 shell 变量。
语　　法：declare [+/-][rxi][变量名称＝设置值] 或 declare -f
补充说明：declare为shell指令，在第一种语法中可用来声明变量并设置变量的属性([rix]即为变量的属性），在第二种语法中可用来显示shell函数。若不加上任何参数，则会显示全部的shell变量与函数(与执行set指令的效果相同)。
参　　数：
　+/- 　"-"可用来指定变量的属性，"+"则是取消变量所设的属性。 
　-f 　仅显示函数。 
　r 　将变量设置为只读。 
　x 　指定的变量会成为环境变量，可供shell以外的程序来使用。 
　i 　[设置值]可以是数值，字符串或运算式。
### depmod(depend module)：分析可载入模块的相依性。
语　　法：depmod [-adeisvV][-m <文件>][--help][模块名称]
补充说明：depmod可检测模块的相依性，供modprobe在安装模块时使用。
参　　数：
　-a或--all 　分析所有可用的模块。 
　-d或debug 　执行排错模式。 
　-e 　输出无法参照的符号。 
　-i 　不检查符号表的版本。 
　-m<文件>或system-map<文件> 　使用指定的符号表文件。 
　-s或--system-log 　在系统记录中记录错误。 
　-v或--verbose 　执行时显示详细的信息。 
　-V或--version 　显示版本信息。 
　--help 　显示帮助。
### dircolors：设置 ls 指令在显示目录或文件时所用的色彩。
语　　法：dircolors [色彩配置文件] 或 dircolors [-bcp][--help][--version]
补充说明：dircolors可根据[色彩配置文件]来设置LS_COLORS环境变量或是显示设置LS_COLORS环境变量的shell指令。
参　　数：
　-b或--sh或--bourne-shell 　显示在Boume shell中，将LS_COLORS设为目前预设置的shell指令。 
　-c或--csh或--c-shell 　显示在C shell中，将LS_COLORS设为目前预设置的shell指令。 
　-p或--print-database 　显示预设置 
　-help 　显示帮助。 
　-version 　显示版本信息。
### dmesg：显示开机信息。
语　　法：dmesg [-cn][-s <缓冲区大小>]
补充说明：kernel会将开机信息存储在ring buffer中。您若是开机时来不及查看信息，可利用dmesg来查看。开机信息亦保存在/var/log目录中，名称为dmesg的文件里。
参　　数：
　-c 　显示信息后，清除ring buffer中的内容。 
　-s<缓冲区大小> 　预设置为8196，刚好等于ring buffer的大小。 
　-n 　设置记录信息的层级。
### enable：启动或关闭 shell 内建指令。
语　　法：enable [-n][-all][内建指令]
补充说明：若要执行的文件名称与shell内建指令相同，可用enable -n来关闭shell内建指令。若不加-n参数，enable可重新启动关闭的指令。
参　　数：
　-n 　关闭指定的shell内建指令。 
　-all 　显示shell所有关闭与启动的指令。
### eval：重新运算求出参数的内容。
语　　法：eval [参数]
补充说明：eval可读取一连串的参数，然后再依参数本身的特性来执行。
参　　数：参数不限数目，彼此之间用分号分开。
### export：设置或显示环境变量。
语　　法：export [-fnp][变量名称]=[变量设置值]
补充说明：在shell中执行程序时，shell会提供一组环境变量。export可新增，修改或删除环境变量，供后续执行的程序使用。export的效力仅及于该此登陆操作。
参　　数：
　-f 　代表[变量名称]中为函数名称。 
　-n 　删除指定的变量。变量实际上并未删除，只是不会输出到后续指令的执行环境中。 
　-p 　列出所有的shell赋予程序的环境变量。
### fbset(frame buffer setup)：设置景框缓冲区。
语　　法：fbset [-ahinsvVx][-db <信息文件>][-fb <外围设备代号>][--test][显示模式]
补充说明：fbset指令可用于设置景框缓冲区的大小，还能调整画面之分辨率，位置，高低宽窄，色彩 深度，并可决定是否启动先卡之各项硬件特性。
参　　数：
　-a或--all 　改变所有使用该设备之虚拟终端机的显示模式。 
　-db<信息文件> 　指定显示模式的信息文件，预设值文件名称为fb.modes，存放在/etc目录下  
　-fb<外围设备代号> 　指定用来做为输出景框缓冲区之外围设备，预设置为"/dev/fd0"。 
　-h或-help 　在线帮助。 
　-i或--info 　列出所有景框缓冲区之相关信息。 
　-ifb<外围设备代号> 　使用另一个景框缓冲区外围设备之设置值。 
　-n或--now 　马上改变显示模式。 
　-ofb<外围设备代号> 　此参数效果和指定"-fb"参数相同。 
　-s或--show 　列出目前显示模式之设置。 
　-v或--verbose 　显示指令执行过程。 
　-V或--version 　显示版本信息。 
　-x或--xfree86 　使用XFree86兼容模式。 
　--test 　仅做测试，并不改变现行的显示模式
### grpconv(group convert to shadow password)：开启群组的投影密码。
语　　法：grpconv
补充说明：Linux系统里的用户和群组密码，分别存放在/etc目录下的passwd和group文件中。因系统运作所需，任何人都得以读取它们，造成安全上的破绽。投影密码将文件内的密码改存在/etc目录下的shadow和gshadow文件内，只允许系统管理者读取，同时把原密码置换为"x"字符。投影密码的功能可随时开启或关闭，您只需执行grpconv指令就能开启群组投影密码。
### grpunconv(group unconvert from shadow password)：关闭群组的投影密码。
语　　法：grpunconv
补充说明：执行grpunconv指令可关闭群组投影密码，它会把密码从gshadow文件内，回存到group文件里。
### hwclock(hardware clock)：显示与设定硬件时钟。
语　　法：
hwclock [--adjust][--debug][--directisa][--hctosys][--show][--systohc][--test]
[--utc][--version][--set --date=<日期与时间>]
补充说明：在Linux中有硬件时钟与系统时钟等两种时钟。硬件时钟是指主机板上的时钟设备，也就是通常可在BIOS画面设定的时钟。系统时钟则是指kernel中的时钟。当Linux启动时，系统时钟会去读取硬件时钟的设定，之后系统时钟即独立运作。所有Linux相关指令与函数都是读取系统时钟的设定。
参　　数：
  --adjust 　hwclock每次更改硬件时钟时，都会记录在/etc/adjtime文件中。使用--adjust参数，可使hwclock根据先前的记录来估算硬件时钟的偏差，并用来校正目前的硬件时钟。 
  --debug 　显示hwclock执行时详细的信息。 
  --directisa 　hwclock预设从/dev/rtc设备来存取硬件时钟。若无法存取时，可用此参数直接以I/O指令来存取硬件时钟。 
  --hctosys 　将系统时钟调整为与目前的硬件时钟一致。 
  --set --date=<日期与时间> 　设定硬件时钟。 
  --show 　显示硬件时钟的时间与日期。 
  --systohc 　将硬件时钟调整为与目前的系统时钟一致。 
  --test 　仅测试程序，而不会实际更改硬件时钟。 
  --utc 　若要使用格林威治时间，请加入此参数，hwclock会执行转换的工作。 
  --version 　显示版本信息。
### insmod(install module)：载入模块。
语　　法：insmod [-fkmpsvxX][-o <模块名称>][模块文件][符号名称 = 符号值]
补充说明：Linux有许多功能是通过模块的方式，在需要时才载入kernel。如此可使kernel较为精简，进而提高效率，以及保有较大的弹性。这类可载入的模块，通常是设备驱动程序。
参　　数：
  -f 　不检查目前kernel版本与模块编译时的kernel版本是否一致，强制将模块载入。 
  -k 　将模块设置为自动卸除。 
  -m 　输出模块的载入信息。 
  -o<模块名称> 　指定模块的名称，可使用模块文件的文件名。 
  -p 　测试模块是否能正确地载入kernel。 
  -s 　将所有信息记录在系统记录文件中。 
  -v 　执行时显示详细的信息。 
  -x 　不要汇出模块的外部符号。 
  -X 　汇出模块所有的外部符号，此为预设置
### kbdconfig：设置键盘类型。
语　　法：kbdconfig [--back][--test]
补充说明：kbdconfig(Red Hat Linux才有的指令)是一个用来设置键盘的程序，提供图形化的操作界面。kbdconfig实际上是修改/etc/sysconfig/keyboard的键盘配置文件。
参　　数：
  --back 　执行时将预设的Cancel按钮更改为Back按钮。 
  --test 　仅作测试，不会实际更改设置。
### lilo(linux loader)：安装核心载入，开机管理程序。
语　　法：lilo [-clqtV][-b<外围设备代号>][-C<配置文件>][-d<延迟时间>][－D<识别标签>][-f<几何参数文件>][-i<开机磁区文件>][-I<识别标签>][-m<映射文件>][-P<fix/ignore>][-r<根目录>][-R<执行指令>...][-s<备份文件>][-S<备份文件>][-uU<外围设备代号>][-v...]
补充说明：lilo是个Linux系统核心载入程序，同时具备管理开机的功能。单独执行lilo指令，它会读取/etc/目录下的lilo.conf配置文件，然后根据其内容安装lilo。
参　　数：
  -b<外围设备代号> 　指定安装lilo之处的外围设备代号。 
  -c 　使用紧致映射模式。 
  -C<配置文件> 　指定lilo的配置文件。 
  -d<延迟时间> 　设置开机延迟时间。 
  -D<识别标签> 　指定开机后预设启动的操作系统，或系统核心识别标签。 
  -f<几何参数文件> 　指定磁盘的几何参数配置文件。 
  -i<开机磁区文件> 　指定欲使用的开机磁区文件，预设是/boot目录里的boot.b文件。 
  -I<识别标签> 　显示系统核心存放之处。 
  -l 　产生线形磁区地址。 
  -m<映射文件> 　指定映射文件。 
  -P<fix/ignore> 　决定要修复或忽略分区表的错误。 
  -q 　列出映射的系统核心文件。 
  -r<根目录> 　设置系统启动时欲挂入成为根目录的目录。 
  -R<执行指令> 　设置下次启动系统时，首先执行的指令。 
  -s<备份文件> 　指定备份文件。 
  -S<备份文件> 　强制指定备份文件。 
  -t 　不执行指令，仅列出实际执行会进行的动作。 
  -u<外围色设备代号> 　删除lilo。 
  -U<外围设备代号> 　此参数的效果和指定"-u"参数类似，当不检查时间戳记。 
  -v 　显示指令执行过程。 
  -V 　显示版本信息。
### liloconfig：设置核心载入，开机管理程序。
语　　法：liloconfig
补充说明：liloconfig是Slackware发行版专门用来调整lilo设置的程序。它通过互动式操作界面，让用户能够利用键盘上的方向键等，轻易地操控lilo的安装，设置作业，而无须下达各种参数或撰写配置文件。
### lsmod(list modules)：显示已载入系统的模块。
语　　法：lsmod
补充说明：执行lsmod指令，会列出所有已载入系统的模块。Linux操作系统的核心具有模块化的特性，应此在编译核心时，务须把全部的功能都放入核心。您可以将这些功能编译成一个个单独的模块，待需要时再分别载入。
### minfo：显示MS-DOS文件系统的各项参数。
语　　法：</>minfo [-v][驱动器代号]
补充说明：minfo为mtools工具指令，可显示MS-DOS系统磁盘的各项参数，包括磁区数，磁头数...等。
参　　数：
  -v 　除了一般信息外，并显示可开机磁区的内容。
### mkkickstart：建立安装的组态文件。
语　　法：mkkickstart [--bootp][--dhcp][--nonet][--nox][--version][--nfs <远端电脑:路径>]
 补充说明：mkkickstart可根据目前系统的设置来建立组态文件，供其他电脑在安装时使用。组态文件的内容包括使用语言，网络环境，系统磁盘状态，以及X Windows的设置等信息。
参　　数：
  --bootp 　安装与开机时，使用BOOTP。 
  --dhcp 　安装与开机时，使用DHCP。 
  --nfs<远端电脑:路径> 　使用指定的网络路径安装。 
  --nonet 　不要进行网络设置，即假设在没有网络环境的状态下。 
  --nox 　不要进行X Windows的环境设置。 
  --version 　显示版本信息。
### modinfo(module infomation)：显示kernel模块的信息。
语　　法：modinfo [-adhpV][模块文件]
补充说明：modinfo会显示kernel模块的对象文件，以显示该模块的相关信息。
参　　数：
  -a或--author 　显示模块开发人员。 
  -d或--description 　显示模块的说明。 
  -h或--help 　显示modinfo的参数使用方法。 
  -p或--parameters 　显示模块所支持的参数。 
  -V或--version 　显示版本信息。
### modprobe(module probe)：自动处理可载入模块。
语　　法：modprobe [-acdlrtvV][--help][模块文件][符号名称 = 符号值]
补充说明：modprobe可载入指定的个别模块，或是载入一组相依的模块。modprobe会根据depmod所产生的相依关系，决定要载入哪些模块。若在载入过程中发生错误，在modprobe会卸载整组的模块。
参　　数：
  -a或--all 　载入全部的模块。 
  -c或--show-conf 　显示所有模块的设置信息。 
  -d或--debug 　使用排错模式。 
  -l或--list 　显示可用的模块。 
  -r或--remove 　模块闲置不用时，即自动卸载模块。 
  -t或--type 　指定模块类型。 
  -v或--verbose 　执行时显示详细的信息。 
  -V或--version 　显示版本信息。 
  -help 　显示帮助。
### mouseconfig：设置鼠标相关参数。
语　　法：
### mouseconfig [--back][--emulthree][--help][--expert][--kickstart][--noprobe][--test][--device <连接端口>][鼠标类型]
补充说明：mouseconfig为鼠标设置程序，可自动设置相关参数，或者用户也可以利用所提供互动模式自行设置鼠标。mouseconfig是Red Hat Linux才有的命令。
参　　数：
  --back 　在设置画面上显示Back按钮，而取代预设的Cancel按钮。 
  --device<连接端口> 　指定硬件连接端口。可用的选项有ttyS0，ttyS1，ttyS2，ttyS3与orpsaux。 
  --emulthree 　将二钮鼠标模拟成三钮鼠标。 
  --help 　显示帮助以及所有支持的鼠标类型。 
  --expert 　程序预设可自动判断部分设置值。若要自行设置，请使用--expert参数。 
  --kickstart 　让程序自动检测并保存所有的鼠标设置。 
  --noprobe 　不要检测鼠标设备。 
  --test 　测试模式，不会改变任何设置。
### ntsysv：设置系统的各种服务。
语　　法：ntsysv [--back][--level <等级代号>]
补充说明：这是Red Hat公司遵循GPL规则所开发的程序，它具有互动式操作界面，您可以轻易地利用方向键和空格键等，开启，关闭操作系统在每个执行等级中，所要执行的系统服务。
参　　数：
  --back 　在互动式界面里，显示Back钮，而非Cancel钮。 
  --level <等级代号> 　在指定的执行等级中，决定要开启或关闭哪些系统服务。
### passwd(password)：设置密码。
语　　法：passwd [-dklS][-u <-f>][用户名称]
补充说明：passwd指令让用户可以更改自己的密码，而系统管理者则能用它管理系统用户的密码。只有管理者可以指定用户名称，一般用户只能变更自己的密码。
参　　数：
  -d 　删除密码。本参数仅有系统管理者才能使用。 
  -f 　强制执行。 
  -k 　设置只有在密码过期失效后，方能更新。 
  -l 　锁住密码。 
  -s 　列出密码的相关信息。本参数仅有系统管理者才能使用。 
  -u 　解开已上锁的帐号。
### pwconv：开启用户的投影密码。
语　　法：pwconv
补充说明：Linux系统里的用户和群组密码，分别存放在名称为passwd和group的文件中，　这两个文件位于/etc目录下。因系统运作所需，任何人都得以读取它们，造成安全上的破绽。投影密码将文件内的密码改存在/etc目录下的shadow和gshadow文件内，只允许系统管理者读取，同时把原密码置换为"x"字符，有效的强化了系统的安全性。
### pwunconv：关闭用户的投影密码。
语　　法：pwunconv
补充说明：执行pwunconv指令可以关闭用户投影密码，它会把密码从shadow文件内，重回存到passwd文件里。
### rdate(receive date)：显示其他主机的日期与时间。
语　　法：rdate [-ps][主机名称或IP地址...]
补充说明：执行rdate指令，向其他主机询问系统时间并显示出来。
参　　数：
  -p 　显示远端主机的日期与时间。 
  -s 　把从远端主机收到的日期和时间，回存到本地主机的系统时间。
### resize：设置终端机视窗的大小。
语　　法：resize [-cu][-s <列数> <行数>]
补充说明：执行resize指令可设置虚拟终端机的视窗大小。
参　　数：
　-c 　就算用户环境并非C Shell，也用C Shell指令改变视窗大小。 
　-s <列数> <行数> 　设置终端机视窗的垂直高度和水平宽度。 
　-u 　就算用户环境并非Bourne Shell，也用Bourne Shell指令改变视窗大小。
### rmmod(remove module)：删除模块。
语　　法：rmmod [-as][模块名称...]
补充说明：执行rmmod指令，可删除不需要的模块。Linux操作系统的核心具有模块化的特性，应此在编译核心时，务须把全部的功能都放如核心。你可以将这些功能编译成一个个单独的模块，待有需要时再分别载入它们。
参　　数：
  -a 　删除所有目前不需要的模块。 
  -s 　把信息输出至syslog常驻服务，而非终端机界面。
### rpm(redhat package manager)：管理套件。
语　　法：rpm [-acdhilqRsv][-b<完成阶段><套间档>+][-e<套件挡>][-f<文件>+][-i<套件档>][-p<套件档>＋][-U<套件档>][-vv][--addsign<套件档>+][--allfiles][--allmatches][--badreloc][--buildroot<根目录>][--changelog][--checksig<套件档>+][--clean][--dbpath<数据库目录>][--dump][--excludedocs][--excludepath<排除目录>][--force][--ftpproxy<主机名称或IP地址>][--ftpport<通信端口>][--help][--httpproxy<主机名称或IP地址>][--httpport<通信端口>][--ignorearch][--ignoreos][--ignoresize][--includedocs][--initdb][justdb][--nobulid][--nodeps][--nofiles][--nogpg][--nomd5][--nopgp][--noorder][--noscripts][--notriggers][--oldpackage][--percent][--pipe<执行指令>][--prefix<目的目录>][--provides][--queryformat<档头格式>][--querytags][--rcfile<配置档>][--rebulid<套件档>][--rebuliddb][--recompile<套件档>][--relocate<原目录>=<新目录>][--replacefiles][--replacepkgs][--requires][--resign<套件档>+][--rmsource][--rmsource<文件>][--root<根目录>][--scripts][--setperms][--setugids][--short-circuit][--sign][--target=<安装平台>+][--test][--timecheck<检查秒数>][--triggeredby<套件档>][--triggers][--verify][--version][--whatprovides<功能特性>][--whatrequires<功能特性>]
补充说明：rmp原本是Red Hat Linux发行版专门用来管理Linux各项套件的程序，由于它遵循GPL规则且功能强大方便，因而广受欢迎。逐渐受到其他发行版的采用。RPM套件管理方式的出现，让Linux易于安装，升级，间接提升了Linux的适用度。
参　　数：
  -a 　查询所有套件。
  -b<完成阶段><套件档>+或-t <完成阶段><套件档>+ 　设置包装套件的完成阶段，并指定套件档的文件名称。
  -c 　只列出组态配置文件，本参数需配合"-l"参数使用。 
  -d 　只列出文本文件，本参数需配合"-l"参数使用。　 
  -e<套件档>或--erase<套件档> 　删除指定的套件。
  -f<文件>+ 　查询拥有指定文件的套件。 
  -h或--hash 　套件安装时列出标记。 
  -i 　显示套件的相关信息。 
  -i<套件档>或--install<套件档> 　安装指定的套件档。 
  -l 　显示套件的文件列表。 
  -p<套件档>+ 　查询指定的RPM套件档。 
  -q 　使用询问模式，当遇到任何问题时，rpm指令会先询问用户。 
  -R 　显示套件的关联性信息。 
  -s 　显示文件状态，本参数需配合"-l"参数使用。 
  -U<套件档>或--upgrade<套件档>   升级指定的套件档。 
  -v 　显示指令执行过程。
  -vv 　详细显示指令执行过程，便于排错。 
  -addsign<套件档>+ 　在指定的套件里加上新的签名认证。 
  --allfiles 　安装所有文件。 
  --allmatches 　删除符合指定的套件所包含的文件。 
  --badreloc 　发生错误时，重新配置文件。 
  --buildroot<根目录> 　设置产生套件时，欲当作根目录的目录。 
  --changelog 　显示套件的更改记录。 
  --checksig<套件档>+ 　检验该套件的签名认证。 
  --clean 　完成套件的包装后，删除包装过程中所建立的目录。 
  --dbpath<数据库目录> 　设置欲存放RPM数据库的目录。 
  --dump 　显示每个文件的验证信息。本参数需配合"-l"参数使用。 
  --excludedocs 　安装套件时，不要安装文件。 
  --excludepath<排除目录> 　忽略在指定目录里的所有文件。 
  --force 　强行置换套件或文件。 
  --ftpproxy<主机名称或IP地址> 　指定FTP代理服务器。 
  --ftpport<通信端口> 　设置FTP服务器或代理服务器使用的通信端口。 
  --help 　在线帮助。 
  --httpproxy<主机名称或IP地址> 　指定HTTP代理服务器。 
  --httpport<通信端口> 　设置HTTP服务器或代理服务器使用的通信端口。 
  --ignorearch 　不验证套件档的结构正确性。 
  --ignoreos 　不验证套件档的结构正确性。 
  --ignoresize 　安装前不检查磁盘空间是否足够。 
  --includedocs 　安装套件时，一并安装文件。 
  --initdb 　确认有正确的数据库可以使用。 
  --justdb 　更新数据库，当不变动任何文件。 
  --nobulid 　不执行任何完成阶段。 
  --nodeps 　不验证套件档的相互关联性。 
  --nofiles 　不验证文件的属性。 
  --nogpg 　略过所有GPG的签名认证。 
  --nomd5 　不使用MD5编码演算确认文件的大小与正确性。 
  --nopgp 　略过所有PGP的签名认证。 
  --noorder 　不重新编排套件的安装顺序，以便满足其彼此间的关联性。 
  --noscripts 　不执行任何安装Script文件。 
  --notriggers 　不执行该套件包装内的任何Script文件。 
  --oldpackage 　升级成旧版本的套件。 
  --percent 　安装套件时显示完成度百分比。 
  --pipe<执行指令> 　建立管道，把输出结果转为该执行指令的输入数据。 
  --prefix<目的目录> 　若重新配置文件，就把文件放到指定的目录下。 
  --provides 　查询该套件所提供的兼容度。 
  --queryformat<档头格式> 　设置档头的表示方式。 
  --querytags 　列出可用于档头格式的标签。 
  --rcfile<配置文件> 　使用指定的配置文件。 
  --rebulid<套件档> 　安装原始代码套件，重新产生二进制文件的套件。 
  --rebuliddb 　以现有的数据库为主，重建一份数据库。 
  --recompile<套件档> 　此参数的效果和指定"--rebulid"参数类似，当不产生套件档。 
  --relocate<原目录>=<新目录> 　把本来会放到原目录下的文件改放到新目录。 
  --replacefiles 　强行置换文件。 
  --replacepkgs 　强行置换套件。 
  --requires 　查询该套件所需要的兼容度。 
  --resing<套件档>+ 　删除现有认证，重新产生签名认证。 
  --rmsource 　完成套件的包装后，删除原始代码。 
  --rmsource<文件> 　删除原始代码和指定的文件。 
  --root<根目录> 　设置欲当作根目录的目录。 
  --scripts 　列出安装套件的Script的变量。 
  --setperms 　设置文件的权限。 
  --setugids 　设置文件的拥有者和所属群组。 
  --short-circuit 　直接略过指定完成阶段的步骤。
  --sign 　产生PGP或GPG的签名认证。 
  --target=<安装平台>+ 　设置产生的套件的安装平台。 
  --test 　仅作测试，并不真的安装套件。 
  --timecheck<检查秒数> 　设置检查时间的计时秒数。 
  --triggeredby<套件档> 　查询该套件的包装者。 
  --triggers 　展示套件档内的包装Script。 
  --verify 　此参数的效果和指定"-q"参数相同。 
  --version 　显示版本信息。 
  --whatprovides<功能特性> 　查询该套件对指定的功能特性所提供的兼容度。 
  --whatrequires<功能特性> 　查询该套件对指定的功能特性所需要的兼容度。
### set：设置shell。
语　　法：set [+-abCdefhHklmnpPtuvx]
补充说明：set指令能设置所使用shell的执行方式，可依照不同的需求来做设置。
参　　数：
　-a 　标示已修改的变量，以供输出至环境变量。 
　-b 　使被中止的后台程序立刻回报执行状态。 
　-C 　转向所产生的文件无法覆盖已存在的文件。 
　-d 　Shell预设会用杂凑表记忆使用过的指令，以加速指令的执行。使用-d参数可取消。 
　-e 　若指令传回值不等于0，则立即退出shell。　　 
　-f　 　取消使用通配符。 
　-h 　自动记录函数的所在位置。 
　-H Shell 　可利用"!"加<指令编号>的方式来执行history中记录的指令。 
　-k 　指令所给的参数都会被视为此指令的环境变量。 
　-l 　记录for循环的变量名称。 
　-m 　使用监视模式。 
　-n 　只读取指令，而不实际执行。 
　-p 　启动优先顺序模式。 
　-P 　启动-P参数后，执行指令时，会以实际的文件或目录来取代符号连接。 
　-t 　执行完随后的指令，即退出shell。 
　-u 　当执行时使用到未定义过的变量，则显示错误信息。 
　-v 　显示shell所读取的输入值。 
　-x 　执行指令后，会先显示该指令及所下的参数。 
　+<参数> 　取消某个set曾启动的参数。
### setconsole：设置系统终端。
语　　法：setconsole [video][serial][ttya][ttyb]
补充说明：setconsole可用来指定系统终端。
参　　数：
　serial 　使用PROM终端。 
　ttya,cua0或ttyS0 　使用第１个串口设备作为终端。 
　ttyb,cua1或ttyS1 　使用第２个串口设备作为终端。 
　video 　使用主机上的现卡作为终端。
### setenv(set environment variable)：查询或显示环境变量。
语　　法：setenv [变量名称][变量值]
补充说明：setenv为tsch中查询或设置环境变量的指令
### setup：设置公用程序。
语　　法：setup
补充说明：setup是一个设置公用程序，提供图形界面的操作方式。在setup中可设置7类的选项：
　１.登陆认证方式 
　２.键盘组态设置 
　３.鼠标组态设置 
　４.开机时所要启动的系统服务 
　５.声卡组态设置 
　６.时区设置 
　７.X Windows组态设置
### sndconfig：设置声卡。
语　　法：sndconfig [--help][--noautoconfig][--noprobe]
补充说明：sndconfig为声卡设置程序，支持PnP柚茫勺远觳獠⑸柚肞nP声卡。
参 　 数：
  --help 　显示帮助。 
  --noautoconfig 　不自动设置PnP的声卡。 
  --noprobe 　不自动检测PnP声卡。
### SVGAText Mode：加强文字模式的显示画面。
语　　法：SVGATextMode [-acdfhmnrsv][-t <配置文件>][模式]
补充说明：SVGATextMode可用来设置文字模式下的显示画面，包括分辨率，字体和更新频率等。
参　　数：
　-a 　如果新显示模式的屏幕大小与原先不同时，SVGATextMode会执行必要的系统设置。 
　-c 　维持原有的VGA时脉。 
　-d 　执行时会显示详细的信息，供排错时参考。 
　-f 　不要执行配置文件中有关字体载入的指令。 
　-h 　显示帮助。 
　-m 　允许1x1的方式来重设屏幕大小。 
　-n 　仅测试指定的模式。 
　-r 　通知或重设与屏幕大小相关的程序。 
　-s 　显示配置文件中所有可用的模式。 
　-t<配置文件> 　指定配置文件。 
　-v 　SVGATextMode在配置新的显示模式时，预设会先检查垂直与水平的更新更新频率是否在配置文件所指定的范围内，如果不在范围内，则不设置新的显示模式。 
　模式］ 　[模式]参数必须是配置文件中模式的名称。
### timeconfig：设置时区。
语　　法：timeconfig [--arc][--back][--test][--utc][时区名称]
补充说明：这是Red Hat公司遵循GPL规则所开发的程序，它具有互动式操作界面，您可以轻易地利用方向键和空格键等，设置系统时间所属的时区。
参　　数：
　--arc 　使用Alpha硬件结构的格式存储系统时间。 
　--back 　在互动式界面里，显示Back钮而非Cancel钮。 
　--test 　仅作测试，并不真的改变系统的时区。 
　--utc 　把硬件时钟上的时间视为CUT，有时也称为UTC或UCT。
### ulimit：控制shell程序的资源。
语　　法：ulimit [-aHS][-c <core文件上限>][-d <数据节区大小>][-f <文件大小>][-m <内存大小>][-n <文件数目>][-p <缓冲区大小>][-s <堆叠大小>][-t <CPU时间>][-u <程序数目>][-v <虚拟内存大小>]
补充说明：ulimit为shell内建指令，可用来控制shell执行程序的资源。
参　　数：
  -a 　显示目前资源限制的设定。 
  -c <core文件上限> 　设定core文件的最大值，单位为区块。 
  -d <数据节区大小> 　程序数据节区的最大值，单位为KB。 
  -f <文件大小> 　shell所能建立的最大文件，单位为区块。 
  -H 　设定资源的硬性限制，也就是管理员所设下的限制。 
  -m <内存大小> 　指定可使用内存的上限，单位为KB。 
  -n <文件数目> 　指定同一时间最多可开启的文件数。 
  -p <缓冲区大小> 　指定管道缓冲区的大小，单位512字节。 
  -s <堆叠大小> 　指定堆叠的上限，单位为KB。 
  -S 　设定资源的弹性限制。 
  -t <CPU时间> 　指定CPU使用时间的上限，单位为秒。 
  -u <程序数目> 　用户最多可开启的程序数目。 
  -v <虚拟内存大小> 　指定可使用的虚拟内存上限，单位为KB。
### unalias：删除别名。
语　　法：unalias [-a][别名]
补充说明：unalias为shell内建指令，可删除别名设置。
参　　数：
  -a 　删除全部的别名。
### unset：删除变量或函数。
语　　法：unset [-fv][变量或函数名称]
补充说明：unset为shell内建指令，可删除变量或函数。
参 　 数：
  -f 　仅删除函数。 
  -v 　仅删除变量。
 
## linux文档编辑命令
### col：过滤控制字符。
语　　法：col [-bfx][-l<缓冲区列数>]
补充说明：在许多UNIX说明文件里，都有RLF控制字符。当我们运用shell特殊字符">"和">>"，把说明文件的内容输出成纯文本文件时，控制字符会变成乱码，col指令则能有效滤除这些控制字符。
参　　数：
  -b   过滤掉所有的控制字符，包括RLF和HRLF。
  -f   滤除RLF字符，但允许将HRLF字符呈现出来。
  -x   以多个空格字符来表示跳格字符。
  -l<缓冲区列数>   预设的内存缓冲区有128列，您可以自行指定缓冲区的大小。
### colrm(column remove)：滤掉指定的行。
语　　法：colrm [开始行数编号<结束行数编号>]
补充说明：colrm指令从标准输入设备读取书记，转而输出到标准输出设备。如果不加任何参数，则该指令不会过滤任何一行。
### comm(common)：比较两个已排过序的文件。
语　　法：comm [-123][--help][--version][第1个文件][第2个文件]
补充说明：这项指令会一列列地比较两个已排序文件的差异，并将其结果显示出来，如果没有指定任何参数，则会把结果分成3行显示：第1行仅是在第1个文件中出现过的列，第2行是仅在第2个文件中出现过的列，第3行则是在第1与第2个文件里都出现过的列。若给予的文件名称为"-"，则comm指令会从标准输入设备读取数据。
参　　数：
  -1   不显示只在第1个文件里出现过的列。
  -2   不显示只在第2个文件里出现过的列。
  -3   不显示只在第1和第2个文件里出现过的列。
  --help   在线帮助。
  --version   显示版本信息。
### csplit(context split)：分割文件。
语　　法：csplit [-kqsz][-b<输出格式>][-f<输出字首字符串>][-n<输出文件名位数>][--help][--version][文件][范本样式...]
补充说明：将文件依照指定的范本样式予以切割后，分别保存成名称为xx00,xx01,xx02...的文件。若给予的文件名称为"-"，则csplit指令会从标准输入设备读取数据。
参　　数：
  -b<输出格式>或--suffix-format=<输出格式>   预设的输出格式其文件名称为xx00,xx01...等，您可以通过改变<输出格式>来改变输出的文件名。
  -f<输出字首字符串>或--prefix=<输出字首字符串>   预设的输出字首字符串其文件名为xx00,xx01...等，如果你指定输出字首字符串为"hello"，则输出的文件名称会变成hello00,hello01...等。
  -k或--keep-files   保留文件，就算发生错误或中断执行，也不能删除已经输出保存的文件。
  -n<输出文件名位数>或--digits=<输出文件名位数>   预设的输出文件名位数其文件名称为xx00,xx01...等，如果你指定输出文件名位数为"3"，则输出的文件名称会变成xx000,xx001...等。
  -q或-s或--quiet或--silent   不显示指令执行过程。
  -z或--elide-empty-files   删除长度为0 Byte文件。
  --help   在线帮助。
  --version   显示版本信息。
### ed(editor)：文本编辑器。
语　　法：ed [-][-Gs][-p<字符串>][--help][--version][文件]
补充说明：ed是Linux中功能最简单的文本编辑程序，一次仅能编辑一行而非全屏幕方式的操作。
参　　数：
  -G或--traditional   提供回兼容的功能。
  -p<字符串>   指定ed在command mode的提示字符。
  -s,-,--quiet或--silent   不执行开启文件时的检查功能。
  --help   显示帮助。
  --version   显示版本信息。
### egrep：在文件内查找指定的字符串。egrep执行效果如grep -E，使用的语法及参数可参照grep指令，与grep不同点在于解读字符串的方法，egrep是用extended regular expression语法来解读，而grep则用basic regular expression语法，extended regular expression比basic regular expression有更完整的表达规范。
### ex：在Ex模式下启动vim文本编辑器。ex执行效果如同vi -E，使用语法及参数可参照vi指令，如要从Ex模式回到普通模式，则在vim中输入:vi或:visual即可。
### fgrep(fixed regexp)：查找文件里符合条件的字符串。
语　　法：fgrep [范本样式][文件或目录...]
补充说明：本指令相当于执行grep指令加上参数"-F"，详见grep指令说明。
### fmt(fromat)：编排文本文件。
语　　法：fmt [-cstu][-p<列起始字符串>][-w<每列字符数>][--help][--version][文件...]
补充说明：fmt指令会从指定的文件里读取内容，将其依照指定格式重新编排后，输出到标准输出设备。若指定的文件名为"-"，则fmt指令会从标准输入设备读取数据。
参　　数：
  -c或--crown-margin   每段前两列缩排。
  -p<列起始字符串>或-prefix=<列起始字符串>   仅合并含有指定字符串的列，通常运用在程序语言的注解方面。
  -s或--split-only   只拆开字数超出每列字符数的列，但不合并字数不足每列字符数的列。
  -t或--tagged-paragraph   每列前两列缩排，但第1列和第2列的缩排格式不同。
  -u或--uniform-spacing   每个字符之间都以一个空格字符间隔，每个句子之间则两个空格字符分隔。
  -w<每列字符数>或--width=<每列字符数>或-<每列字符数>   设置每列的最大字符数。
  --help   在线帮助。
  --version   显示版本信息。
### fold：限制文件列宽。
语　　法：fold [-bs][-w<每列行数>][--help][--version][文件...]
补充说明：fold指令会从指定的文件里读取内容，将超过限定列宽的列加入增列字符后，输出到标准输出设备。若不指定任何文件名称，或是所给予的文件名为“-”，则fold指令会从标准输入设备读取数据。
参　　数：
  -b或--bytes   以Byte为单位计算列宽，而非采用行数编号为单位。
  -s或--spaces   以空格字符作为换列点。
  -w<每列行数>或--width<每列行数>   设置每列的最大行数。
  --help   在线帮助。
  --version   显示版本信息。
### grep：查找文件里符合条件的字符串。
语　　法：grep [-abcEFGhHilLnqrsvVwxy][-A<显示列数>][-B<显示列数>][-C<显示列数>][-d<进行动作>][-e<范本样式>][-f<范本文件>][--help][范本样式][文件或目录...]
补充说明：grep指令用于查找内容包含指定的范本样式的文件，如果发现某文件的内容符合所指定的范本样式，预设grep指令会把含有范本样式的那一列显示出来。若不指定任何文件名称，或是所给予的文件名为“-”，则grep指令会从标准输入设备读取数据。
参　　数：
  -a或--text   不要忽略二进制的数据。
  -A<显示列数>或--after-context=<显示列数>   除了显示符合范本样式的那一列之外，并显示该列之后的内容。
  -b或--byte-offset   在显示符合范本样式的那一列之前，标示出该列第一个字符的位编号。
  -B<显示列数>或--before-context=<显示列数>   除了显示符合范本样式的那一列之外，并显示该列之前的内容。
  -c或--count   计算符合范本样式的列数。
  -C<显示列数>或--context=<显示列数>或-<显示列数>   除了显示符合范本样式的那一列之外，并显示该列之前后的内容。
  -d<进行动作>或--directories=<进行动作>   当指定要查找的是目录而非文件时，必须使用这项参数，否则grep指令将回报信息并停止动作。
  -e<范本样式>或--regexp=<范本样式>   指定字符串做为查找文件内容的范本样式。
  -E或--extended-regexp   将范本样式为延伸的普通表示法来使用。
  -f<范本文件>或--file=<范本文件>   指定范本文件，其内容含有一个或多个范本样式，让grep查找符合范本条件的文件内容，格式为每列一个范本样式。
  -F或--fixed-regexp   将范本样式视为固定字符串的列表。
  -G或--basic-regexp   将范本样式视为普通的表示法来使用。
  -h或--no-filename   在显示符合范本样式的那一列之前，不标示该列所属的文件名称。
  -H或--with-filename   在显示符合范本样式的那一列之前，表示该列所属的文件名称。
  -i或--ignore-case   忽略字符大小写的差别。
  -l或--file-with-matches   列出文件内容符合指定的范本样式的文件名称。
  -L或--files-without-match   列出文件内容不符合指定的范本样式的文件名称。
  -n或--line-number   在显示符合范本样式的那一列之前，标示出该列的列数编号。
  -q或--quiet或--silent   不显示任何信息。
  -r或--recursive   此参数的效果和指定“-d recurse”参数相同。
  -s或--no-messages   不显示错误信息。
  -v或--revert-match   反转查找。
  -V或--version   显示版本信息。
  -w或--word-regexp   只显示全字符合的列。
  -x或--line-regexp   只显示全列符合的列。
  -y   此参数的效果和指定“-i”参数相同。
  --help   在线帮助。
### ispell(interactive spelling checking)：拼字检查程序。
语　　法：ispell [-aAbBClmMnNPStVx][-d<字典文件>][-L<行数>][-p<字典文件>][-w<非字母字符>][-W<字符串长度>][要检查的文件]
补充说明：ispell预设会使用/usr/lib/ispell/english.hash字典文件来检查文本文件。若在检查的文件中找到字典没有的词汇，ispell会建议使用的词汇，或是让你将新的词汇加入个人字典。
参　　数：
  -a   当其他程序输出送到ispell时，必须使用此参数。
  -A   读取到"&Include File&"字符串时，就去检查字符串后所指定文件的内容。
  -b   产生备份文件，文件名为.bak。
  -B   检查连字错误。
  -C   不检查连字错误。
  -d<字典文件>   指定字典文件。
  -l   从标准输入设备读取字符串，结束后显示拼错的词汇。
  -L<行数>   指定内文显示的行数。
  -m   自动考虑字尾的变化。
  -M   进入ispell后，在画面下方显示指令的按键。
  -n   检查的文件为noff或troff的格式。
  -N   进入ispell后，在画面下方不显示指令的按键。
  -p<字典文件>   指定个人字典文件。
  -P   不考虑字尾变化的情形。
  -S   不排序建议取代的词汇。
  -t   检查的文件为TeX或LaTeX的格式。
  -V   非ANSI标准的字符会以"M-^"的方式来显示。
  -w<非字母字符>   检查时，特别挑出含有指定的字符。
  -W<字符串长度>   不检查指定长度的词汇。
  -x   不要产生备份文件。
### jed：编辑文本文件。
语　　法：jed [-2n][-batch][-f<函数>][-g<行数>][-i<文件>][-I<文件>][-s<字符串>][文件]
补充说明：Jed是以Slang所写成的程序，适合用来编辑程序原始代码。
参　　数：
  -2   显示上下两个编辑区。
  -batch   以批处理模式来执行。
  -f<函数>   执行Slang函数。
  -g<行数>   移到缓冲区中指定的行数。
  -i<文件>   将指定的文件载入缓冲区。
  -i<文件>   载入Slang原始代码文件。
  -n   不要载入jed.rc配置文件。
  -s<字符串>   查找并移到指定的字符串。
### joe：编辑文本文件。
语　　法：
joe [-asis][-beep][-csmode][-dopadding][-exask][-force][-help][-keepup][-lightoff][-arking][-mid][-nobackups][-nonotice][-nosta][-noxon][-orphan][-backpath<目录>][-columns<栏位>][-lines<行数>][-pg<行数>][-skiptop<行数>][-autoindent crlf linums overwrite rdonly wordwrap][+<行数>][-indentc<缩排字符>][-istep<缩排字符数>][-keymap<按键配置文件>][-lmargin<栏数>][-rmargin<栏数>][-tab<栏数>][要编辑的文件]
补充说明：Joe是一个功能强大的全屏幕文本编辑程序。操作的复杂度要比Pico高一点，但是功能较为齐全。Joe一次可开启多个文件，每个文件各放在一个编辑区内，并可在文件之间执行剪贴的动作。
参　　数：
以下为程序参数
  -asis   字符码超过127的字符不做任何处理。
  -backpath<目录>   指定备份文件的目录。
  -beep   编辑时，若有错误即发出哗声。
  -columns<栏位>   设置栏数。
  -csmode   可执行连续查找模式。
  -dopadding   是程序跟tty间存在缓冲区。
  -exask   在程序中，执行"Ctrl+k+x"时，会先确认是否要保存文件。
  -force   强制在最后一行的结尾处加上换行符号。
  -help   执行程序时一并显示帮助。
  -keepup   在进入程序后，画面上方为状态列。
  -lightoff   选取的区块在执行完区块命令后，就会回复成原来的状态。
  -lines<行数>   设置行数。
  -marking   在选取区块时，反白区块会随着光标移动。
  -mid   当光标移出画面时，即自动卷页，使光标回到中央。
  -nobackups   不建立备份文件。
  -nonotice   程序执行时，不显示版权信息。
  -nosta   程序执行时，不显示状态列。
  -noxon   尝试取消“Ctrl+s”与“Ctrl+q”键的功能。
  -orphan   若同时开启一个以上的文件，则其他文件会置于独立的缓冲区，而不会另外开启编辑区。
  -pg<行数>   按“PageUp”或“PageDown”换页时，所要保留前一页的行数。
  -skiptop<行数>   不使用屏幕上方指定的行数。
以下为文件参数
  +<行数>   指定开启文件时，光标所在的行数。
  -autoindent   自动缩排。
  -crlf   在换行时，使用CR-LF字符。
  -indentc<缩排字符>   执行缩排时，实际插入的字符。
  -istep<缩排字符数>   每次执行缩排时，所移动的缩排字符数。
  -keymap<按键配置文件>   使用不同的按键配置文件。
  -linums   在每行前面加上行号。
  -lmargin<栏数>   设置左侧边界。
  -overwrite   设置覆盖模式。
  -rmargin<栏数>   设置右侧边界。
  -tab<栏数>   设置tab的宽度。
  -rdonly   以只读的方式开启文件-wordwrap编辑时若超过右侧边界，则自动换行。
### join：将两个文件中，指定栏位内容相同的行连接起来。
语　　法：join [-i][-a<1或2>][-e<字符串>][-o<格式>][-t<字符>][-v<1或2>][-1<栏位>][-2<栏位>][--help][--version][文件1][文件2]
补充说明：找出两个文件中，指定栏位内容相同的行，并加以合并，再输出到标准输出设备。
参　　数：
  -a<1或2>   除了显示原来的输出内容之外，还显示指令文件中没有相同栏位的行。
  -e<字符串>   若[文件1]与[文件2]中找不到指定的栏位，则在输出中填入选项中的字符串。
  -i或--igore-case   比较栏位内容时，忽略大小写的差异。
  -o<格式>   按照指定的格式来显示结果。
  -t<字符>   使用栏位的分隔字符。
  -v<1或2>   跟-a相同，但是只显示文件中没有相同栏位的行。
  -1<栏位>   连接[文件1]指定的栏位。
  -2<栏位>   连接[文件2]指定的栏位。
  --help   显示帮助。
  --version   显示版本信息。
### look：查询单字。
语　　法：look [-adf][-t<字尾字符串>][字首字符串][字典文件]
补充说明：look指令用于英文单字的查询。您仅需给予它欲查询的字首字符串，它会显示所有开头字符串符合该条件的单字。
参　　数：
  -a   使用另一个字典文件web2，该文件也位于/usr/dict目录下。
  -d   只对比英文字母和数字，其余一慨忽略不予比对。
  -f   忽略字符大小写差别。
  -t<字尾字符串>   设置字尾字符串。
### mtype：显示MS-DOS文件的内容。
语　　法：mtype [-st][文件]
补充说明：mtype为mtools工具指令，模拟MS-DOS的type指令，可显示MS-DOS文件的内容。
参　　数：
  -s   去除8位字符码集的第一个位，使它兼容于7位的ASCII。
  -t   将MS-DOS文本文件中的“换行+光标移至行首”字符转换成Linux的换行字符。
### pico：编辑文字文件。
语　　法：pico [-bdefghjkmqtvwxz][-n<间隔秒数>][-o<工作目录>][-r<编辑页宽>][-s<拼字检查器>][+<列数编号>][文件]
补充说明：pico是个简单易用、以显示导向为主的文字编辑程序，它伴随着处理电子邮件和新闻组的程序pine而来。
参　　数：
  -b   开启置换的功能。
  -d   开启删除的功能。
  -e   使用完整的文件名称。
  -f   支持键盘上的F1、F2...等功能键。
  -g   显示光标。
  -h   在线帮助。
  -j   开启切换的功能。
  -k   预设pico在使用剪下命令时，会把光标所在的列的内容全部删除。
  -m   开启鼠标支持的功能，您可用鼠标点选命令列表。
  -n<间隔秒数>   设置多久检查一次新邮件。
  -o<工作目录>   设置工作目录。
  -q   忽略预设值。
  -r<编辑页宽>   设置编辑文件的页宽。
  -s<拼字检查器>   另外指定拼字检查器。
  -t   启动工具模式。
  -v   启动阅读模式，用户只能观看，无法编辑文件的内容。
  -w   关闭自动换行，通过这个参数可以编辑内容很长的列。
  -x   关闭换面下方的命令列表。
  -z   让pico可被Ctrl+z中断，暂存在后台作业里。
  +<列数编号>   执行pico指令进入编辑模式时，从指定的列数开始编辑。
### rgrep(recursive grep)：递归查找文件里符合条件的字符串。
语　　法：rgrep [-?BcDFhHilnNrv][-R<范本样式>][-W<列长度>][-x<扩展名>][--help][--version][范本样式][文件或目录...]
补充说明：rgrep指令的功能和grep指令类似，可查找内容包含指定的范本样式的文件，如果发现某文件的内容符合所指定的范本样式，预设rgrep指令会把含有范本样式的那一列显示出来。
参　　数：
  -?   显示范本样式与范例的说明。
  -B   忽略二进制的数据。
  -c   计算符合范本样式的列数。
  -D   排错模式，只列出指令搜寻的目录清单，而不会读取文件内容。
  -F   当遇到符号连接时，rgrep预设是忽略不予处理，加上本参数后，rgrep指令就会读取该连接所指向的原始文件的内容。
  -h   特别将符合范本样式的字符串标示出来。
  -H   只列出符合范本样式的字符串，而非显示整列的内容。
  -i   忽略字符大小写的差别。
  -l   列出文件内容符合指定的范本样式的文件名称。
  -n   在显示符合坊本样式的那一列之前，标示出该列的列数编号。
  -N   不要递归处理。
  -r   递归处理，将指定目录下的所有文件及子目录一并处理。
  -R<范本样式>   此参数的效果和指定“-r”参数类似，但只主力符合范本样式文件名称的文件。
  -v   反转查找。
  -W<列长度>   限制符合范本样式的字符串所在列，必须拥有的字符数。
  -x<扩展名>   只处理符合指定扩展名的文件名称的文件。
  --help   在线帮助。
  --version   显示版本信息。
### sed(stream editor)：利用script来处理文本文件。
语　　法：sed [-hnV][-e<script>][-f<script文件>][文本文件]
补充说明：sed可依照script的指令，来处理、编辑文本文件。
参　　数：
  -e<script>或--expression=<script>   以选项中指定的script来处理输入的文本文件。
  -f<script文件>或--file=<script文件>  以选项中指定的script文件来处理输入的文本文件。
  -h或--help   显示帮助。
  -n或--quiet或--silent   仅显示script处理后的结果。
  -V或--version   显示版本信息。
### sort：将文本文件内容加以排序。
语　　法：sort [-bcdfimMnr][-o<输出文件>][-t<分隔字符>][+<起始栏位>-<结束栏位>][--help][--verison][文件]
补充说明：sort可针对文本文件的内容，以行为单位来排序。
参　　数：
  -b   忽略每行前面开始出的空格字符。
  -c   检查文件是否已经按照顺序排序。
  -d   排序时，处理英文字母、数字及空格字符外，忽略其他的字符。
  -f   排序时，将小写字母视为大写字母。
  -i   排序时，除了040至176之间的ASCII字符外，忽略其他的字符。
  -m   将几个排序好的文件进行合并。
  -M   将前面3个字母依照月份的缩写进行排序。
  -n   依照数值的大小排序。
  -o<输出文件>   将排序后的结果存入指定的文件。
  -r   以相反的顺序来排序。
  -t<分隔字符>   指定排序时所用的栏位分隔字符。
  +<起始栏位>-<结束栏位>   以指定的栏位来排序，范围由起始栏位到结束栏位的前一栏位。
  --help   显示帮助。
  --version   显示版本信息。
### spell：拼字检查程序。
语　　法：spell
补充说明：spell可从标准输入设备读取字符串，结束后显示拼错的词汇。
### tr(translate character)：转换字符。
语　　法：tr [-cdst][--help][--version][第一字符集][第二字符集]
补充说明：tr指令从标准输入设备读取数据，经过字符串转译后，输出到标准输出设备。
参　　数：
  -c或--complerment   取代所有不属于第一字符集的字符。
  -d或--delete   删除所有属于第一字符集的字符。
  -s或--squeeze-repeats   把连续重复的字符以单独一个字符表示。
  -t或--truncate-set1   先删除第一字符集较第二字符集多出的字符。
  --help   在线帮助。
  --version   显示版本信息。
### uniq：检查及删除文本文件中重复出现的行列。
语　　法：uniq [-cdu][-f<栏位>][-s<字符位置>][-w<字符位置>][--help][--version][输入文件][输出文件]
补充说明：uniq可检查文本文件中重复出现的行列。
参　　数：
  -c或--count   在每列旁边显示该行重复出现的次数。
  -d或--repeated   仅显示重复出现的行列。
  -f<栏位>或--skip-fields=<栏位>   忽略比较指定的栏位。
  -s<字符位置>或--skip-chars=<字符位置>   忽略比较指定的字符。
  -u或--unique   仅显示出一次的行列。
  -w<字符位置>或--check-chars=<字符位置>   指定要比较的字符。
  --help   显示帮助。
  --version   显示版本信息。
  [输入文件]   指定已排序好的文本文件。
  [输出文件]   指定输出的文件。
### wc(word count)：计算字数。
语　　法：wc [-clw][--help][--version][文件...]
补充说明：利用wc指令我们可以计算文件的Byte数、字数、或是列数，若不指定文件名称、或是所给予的文件名为“-”，则wc指令会从标准输入设备读取数据。
参　　数：
  -c或--bytes或--chars   只显示Bytes数。
  -l或--lines   只显示列数。
  -w或--words   只显示字数。
  --help   在线帮助。
  --version   显示版本信息。

 
4. linux压缩备份命令
### ar：建立或修改备存文件，或是从备存文件中抽取文件。
语　　法：ar[-dmpqrtx][cfosSuvV][a<成员文件>][b<成员文件>][i<成员文件>][备存文件][成员文件]
补充说明：ar可让您集合许多文件，成为单一的备存文件。在备存文件中，所有成员文件皆保有原来的属性与权限。
参　　数：
　指令参数 
　-d 　删除备存文件中的成员文件。 
　-m 　变更成员文件在备存文件中的次序。 
　-p 　显示备存文件中的成员文件内容。 
　-q 　将问家附加在备存文件末端。 
　-r 　将文件插入备存文件中。 
　-t 　显示备存文件中所包含的文件。 
　-x 　自备存文件中取出成员文件。 
　选项参数 
　a<成员文件> 　将文件插入备存文件中指定的成员文件之后。 
　b<成员文件> 　将文件插入备存文件中指定的成员文件之前。 
　c 　建立备存文件。 
　f 　为避免过长的文件名不兼容于其他系统的ar指令指令，因此可利用此参数，截掉要放入备存文件中过长的成员文件名称。 
　i<成员文件> 　将问家插入备存文件中指定的成员文件之前。 
　o 　保留备存文件中文件的日期。 
　s 　若备存文件中包含了对象模式，可利用此参数建立备存文件的符号表。 
　S 　不产生符号表。 
　u 　只将日期较新文件插入备存文件中。 
　v 　程序执行时显示详细的信息。 
　V 　显示版本信息。
### bunzip2：.bz2文件的解压缩程序。
语　　法：bunzip2 [-fkLsvV][.bz2压缩文件]
补充说明：bunzip2可解压缩.bz2格式的压缩文件。bunzip2实际上是bzip2的符号连接，执行bunzip2与bzip2 -d的效果相同。
参　　数：
　-f或--force 　解压缩时，若输出的文件与现有文件同名时，预设不会覆盖现有的文件。若要覆盖，请使用此参数。 
　-k或--keep 　在解压缩后，预设会删除原来的压缩文件。若要保留压缩文件，请使用此参数。 
　-s或--small 　降低程序执行时，内存的使用量。 
　-v或--verbose 　解压缩文件时，显示详细的信息。 
　-l,--license,-V或--version 　显示版本信息。
### bzip2：.bz2文件的压缩程序。
语　　法：bzip2 [-cdfhkLstvVz][--repetitive-best][--repetitive-fast][- 压缩等级][要压缩的文件]
补充说明：bzip2采用新的压缩演算法，压缩效果比传统的LZ77/LZ78压缩演算法来得好。若没有加上任何参数，bzip2压缩完文件后会产生.bz2的压缩文件，并删除原始的文件。
参　　数：
　-c或--stdout 　将压缩与解压缩的结果送到标准输出。 
　-d或--decompress 　执行解压缩。 
　-f或--force 　bzip2在压缩或解压缩时，若输出文件与现有文件同名，预设不会覆盖现有文件。若要覆盖，请使用此参数。 
　-h或--help 　显示帮助。 
　-k或--keep 　bzip2在压缩或解压缩后，会删除原始的文件。若要保留原始文件，请使用此参数。 
　-s或--small 　降低程序执行时内存的使用量。 
　-t或--test 　测试.bz2压缩文件的完整性。 
　-v或--verbose 　压缩或解压缩文件时，显示详细的信息。 
　-z或--compress 　强制执行压缩。 
　-L,--license,
　-V或--version 　显示版本信息。 
　--repetitive-best 　若文件中有重复出现的资料时，可利用此参数提高压缩效果。 
　--repetitive-fast 　若文件中有重复出现的资料时，可利用此参数加快执行速度。 
　-压缩等级 　压缩时的区块大小。
### bzip2recover：用来修复损坏的.bz2文件。
语　　法：bzip2recover [.bz2 压缩文件]
补充说明：bzip2是以区块的方式来压缩文件，每个区块视为独立的单位。因此，当某一区块损坏时，便可利用bzip2recover，试着将文件中的区块隔开来，以便解压缩正常的区块。通常只适用在压缩文件很大的情况。
### compress：压缩或解压文件。
语　　法：compress [-cdfrvV][-b <压缩效率>][文件或目录...]
补充说明：compress是个历史悠久的压缩程序，文件经它压缩后，其名称后面会多出".Z"的扩展名。当要解压缩时，可执行uncompress指令。事实上uncompress是指向compress的符号连接，因此不论是压缩或解压缩，都可通过compress指令单独完成。
参　　数：
　-b<压缩效率> 　压缩效率是一个介于9－16的数值，预设值为"16"，指定愈大的数值，压缩效率就愈高。 
　-c 　把压缩后的文件输出到标准输出设备，不去更动原始文件。 
　-d 　对文件进行解压缩而非压缩。 
　-f 　强制保存压缩文件，不理会文件名称或硬连接是否存在，该文件是否为符号连接以及压缩效率高低的问题。 
　-r 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v 　显示指令执行过程。 
　-V 　显示指令版本及程序预设值。
### cpio(copy in/out)：备份文件。
语　　法：cpio [-0aABckLovV][-C <输入/输出大小>][-F <备份档>][-H <备份格式>][-O <备份档>][--block-size=<区块大小>][--force-local][--help][--quiet][--version] 或 cpio [-bBcdfikmnrsStuvV][-C <输入/输出大小>][-E <范本文件>][-F <备份档>][-H <备份格式>][-I <备份档>][-M <回传信息>][-R <拥有者><:/.><所属群组>][--block-size=<区块大小>][--force-local][--help][--no-absolute-filenames][--no-preserve-owner][--only-verify-crc][--quiet][--sparse][--version][范本样式...] 或 cpio [-0adkiLmpuvV][-R <拥有者><:/.><所属群组>][--help][--no-preserve-owner][--quiet][--sparse][--version][目的目]
补充说明：cpio是用来建立，还原备份档的工具程序，它可以加入，解开cpio或tra备份档内的文件。
参　　数：
　-0或--null 　接受新增列控制字符，通常配合find指令的"-print0"参数使用。 
　-a或--reset-access-time 　重新设置文件的存取时间。 
　-A或--append 　附加到已存在的备份档中，且这个备份档必须存放在磁盘上，而不能放置于磁带机里。 
　-b或--swap 　此参数的效果和同时指定"-sS"参数相同。 
　-B 　将输入/输出的区块大小改成5210 Bytes。 
　-c 　使用旧ASCII备份格式。 
　-C<区块大小>或--io-size=<区块大小> 　设置输入/输出的区块大小，单位是Byte。 
　-d或--make-directories 　如有需要cpio会自行建立目录。 
　-E<范本文件>或--pattern-file=<范本文件> 　指定范本文件，其内含有一个或多个范本样式，让cpio解开符合范本条件的文件，格式为每列一个范本样式。 
　-f或--nonmatching 　让cpio解开所有不符合范本条件的文件。 
　-F<备份档>或--file=<备份档> 　指定备份档的名称，用来取代标准输入或输出，也能借此通过网络使用另一台主机的保存设备存取备份档。 
　-H<备份格式> 　指定备份时欲使用的文件格式。 
　-i或--extract 　执行copy-in模式，还原备份档。 
　-l<备份档> 　指定备份档的名称，用来取代标准输入，也能借此通过网络使用另一台主机的保存设备读取备份档。 
　-k 　此参数将忽略不予处理，仅负责解决cpio不同版本间的兼容性问题。 
　-l或--link 　以硬连接的方式取代复制文件，可在copy-pass模式下运用。 
　-L或--dereference 　不建立符号连接，直接复制该连接所指向的原始文件。 
　-m或preserve-modification-time 　不去更换文件的更改时间。 
　-M<回传信息>或--message=<回传信息> 　设置更换保存媒体的信息。 
　-n或--numeric-uid-gid 　使用"-tv"参数列出备份档的内容时，若再加上参数"-n"，则会以用户识别码和群组识别码替代拥有者和群组名称列出文件清单。 
　-o或--create 　执行copy-out模式，建立备份档。 
　-O<备份档> 　指定备份档的名称，用来取代标准输出，也能借此通过网络　使用另一台主机的保存设备存放备份档。 
　-p或--pass-through 　执行copy-pass模式，略过备份步骤，直接将文件复制到目的目录。 
　-r或--rename 　当有文件名称需要更动时，采用互动模式。 
　-R<拥有者><:/.><所属群组>或
　----owner<拥有者><:/.><所属群组> 　在copy-in模式还原备份档，或copy-pass模式复制文件时，可指定这些备份，复制的文件的拥有者与所属群组。 
　-s或--swap-bytes 　交换每对字节的内容。 
　-S或--swap-halfwords 　交换每半个字节的内容。 
　-t或--list 　将输入的内容呈现出来。 
　-u或--unconditional 　置换所有文件，不论日期时间的新旧与否，皆不予询问而直接覆盖。 
　-v或--verbose 　详细显示指令的执行过程。 
　-V或--dot 　执行指令时，在每个文件的执行程序前面加上"."号 
　--block-size=<区块大小> 　设置输入/输出的区块大小，假如设置数值为5，则区块大小为2500，若设置成10，则区块大小为5120，依次类推。 
　--force-local 　强制将备份档存放在本地主机。 
　--help 　在线帮助。 
　--no-absolute-filenames 　使用相对路径建立文件名称。 
　--no-preserve-owner 　不保留文件的拥有者，谁解开了备份档，那些文件就归谁所有。 
　-only-verify-crc 　当备份档采用CRC备份格式时，可使用这项参数检查备份档内的每个文件是否正确无误。 
　--quiet 　不显示复制了多少区块。 
　--sparse 　倘若一个文件内含大量的连续0字节，则将此文件存成稀疏文件。 
　--version 　显示版本信息。
### dump：备份文件系统。
语　　法：dump [-cnu][-0123456789][-b <区块大小>][-B <区块数目>][-d <密度>][-f <设备名称>][-h <层级>][-s <磁带长度>][-T <日期>][目录或文件系统] 或 dump [-wW]
补充说明：dump为备份工具程序，可将目录或整个文件系统备份至指定的设备，或备份成一个大文件。
参　　数：
　-0123456789 　备份的层级。 
　-b<区块大小> 　指定区块的大小，单位为KB。 
　-B<区块数目> 　指定备份卷册的区块数目。 
　-c 　修改备份磁带预设的密度与容量。 
　-d<密度> 　设置磁带的密度。单位为BPI。 
　-f<设备名称> 　指定备份设备。 
　-h<层级> 　当备份层级等于或大雨指定的层级时，将不备份用户标示为"nodump"的文件。 
　-n 　当备份工作需要管理员介入时，向所有"operator"群组中的使用者发出通知。 
　-s<磁带长度> 　备份磁带的长度，单位为英尺。 
　-T<日期> 　指定开始备份的时间与日期。 
　-u 　备份完毕后，在/etc/dumpdates中记录备份的文件系统，层级，日期与时间等。 
　-w 　与-W类似，但仅显示需要备份的文件。 
　-W 　显示需要备份的文件及其最后一次备份的层级，时间与日期。
### gunzip(gnu unzip)：解压文件。
语　　法：gunzip [-acfhlLnNqrtvV][-s <压缩字尾字符串>][文件...] 或 gunzip [-acfhlLnNqrtvV][-s <压缩字尾字符串>][目录]
补充说明：gunzip是个使用广泛的解压缩程序，它用于解开被gzip压缩过的文件，这些压缩文件预设最后的扩展名为".gz"。事实上gunzip就是gzip的硬连接，因此不论是压缩或解压缩，都可通过gzip指令单独完成。
参　　数：
　-a或--ascii 　使用ASCII文字模式。 
　-c或--stdout或--to-stdout 　把解压后的文件输出到标准输出设备。 
　-f或-force 　强行解开压缩文件，不理会文件名称或硬连接是否存在以及该文件是否为符号连接。 
　-h或--help 　在线帮助。 
　-l或--list 　列出压缩文件的相关信息。 
　-L或--license 　显示版本与版权信息。 
　-n或--no-name 　解压缩时，若压缩文件内含有远来的文件名称及时间戳记，则将其忽略不予处理。 
　-N或--name 　解压缩时，若压缩文件内含有原来的文件名称及时间戳记，则将其回存到解开的文件上。 
　-q或--quiet 　不显示警告信息。 
　-r或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-S<压缩字尾字符串>或--suffix<压缩字尾字符串> 　更改压缩字尾字符串。 
　-t或--test 　测试压缩文件是否正确无误。 
　-v或--verbose 　显示指令执行过程。 
　-V或--version 显示版本信息。
### gzexe(gzip executable)：压缩执行文件。
语　　法：gzexe [-d][执行文件...]
补充说明：gzexe是用来压缩执行文件的程序。当您去执行被压缩过的执行文件时，该文件会自动解压然后继续执行，和使用一般的执行文件相同。
参　　数：
　-d 　解开压缩文件。
### gzip(gnu zip)：压缩文件。
语　　法：gzip [-acdfhlLnNqrtvV][-S <压缩字尾字符串>][-<压缩效率>][--best/fast][文件...] 或 gzip [-acdfhlLnNqrtvV][-S <压缩字尾字符串>][-<压缩效率>][--best/fast][目录]
补充说明：gzip是个使用广泛的压缩程序，文件经它压缩过后，其名称后面会多出".gz"的扩展名。
参　　数：
　-a或--ascii 　使用ASCII文字模式。 
　-c或--stdout或--to-stdout 　把压缩后的文件输出到标准输出设备，不去更动原始文件。 
　-d或--decompress或----uncompress 　解开压缩文件。 
　-f或--force 　强行压缩文件。不理会文件名称或硬连接是否存在以及该文件是否为符号连接。 
　-h或--help 　在线帮助。 
　-l或--list 　列出压缩文件的相关信息。 
　-L或--license 　显示版本与版权信息。 
　-n或--no-name 　压缩文件时，不保存原来的文件名称及时间戳记。 
　-N或--name 　压缩文件时，保存原来的文件名称及时间戳记。 
　-q或--quiet 　不显示警告信息。 
　-r或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-S<压缩字尾字符串>或----suffix<压缩字尾字符串> 　更改压缩字尾字符串。 
　-t或--test 　测试压缩文件是否正确无误。 
　-v或--verbose 　显示指令执行过程。 
　-V或--version 　显示版本信息。 
　-<压缩效率> 　压缩效率是一个介于1－9的数值，预设值为"6"，指定愈大的数值，压缩效率就会愈高。 
　--best 　此参数的效果和指定"-9"参数相同。 
　--fast 　此参数的效果和指定"-1"参数相同。
### lha：压缩或解压缩文件。
语　　法：
lha [-acdfglmnpqtuvx][-a <0/1/2>/u</0/1/2>][-<a/c/u>d][-<e/x>i][-<a/u>o][-<e/x>w=<目的目录>][-<a/u>z][压缩文件][文件...] 或 lha [-acdfglmnpqtuvx][-a <0/1/2>/u</0/1/2>][-<a/c/u>d][-<e/x>i][-<a/u>o][-<e/x>w=<目的目录>][-<a/u>z][压缩文件][目录...]
补充说明：lha是从lharc演变而来的压缩程序，文件经它压缩后，会另外产生具有".lzh"扩展名的压缩文件。
参　　数：
  -a或a   压缩文件，并加入到压缩文件内。
  -a<0/1/2>/u</0/1/2>   压缩文件时，采用不同的文件头。
  -c或c   压缩文件，重新建构新的压缩文件后，再将其加入。
  -d或d   从压缩文件内删除指定的文件。
  -<a/c/u>d或<a/c/u>d   压缩文件，然后将其加入，重新建构，更新压缩文件或，删除原始文件，也就是把文件移到压缩文件中。
  -e或e   解开压缩文件。
  -f或f   强制执行lha命令，在解压时会直接覆盖已有的文件而不加以询问。
  -g或g   使用通用的压缩格式，便于解决兼容性的问题。
  -<e/x>i或<e/x>i   解开压缩文件时，忽略保存在压缩文件内的文件路径，直接将其解压后存放在现行目录下或是指定的目录中。
  -l或l   列出压缩文件的相关信息。
  -m或m   此参数的效果和同时指定"-ad"参数相同。
  -n或n   不执行指令，仅列出实际执行会进行的动作。
  -<a/u>o或<a/u>o   采用lharc兼容格式，将压缩后的文件加入，更新压缩文件。
  -p或p   从压缩文件内输出到标准输出设备。
  -q或q   不显示指令执行过程。
  -t或t   检查备份文件内的每个文件是否正确无误。
  -u或u   更换较新的文件到压缩文件内。
  -u</0/1/2>或u</0/1/2>   在文件压缩时采用不同的文件头，然后更新到压缩文件内。
  -v或v   详细列出压缩文件的相关信息。
  -<e/x>w=<目的目录>或<e/x>w=<目的目录>   指定解压缩的目录。
  -x或x   解开压缩文件。
  -<a/u>z或<a/u>z   不压缩文件，直接把它加入，更新压缩文件。
### restore：还原(Restore)由倾倒(Dump)操作所备份下来的文件或整个文件系统(一个分区)。
语　　法：restore [-cCvy][-b <区块大小>][-D <文件系统>][-f <备份文件>][-s <文件编号>] 或 restore [-chimvy][-b <区块大小>][-f <备份文件>][-s <文件编号>] 或 restore [-crvy][-b <区块大小>][-f <备份文件>][-s <文件编号>] 或 restore [-cRvy][-b <区块大小>][-D <文件系统>][-f <备份文件>][-s <文件编号>] 或 restore [chtvy][-b <区块大小>][-D <文件系统>][-f <备份文件>][-s <文件编号>][文件...] 或 restore [-chmvxy][-b <区块大小>][-D <文件系统>][-f <备份文件>][-s <文件编号>][文件...]
补充说明：restore 指令所进行的操作和dump指令相反，倾倒操作可用来备份文件，而还原操作则是写回这些已备份的文件。
参　　数：
  -b<区块大小>   设置区块大小，单位是Byte。
  -c   不检查倾倒操作的备份格式，仅准许读取使用旧格式的备份文件。
  -C   使用对比模式，将备份的文件与现行的文件相互对比。
  -D<文件系统>   允许用户指定文件系统的名称。
  -f<备份文件>   从指定的文件中读取备份数据，进行还原操作。
  -h   仅解出目录而不包括与该目录相关的所有文件。
  -i   使用互动模式，在进行还原操作时，restore指令将依序询问用户。
  -m   解开符合指定的inode编号的文件或目录而非采用文件名称指定。
  -r   进行还原操作。
  -R   全面还原文件系统时，检查应从何处开始进行。
  -s<文件编号>   当备份数据超过一卷磁带时，您可以指定备份文件的编号。
  -t   指定文件名称，若该文件已存在备份文件中，则列出它们的名称。
  -v   显示指令执行过程。
  -x   设置文件名称，且从指定的存储媒体里读入它们，若该文件已存在在备份文件中，则将其还原到文件系统内。
  -y   不询问任何问题，一律以同意回答并继续执行指令。
### tar(tape archive)：备份文件。
语　　法：tar [-ABcdgGhiklmMoOpPrRsStuUvwWxzZ][-b <区块数目>][-C <目的目录>][-f <备份文件>][-F <Script文件>][-K <文件>][-L <媒体容量>][-N <日期时间>][-T <范本文件>][-V <卷册名称>][-X <范本文件>][-<设备编号><存储密度>][--after-date=<日期时间>][--atime-preserve][--backuup=<备份方式>][--checkpoint][--concatenate][--confirmation][--delete][--exclude=<范本样式>][--force-local][--group=<群组名称>][--help][--ignore-failed-read][--new-volume-script=<Script文件>][--newer-mtime][--no-recursion][--null][--numeric-owner][--owner=<用户名称>][--posix][--erve][--preserve-order][--preserve-permissions][--record-size=<区块数目>][--recursive-unlink][--remove-files][--rsh-command=<执行指令>][--same-owner][--suffix=<备份字尾字符串>][--totals][--use-compress-program=<执行指令>][--version][--volno-file=<编号文件>][文件或目录...]
补充说明：tar是用来建立，还原备份文件的工具程序，它可以加入，解开备份文件内的文件。
参　　数：
  -A或--catenate   新增温暖件到已存在的备份文件。
  -b<区块数目>或--blocking-factor=<区块数目>   设置每笔记录的区块数目，每个区块大小为12Bytes。
  -B或--read-full-records   读取数据时重设区块大小。
  -c或--create   建立新的备份文件。
  -C<目的目录>或--directory=<目的目录>   切换到指定的目录。
  -d或--diff或--compare   对比备份文件内和文件系统上的文件的差异。
  -f<备份文件>或--file=<备份文件>   指定备份文件。
  -F<Script文件>或--info-script=<Script文件>   每次更换磁带时，就执行指定的Script文件。
  -g或--listed-incremental   处理GNU格式的大量备份。
  -G或--incremental   处理旧的GNU格式的大量备份。
  -h或--dereference   不建立符号连接，直接复制该连接所指向的原始文件。
  -i或--ignore-zeros   忽略备份文件中的0 Byte区块，也就是EOF。
  -k或--keep-old-files   解开备份文件时，不覆盖已有的文件。
  -K<文件>或--starting-file=<文件>   从指定的文件开始还原。
  -l或--one-file-system   复制的文件或目录存放的文件系统，必须与tar指令执行时所处的文件系统相同，否则不予复制。
  -L<媒体容量>或-tape-length=<媒体容量>   设置存放每体的容量，单位以1024 Bytes计算。
  -m或--modification-time   还原文件时，不变更文件的更改时间。
  -M或--multi-volume   在建立，还原备份文件或列出其中的内容时，采用多卷册模式。
  -N<日期格式>或--newer=<日期时间>   只将较指定日期更新的文件保存到备份文件里。
  -o或--old-archive或--portability   将资料写入备份文件时使用V7格式。
  -O或--stdout   把从备份文件里还原的文件输出到标准输出设备。
  -p或--same-permissions   用原来的文件权限还原文件。
  -P或--absolute-names   文件名使用绝对名称，不移除文件名称前的"/"号。
  -r或--append   新增文件到已存在的备份文件的结尾部分。
  -R或--block-number   列出每个信息在备份文件中的区块编号。
  -s或--same-order   还原文件的顺序和备份文件内的存放顺序相同。
  -S或--sparse   倘若一个文件内含大量的连续0字节，则将此文件存成稀疏文件。
  -t或--list   列出备份文件的内容。
  -T<范本文件>或--files-from=<范本文件>   指定范本文件，其内含有一个或多个范本样式，让tar解开或建立符合设置条件的文件。
  -u或--update   仅置换较备份文件内的文件更新的文件。
  -U或--unlink-first   解开压缩文件还原文件之前，先解除文件的连接。
  -v或--verbose   显示指令执行过程。
  -V<卷册名称>或--label=<卷册名称>   建立使用指定的卷册名称的备份文件。
  -w或--interactive   遭遇问题时先询问用户。
  -W或--verify   写入备份文件后，确认文件正确无误。
  -x或--extract或--get  从备份文件中还原文件。
  -X<范本文件>或--exclude-from=<范本文件>  指定范本文件，其内含有一个或多个范本样式，让ar排除符合设置条件的文件。
  -z或--gzip或--ungzip   通过gzip指令处理备份文件。
  -Z或--compress或--uncompress   通过compress指令处理备份文件。
  -<设备编号><存储密度>   设置备份用的外围设备编号及存放数据的密度。
  --after-date=<日期时间>   此参数的效果和指定"-N"参数相同。
  --atime-preserve   不变更文件的存取时间。
  --backup=<备份方式>或--backup   移除文件前先进行备份。
  --checkpoint   读取备份文件时列出目录名称。
  --concatenate   此参数的效果和指定"-A"参数相同。
  --confirmation   此参数的效果和指定"-w"参数相同。
  --delete   从备份文件中删除指定的文件。
  --exclude=<范本样式>   排除符合范本样式的问家。
  --group=<群组名称>   把加入设备文件中的文件的所属群组设成指定的群组。
  --help   在线帮助。
  --ignore-failed-read   忽略数据读取错误，不中断程序的执行。
  --new-volume-script=<Script文件>   此参数的效果和指定"-F"参数相同。
  --newer-mtime   只保存更改过的文件。
  --no-recursion   不做递归处理，也就是指定目录下的所有文件及子目录不予处理。
  --null   从null设备读取文件名称。
  --numeric-owner   以用户识别码及群组识别码取代用户名称和群组名称。
  --owner=<用户名称>   把加入备份文件中的文件的拥有者设成指定的用户。
  --posix   将数据写入备份文件时使用POSIX格式。
  --preserve      此参数的效果和指定"-ps"参数相同。
  --preserve-order      此参数的效果和指定"-A"参数相同。
  --preserve-permissions      此参数的效果和指定"-p"参数相同。
  --record-size=<区块数目>      此参数的效果和指定"-b"参数相同。
  --recursive-unlink   解开压缩文件还原目录之前，先解除整个目录下所有文件的连接。
  --remove-files   文件加入备份文件后，就将其删除。
  --rsh-command=<执行指令>   设置要在远端主机上执行的指令，以取代rsh指令。
  --same-owner   尝试以相同的文件拥有者还原问家你。
  --suffix=<备份字尾字符串>   移除文件前先行备份。
  --totals   备份文件建立后，列出文件大小。
  --use-compress-program=<执行指令>   通过指定的指令处理备份文件。
  --version   显示版本信息。
  --volno-file=<编号文件>   使用指定文件内的编号取代预设的卷册编号。
### unarj：解压缩.arj文件。
语　　法：unarj [eltx][.arj压缩文件]
补充说明：unarj为.arj压缩文件的压缩程序。
参　　数：
  e   解压缩.arj文件。
  l   显示压缩文件内所包含的文件。
  t   检查压缩文件是否正确。
  x   解压缩时保留原有的路径。
### unzip：解压缩zip文件
语　　法：unzip [-cflptuvz][-agCjLMnoqsVX][-P <密码>][.zip文件][文件][-d <目录>][-x <文件>] 或 unzip [-Z]
补充说明：unzip为.zip压缩文件的解压缩程序。
参　　数：
  -c   将解压缩的结果显示到屏幕上，并对字符做适当的转换。
  -f   更新现有的文件。
  -l   显示压缩文件内所包含的文件。
  -p   与-c参数类似，会将解压缩的结果显示到屏幕上，但不会执行任何的转换。
  -t   检查压缩文件是否正确。
  -u   与-f参数类似，但是除了更新现有的文件外，也会将压缩文件中的其他文件解压缩到目录中。
  -v   执行是时显示详细的信息。
  -z   仅显示压缩文件的备注文字。
  -a   对文本文件进行必要的字符转换。
  -b   不要对文本文件进行字符转换。 
  -C   压缩文件中的文件名称区分大小写。
  -j   不处理压缩文件中原有的目录路径。
  -L   将压缩文件中的全部文件名改为小写。
  -M   将输出结果送到more程序处理。
  -n   解压缩时不要覆盖原有的文件。
  -o   不必先询问用户，unzip执行后覆盖原有文件。
  -P<密码>   使用zip的密码选项。
  -q   执行时不显示任何信息。
  -s   将文件名中的空白字符转换为底线字符。
  -V   保留VMS的文件版本信息。
  -X   解压缩时同时回存文件原来的UID/GID。
  [.zip文件]   指定.zip压缩文件。
  [文件]   指定要处理.zip压缩文件中的哪些文件。
  -d<目录>   指定文件解压缩后所要存储的目录。
  -x<文件>   指定不要处理.zip压缩文件中的哪些文件。
  -Z   unzip -Z等于执行zipinfo指令。
### zip：压缩文件。
语　　法：zip [-AcdDfFghjJKlLmoqrSTuvVwXyz$][-b <工作目录>][-ll][-n <字尾字符串>][-t <日期时间>][-<压缩效率>][压缩文件][文件...][-i <范本样式>][-x <范本样式>]
补充说明：zip是个使用广泛的压缩程序，文件经它压缩后会另外产生具有".zip"扩展名的压缩文件。
参　　数：
  -A   调整可执行的自动解压缩文件。
  -b<工作目录>   指定暂时存放文件的目录。
  -c   替每个被压缩的文件加上注释。
  -d   从压缩文件内删除指定的文件。
  -D   压缩文件内不建立目录名称。
  -f   此参数的效果和指定"-u"参数类似，但不仅更新既有文件，如果某些文件原本不存在于压缩文件内，使用本参数会一并将其加入压缩文件中。
  -F   尝试修复已损坏的压缩文件。
  -g   将文件压缩后附加在既有的压缩文件之后，而非另行建立新的压缩文件。
  -h   在线帮助。
  -i<范本样式>   只压缩符合条件的文件。
  -j   只保存文件名称及其内容，而不存放任何目录名称。
  -J   删除压缩文件前面不必要的数据。
  -k   使用MS-DOS兼容格式的文件名称。
  -l   压缩文件时，把LF字符置换成LF+CR字符。
  -ll   压缩文件时，把LF+CR字符置换成LF字符。
  -L   显示版权信息。
  -m   将文件压缩并加入压缩文件后，删除原始文件，即把文件移到压缩文件中。
  -n<字尾字符串>   不压缩具有特定字尾字符串的文件。
  -o   以压缩文件内拥有最新更改时间的文件为准，将压缩文件的更改时间设成和该文件相同。
  -q   不显示指令执行过程。
  -r   递归处理，将指定目录下的所有文件和子目录一并处理。
  -S   包含系统和隐藏文件。
  -t<日期时间>   把压缩文件的日期设成指定的日期。
  -T   检查备份文件内的每个文件是否正确无误。
  -u   更换较新的文件到压缩文件内。
  -v   显示指令执行过程或显示版本信息。
  -V   保存VMS操作系统的文件属性。
  -w   在文件名称里假如版本编号，本参数仅在VMS操作系统下有效。
  -x<范本样式>   压缩时排除符合条件的文件。
  -X   不保存额外的文件属性。
  -y   直接保存符号连接，而非该连接所指向的文件，本参数仅在UNIX之类的系统下有效。
  -z   替压缩文件加上注释。
  -$   保存第一个被压缩文件所在磁盘的卷册名称。
  -<压缩效率>   压缩效率是一个介于1-9的数值。
### zipinfo：列出压缩文件信息。
语　　法：zipinfo [-12hlmMstTvz][压缩文件][文件...][-x <范本样式>]
补充说明：执行zipinfo指令可得知zip压缩文件的详细信息。
参　　数：
  -1   只列出文件名称。
  -2   此参数的效果和指定"-1"参数类似，但可搭配"-h","-t"和"-z"参数使用。
  -h   只列出压缩文件的文件名称。
  -l   此参数的效果和指定"-m"参数类似，但会列出原始文件的大小而非每个文件的压缩率。
  -m   此参数的效果和指定"-s"参数类似，但多会列出每个文件的压缩率。
  -M   若信息内容超过一个画面，则采用类似more指令的方式列出信息。
  -s   用类似执行"ls -l"指令的效果列出压缩文件内容。
  -t   只列出压缩文件内所包含的文件数目，压缩前后的文件大小及压缩率。
  -T   将压缩文件内每个文件的日期时间用年，月，日，时，分，秒的顺序列出。
  -v   详细显示压缩文件内每一个文件的信息。
  -x<范本样式>   不列出符合条件的文件的信息。
  -z   如果压缩文件内含有注释，就将注释显示出来。



 
## linux文件管理命令
### diff(differential)：比较文件的差异。 
语　　法：diff [-abBcdefHilnNpPqrstTuvwy][-<行数>][-C <行数>][-D <巨集名称>][-I <字符或字符串>][-S <文件>][-W <宽度>][-x <文件或目录>][-X <文件>][--help][--left-column][--suppress-common-line][文件或目录1][文件或目录2] 
补充说明：diff以逐行的方式，比较文本文件的异同处。所是指定要比较目录，则diff会比较目录中相同文件名的文件，但不会比较其中子目录。
参　　数：
　-<行数> 　指定要显示多少行的文本。此参数必须与-c或-u参数一并使用。 
　-a或--text 　diff预设只会逐行比较文本文件。
　-b或--ignore-space-change 　不检查空格字符的不同。 
　-B或--ignore-blank-lines 　不检查空白行。 
　-c 　显示全部内文，并标出不同之处。 
　-C<行数>或--context<行数> 　与执行"-c-<行数>"指令相同。 
　-d或--minimal 　使用不同的演算法，以较小的单位来做比较。 
　-D<巨集名称>或ifdef<巨集名称> 　此参数的输出格式可用于前置处理器巨集。 
　-e或--ed 　此参数的输出格式可用于ed的script文件。 
　-f或-forward-ed 　输出的格式类似ed的script文件，但按照原来文件的顺序来显示不同处。 
　-H或--speed-large-files 　比较大文件时，可加快速度。 
　-l<字符或字符串>或--ignore-matching-lines<字符或字符串> 　若两个文件在某几行有所不同，而这几行同时都包含了选项中指定的字符或字符串，则不显示这两个文件的差异。 
　-i或--ignore-case 　不检查大小写的不同。 
　-l或--paginate 　将结果交由pr程序来分页。 
　-n或--rcs 　将比较结果以RCS的格式来显示。 
　-N或--new-file 　在比较目录时，若文件A仅出现在某个目录中，预设会显示：
Only in目录：文件A若使用-N参数，则diff会将文件A与一个空白的文件比较。  
　-p 　若比较的文件为C语言的程序码文件时，显示差异所在的函数名称。 
　-P或--unidirectional-new-file 　与-N类似，但只有当第二个目录包含了一个第一个目录所没有的文件时，才会将这个文件与空白的文件做比较。 
　-q或--brief 　仅显示有无差异，不显示详细的信息。 
　-r或--recursive 　比较子目录中的文件。 
　-s或--report-identical-files 　若没有发现任何差异，仍然显示信息。 
　-S<文件>或--starting-file<文件> 　在比较目录时，从指定的文件开始比较。 
　-t或--expand-tabs 　在输出时，将tab字符展开。 
　-T或--initial-tab 　在每行前面加上tab字符以便对齐。 
　-u,-U<列数>或--unified=<列数> 　以合并的方式来显示文件内容的不同。 
　-v或--version 　显示版本信息。 
　-w或--ignore-all-space 　忽略全部的空格字符。 
　-W<宽度>或--width<宽度> 　在使用-y参数时，指定栏宽。 
　-x<文件名或目录>或--exclude<文件名或目录> 　不比较选项中所指定的文件或目录。 
　-X<文件>或--exclude-from<文件> 　您可以将文件或目录类型存成文本文件，然后在=<文件>中指定此文本文件。 
　-y或--side-by-side 　以并列的方式显示文件的异同之处。 
　--help 　显示帮助。 
　--left-column 　在使用-y参数时，若两个文件某一行内容相同，则仅在左侧的栏位显示该行内容。
　--suppress-common-lines 　在使用-y参数时，仅显示不同之处。
### diffstat(differential status)：根据diff的比较结果，显示统计数字。
语　　法：diff [-wV][-n <文件名长度>][-p <文件名长度>]
补充说明：diffstat读取diff的输出结果，然后统计各文件的插入，删除，修改等差异计量。
参　　数：
　-n<文件名长度> 　指定文件名长度，指定的长度必须大于或等于所有文件中最长的文件名。 
　-p<文件名长度> 　与-n参数相同，但此处的<文件名长度>包括了文件的路径。 
　-w 　指定输出时栏位的宽度。 
　-V 　显示版本信息。
### file：辨识文件类型。 
语　　法：file [-beLvz][-f <名称文件>][-m <魔法数字文件>...][文件或目录...]  
补充说明：通过file指令，我们得以辨识该文件的类型。 
参　　数：
　-b 　列出辨识结果时，不显示文件名称。 
　-c 　详细显示指令执行过程，便于排错或分析程序执行的情形。 
　-f<名称文件> 　指定名称文件，其内容有一个或多个文件名称呢感，让file依序辨识这些文件，格式为每列一个文件名称。  
　-L 　直接显示符号连接所指向的文件的类别。 
　-m<魔法数字文件> 　指定魔法数字文件。 
　-v 　显示版本信息。 
　-z 　尝试去解读压缩文件的内容。
### find：查找文件或目录。
语　　法：find [目录...][-amin <分钟>][-anewer <参考文件或目录>][-atime <24小时数>][-cmin <分钟>][-cnewer <参考文件或目录>][-ctime <24小时数>][-daystart][-depyh][-empty][-exec <执行指令>][-false][-fls <列表文件>][-follow][-fprint <列表文件>][-fprint0 <列表文件>][-fprintf <列表文件><输出格式>][-fstype <文件系统类型>][-gid <群组识别码>][-group <群组名称>][-help][-ilname <范本样式>][-iname <范本样式>][-inum <inode编号>][-ipath <范本样式>][-iregex <范本样式>][-links <连接数目>][-lname <范本样式>][-ls][-maxdepth <目录层级>][-mindepth <目录层级>][-mmin <分钟>][-mount] 
[-mtime <24小时数>][-name <范本样式>][-newer <参考文件或目录>][-nogroup][noleaf] [-nouser][-ok <执行指令>][-path <范本样式>][-perm <权限数值>][-print][-print0][-printf <输出格式>][-prune][-regex <范本样式>][-size <文件大小>][-true][-type <文件类型>][-uid <用户识别码>][-used <日数>][-user <拥有者名称>][-version][-xdev][-xtype <文件类型>]
补充说明：find指令用于查找符合条件的文件。任何位于参数之前的字符串都将被视为欲查找的目录。
参　　数：
　-amin<分钟> 　查找在指定时间曾被存取过的文件或目录，单位以分钟计算。 
　-anewer<参考文件或目录> 　查找其存取时间较指定文件或目录的存取时间更接近现在的文件或目录。 
　-atime<24小时数> 　查找在指定时间曾被存取过的文件或目录，单位以24小时计算。 
　-cmin<分钟> 　查找在指定时间之时被更改的文件或目录。 
　-cnewer<参考文件或目录> 　查找其更改时间较指定文件或目录的更改时间更接近现在的文件或目录。 
　-ctime<24小时数> 　查找在指定时间之时被更改的文件或目录，单位以24小时计算。 
　-daystart 　从本日开始计算时间。 
　-depth 　从指定目录下最深层的子目录开始查找。 
　-expty 　寻找文件大小为0 Byte的文件，或目录下没有任何子目录或文件的空目录。 
　-exec<执行指令> 　假设find指令的回传值为True，就执行该指令。 
　-false 　将find指令的回传值皆设为False。 
　-fls<列表文件> 　此参数的效果和指定"-ls"参数类似，但会把结果保存为指定的列表文件。 
　-follow 　排除符号连接。 
　-fprint<列表文件> 　此参数的效果和指定"-print"参数类似，但会把结果保存成指定的列表文件。 
　-fprint0<列表文件> 　此参数的效果和指定"-print0"参数类似，但会把结果保存成指定的列表文件。 
　-fprintf<列表文件><输出格式> 　此参数的效果和指定"-printf"参数类似，但会把结果保存成指定的列表文件。 
　-fstype<文件系统类型> 　只寻找该文件系统类型下的文件或目录。 
　-gid<群组识别码> 　查找符合指定之群组识别码的文件或目录。 
　-group<群组名称> 　查找符合指定之群组名称的文件或目录。 
　-help或--help 　在线帮助。 
　-ilname<范本样式> 　此参数的效果和指定"-lname"参数类似，但忽略字符大小写的差别。 
　-iname<范本样式> 　此参数的效果和指定"-name"参数类似，但忽略字符大小写的差别。 
　-inum<inode编号> 　查找符合指定的inode编号的文件或目录。 
　-ipath<范本样式> 　此参数的效果和指定"-ipath"参数类似，但忽略字符大小写的差别。 
　-iregex<范本样式> 　此参数的效果和指定"-regexe"参数类似，但忽略字符大小写的差别。 
　-links<连接数目> 　查找符合指定的硬连接数目的文件或目录。 
　-iname<范本样式> 　指定字符串作为寻找符号连接的范本样式。 
　-ls 　假设find指令的回传值为True，就将文件或目录名称列出到标准输出。 
　-maxdepth<目录层级> 　设置最大目录层级。 
　-mindepth<目录层级> 　设置最小目录层级。 
　-mmin<分钟> 　查找在指定时间曾被更改过的文件或目录，单位以分钟计算。 
　-mount 　此参数的效果和指定"-xdev"相同。 
　-mtime<24小时数> 　查找在指定时间曾被更改过的文件或目录，单位以24小时计算。 
　-name<范本样式> 　指定字符串作为寻找文件或目录的范本样式。 
　-newer<参考文件或目录> 　查找其更改时间较指定文件或目录的更改时间更接近现在的文件或目录。 
　-nogroup 　找出不属于本地主机群组识别码的文件或目录。 
　-noleaf 　不去考虑目录至少需拥有两个硬连接存在。 
　-nouser 　找出不属于本地主机用户识别码的文件或目录。 
　-ok<执行指令> 　此参数的效果和指定"-exec"参数类似，但在执行指令之前会先询问用户，若回答"y"或"Y"，则放弃执行指令。 
　-path<范本样式> 　指定字符串作为寻找目录的范本样式。 
　-perm<权限数值> 　查找符合指定的权限数值的文件或目录。 
　-print 　假设find指令的回传值为True，就将文件或目录名称列出到标准输出。格式为每列一个名称，每个名称之前皆有"./"字符串。 
　-print0 　假设find指令的回传值为True，就将文件或目录名称列出到标准输出。格式为全部的名称皆在同一行。 
　-printf<输出格式> 　假设find指令的回传值为True，就将文件或目录名称列出到标准输出。格式可以自行指定。 
　-prune 　不寻找字符串作为寻找文件或目录的范本样式。 
　-regex<范本样式> 　指定字符串作为寻找文件或目录的范本样式。 
　-size<文件大小> 　查找符合指定的文件大小的文件。 
　-true 　将find指令的回传值皆设为True。 
　-typ<文件类型> 　只寻找符合指定的文件类型的文件。 
　-uid<用户识别码> 　查找符合指定的用户识别码的文件或目录。 
　-used<日数> 　查找文件或目录被更改之后在指定时间曾被存取过的文件或目录，单位以日计算。 
　-user<拥有者名称> 　查找符合指定的拥有者名称的文件或目录。 
　-version或--version 　显示版本信息。 
　-xdev 　将范围局限在先行的文件系统中。 
　-xtype<文件类型> 　此参数的效果和指定"-type"参数类似，差别在于它针对符号连接检查。
### git(gnu interactive tools)：文字模式下的文件管理员。 
语　　法：git
补充说明：git是用来管理文件的程序，它十分类似DOS下的Norton Commander，具有互动式操作界面。它的操作方法和Norton Commander几乎一样，略诉如下： 
  F1 ：执行info指令，查询指令相关信息，会要求您输入欲查询的名称。 
  F2 ：执行cat指令，列出文件内容。 
  F3 ：执行gitview指令，观看文件内容。 
  F4 ：执行vi指令，编辑文件内容。 
  F5 ：执行cp指令，复制文件或目录，会要求您输入目标文件或目录。 
  F6 ：执行mv指令，移动文件或目录，或是更改其名称，会要求您输入目标文件或目录。 
  F7 ：执行mkdir指令，建立目录。 
  F8 ：执行rm指令，删除文件或目录。 
  F9 ：执行make指令，批处理执行指令或编译程序时，会要求您输入相关命令。 
  F10 ：离开git文件管理员。
### gitview(gnu interactie tools viewer)：Hex/ASCII的看文件程序。 
语　　法：gitview [-bchilv][文件]  
补充说明：gitview指令可用于观看文件的内容，它会同时显示十六进制和ASCII格式的字码。 
参　　数：
  -b 　单色模式，不使用ANSI控制码显示彩色。 
  -c 　彩色模式，使用ANSI控制码显示色彩。 
  -h 　在线帮助。 
  -i 　显示存放gitview程序的所在位置。 
  -l 　不使用先前的显示字符。 
  -v 　显示版本信息。
### ln(link)：连接文件或目录。 
语　　法：ln [-bdfinsv][-S <字尾备份字符串>][-V <备份方式>][--help][--version][源文件或目录][目标文件或目录] 或 ln [-bdfinsv][-S <字尾备份字符串>][-V <备份方式>][--help][--version][源文件或目录...][目的目录] 
补充说明：ln指令用在连接文件或目录，如同时指定两个以上的文件或目录，且最后的目的地是一个已经存在的目录，则会把前面指定的所有文件或目录复制到该目录中。若同时指定多个文件或目录，且最后的目的地并非是一个已存在的目录，则会出现错误信息。 
参　　数：
　-b或--backup 　删除，覆盖目标文件之前的备份。 
　-d或-F或--directory 　建立目录的硬连接。 
　-f或--force 　强行建立文件或目录的连接，不论文件或目录是否存在。 
　-i或--interactive 　覆盖既有文件之前先询问用户。 
　-n或--no-dereference 　把符号连接的目的目录视为一般文件。 
　-s或--symbolic 　对源文件建立符号连接，而非硬连接。 
　-S<字尾备份字符串>或--suffix=<字尾备份字符串> 　用"-b"参数备份目标文件后，备份文件的字尾会被加上一个备份字符串，预设的字尾备份字符串是符号"~"，您可通过"-S"参数来改变它。 
　-v或--verbose 　显示指令执行过程。 
　-V<备份方式>或--version-control=<备份方式> 　用"-b"参数备份目标文件后，备份文件的字尾会被加上一个备份字符串，这个字符串不仅可用"-S"参数变更，当使用"-V"参数<备份方式>指定不同备份方式时，也会产生不同字尾的备份字符串。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### locate：查找文件。 
语　　法：locate [-d <数据库文件>][--help][--version][范本样式...]
补充说明：locate指令用于查找符合条件的文件，它会去保存文件与目录名称的数据库内，查找合乎范本样式条件的文件或目录。 
参　　数：
　-d<数据库文件>或--database=<数据库文件>  设置locate指令使用的数据库。locate指令预设的数据库位于/var/lib/slocate目录里，文件名为slocate.db，您可使用这个参数另行指定。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### lsattr(list attribute)：显示文件属性。 
语　　法：lsattr [-adlRvV][文件或目录...]  
补充说明：用chattr执行改变文件或目录的属性，可执行lsattr指令查询其属性。 
参　　数：
  -a 　显示所有文件和目录，包括以"."为名称开头字符的额外内建，现行目录"."与上层目录".."。 
  -d 　显示，目录名称，而非其内容。 
  -l 　此参数目前没有任何作用。 
  -R 　递归处理，将指定目录下的所有文件及子目录一并处理。 
  -v 　显示文件或目录版本。 
  -V 　显示版本信息。
### mattrib：变更或显示MS-DOS文件的属性。 
语　　法：mattrib [+/-][ahrs][-/X][文件] 或 mattrib [-/X][文件] 
补充说明：mattrib为mtools工具指令，模拟MS-DOS的attrib指令，可变更MS-DOS文件的属性。 
参　　数：
　+/- 　"+"代表开启指定的文件属性。"-"代表关闭指定的文件属性。 
　a 　备份属性。使备份程序可用来判断文件是否已经备份过。 
　h 　隐藏属性。 
　r 　只读属性。 
　s 　系统属性。 
　-/ 　列出指定目录及子目录下所有文件的属性。 
　-X 　显示文件属性时，中间不输入空格。
### mc(midnight commander)：提供一个菜单式的文件管理程序。
语　　法：mc [-abcdfhkPstuUVx][-C <参数>][-l <文件>][-v <文件>][目录] 
补充说明：执行mc之后，将会看到菜单式的文件管理程序，共分成４个部分。 
参　　数：
　-a 　当mc程序画线时不用绘图字符画线。 
　-b 　使用单色模式显示。 
　-c 　使用彩色模式显示。 
　-C<参数> 　指定显示的颜色。 
　-d 　不使用鼠标。 
　-f 　显示mc函数库所在的目录。 
　-h 　显示帮助。 
　-k 　重设softkeys成预设置。 
　-l<文件> 　在指定文件中保存ftpfs对话窗的内容。 
　-P 　程序结束时，列出最后的工作目录。 
　-s 　用慢速的终端机模式显示，在这模式下将减少大量的绘图及文字显示。 
　-t 　使用TEMPCAP变量设置终端机，而不使用预设置。 
　-u 　不用目前的shell程序。 
　-U 　使用目前的shell程序。 
　-v<文件> 　使用mc的内部编辑器来显示指定的文件。 
　-V 　显示版本信息。 
　-x 　指定以xterm模式显示。
### mcopy：复制MS-DOS文件。
语　　法：mcopy [-bnmpQt/][源文件][目标文件或目录]
补充说明：mcopy为mtools工具指令，可在MS-DOS文件系统中复制文件或是在MS-DOS与Linux的文件系统之间复制文件。mcopy会MS-DOS磁盘驱动器代号出现的位置来判断如何复制文件。
参　　数：
  -b 　批处理模式。 
  -n 　覆盖其他文件时，不需进行确认而直接覆盖。 
  -m 　将源文件修改时间设为目标文件的修改时间。 
  -p 　将源文件属性设为目标文件的属性。 
  -Q 　复制多个文件时，若发生错误，则立即结束程序。 
  -t 　转换文本文件。 
  -/ 　复制子目录以及其中的所有文件。
### mdel：删除MS-DOS文件。
语　　法：mdel [文件...]
补充说明：mdel为mtools工具指令，模拟MS-DOS的del指令，可删除MS-DOS文件系统中的文件。

### mdir：显示MS-DOS目录。
语　　法：mdir [-afwx/][目录]
补充说明：mdir为mtools工具指令，模拟MS-DOS的dir指令，可显示MS-DOS文件系统中的目录内容。
参　　数：
  -/   显示目录下所有子目录与文件。 
  -a 　显示隐藏文件。 
  -f 　不显示磁盘所剩余的可用空间。 
  -w 　仅显示目录或文件名称，并以横排方式呈现，以便一次能显示较多的目录或文件。 
  -X 　仅显示目录下所有子目录与文件的完整路径，不显示其他信息。
### mktemp：建立暂存文件。
语　　法：mktemp [-qu][文件名参数]
补充说明：mktemp可建立一个暂存文件，供shell script使用。
参　　数：
　-q 　执行时若发生错误，不会显示任何信息。 
　-u 　暂存文件会在mktemp结束前先行删除。 
　[文件名参数] 　文件名参数必须是以"自订名称.XXXXXX"的格式。
### mmove：在MS-DOS文件系统中，移动文件或目录，或更改名称。
语　　法：mmove [源文件或目录...][目标文件或目录]
补充说明：mmove为mtools工具指令，模拟MS-DOS的move指令，可在MS-DOS文件系统中移动现有的文件或目录，或是更改现有文件或目录的名称。
### mread：将MS-DOS文件复制到Linux/Unix的目录中。
语　　法：mread [MS-DOS文件...][Linux文件或目录]
补充说明：mread为mtools工具指令，可将MS-DOS文件复制到Linux的文件系统中。这个指令目前已经不常用，一般都使用mcopy指令来代替。
### mren：更改MS-DOS文件或目录的名称，或是移动文件或目录。
语　　法：mren [源文件或目录...][目标文件或目录]
补充说明：mren为MS-DOS工具指令，模拟MS-DOS的ren指令，可更改MS-DOS文件或目录名称。除此之外，ren也可移动文件或目录，但仅限于在同一设备内。
### mshowfat：显示MS-DOS文件在FAT中的记录。
语　　法：mshowfat [文件...]
补充说明：mshowfat为mtools工具指令，可显示MS-DOS文件在FAT中的记录编号。
### mtools：显示mtools支持的指令。
语　　法：mtools
补充说明：mtools为MS-DOS文件系统的工具程序，可模拟许多MS-DOS的指令。这些指令都是mtools的符号连接，因此会有一些共同的特性。
参　　数：
  -a 　长文件名重复时自动更改目标文件的长文件名。 
  -A 　短文件名重复但长文件名不同时自动更改目标文件的短文件名。 
  -o 　长文件名重复时，将目标文件覆盖现有的文件。 
  -O 　短文件名重复但长文件名不同时，将目标文件覆盖现有的文件。 
  -r 　长文件名重复时，要求用户更改目标文件的长文件名。 
  -R 　短文件名重复但长文件名煌保笥没Ц哪勘晡募亩涛募?nbsp;
  -s 　长文件名重复时，则不处理该目标文件。 
  -S 　短文件名重复但长文件名不同时，则不处理该目标文件。 
  -v 　执行时显示详细的说明。 
  -V 　显示版本信息。
### mtoolstest：测试并显示mtools的相关设置。
语　　法：mtoolstest
补充说明：mtoolstest为mtools工具指令，可读取与分析mtools的配置文件，并在屏幕上显示结果
### mv：移动或更名现有的文件或目录。
语　　法：mv [-bfiuv][--help][--version][-S <附加字尾>][-V <方法>][源文件或目录][目标文件或目录]
补充说明：mv可移动文件或目录，或是更改文件或目录的名称。
参　　数：
　-b或--backup 　若需覆盖文件，则覆盖前先行备份。 
　-f或--force 　若目标文件或目录与现有的文件或目录重复，则直接覆盖现有的文　件或目录。 
　-i或--interactive 　覆盖前先行询问用户。 
　-S<附加字尾>或
　--suffix=<附加字尾> 　与-b参数一并使用，可指定备份文件的所要附加的字尾。 
　-u或--update 　在移动或更改文件名时，若目标文件已存在，且其文件日期比源文件新，则不覆盖目标文件。 
　-v或--verbose 　执行时显示详细的信息。 
　-V=<方法>或
　--version-control=<方法> 　与-b参数一并使用，可指定备份的方法。 
　--help 　显示帮助。 
　--version 　显示版本信息。
### od(octal dump)：输出文件内容。
语　　法：od [-abcdfhilovx][-A <字码基数>][-j <字符数目>][-N <字符数目>][-s <字符串字符数>][-t <输出格式>][-w <每列字符数>][--help][--version][文件...]
补充说明：od指令会读取所给予的文件的内容，并将其内容以八进制字码呈现出来。
参　　数：
　-a 　此参数的效果和同时指定"-ta"参数相同。 
　-A<字码基数> 　选择要以何种基数计算字码。 
　-b 　此参数的效果和同时指定"-toC"参数相同。 
　-c 　此参数的效果和同时指定"-tC"参数相同。 
　-d 　此参数的效果和同时指定"-tu2"参数相同。 
　-f 　此参数的效果和同时指定"-tfF"参数相同。 
　-h 　此参数的效果和同时指定"-tx2"参数相同。 
　-i 　此参数的效果和同时指定"-td2"参数相同。 
　-j<字符数目>或--skip-bytes=<字符数目> 　略过设置的字符数目。 
　-l 　此参数的效果和同时指定"-td4"参数相同。 
　-N<字符数目>或--read-bytes=<字符数目> 　到设置的字符数目为止。 
　-o 　此参数的效果和同时指定"-to2"参数相同。 
　-s<字符串字符数>或--strings=<字符串字符数> 　只显示符合指定的字符数目的字符串。 
　-t<输出格式>或--format=<输出格式> 　设置输出格式。 
　-v或--output-duplicates 　输出时不省略重复的数据。 
　-w<每列字符数>或--width=<每列字符数> 　设置每列的最大字符数。 
　-x 　此参数的效果和同时指定"-h"参数相同。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### paste：合并文件的列。
语　　法：paste [-s][-d <间隔字符>][--help][--version][文件...]
补充说明：paste指令会把每个文件以列对列的方式，一列列地加以合并。
参　　数：
　-d<间隔字符>或--delimiters=<间隔字符> 　用指定的间隔字符取代跳格字符。 
　-s或--serial 　串列进行而非平行处理。 
　--help 　在线帮助。 
　--version 　显示帮助信息。
### patch：修补文件。
语　　法：patch [-bceEflnNRstTuvZ][-B <备份字首字符串>][-d <工作目录>][-D <标示符号>][-F <监别列数>][-g <控制数值>][-i <修补文件>][-o <输出文件>][-p <剥离层级>][-r <拒绝文件>][-V <备份方式>][-Y <备份字首字符串>][-z <备份字尾字符串>][--backup-if　　 -mismatch][--binary][--help][--nobackup-if-mismatch][--verbose][原始文件 <修补文件>] 或 path [-p <剥离层级>] < [修补文件]
补充说明：patch指令让用户利用设置修补文件的方式，修改，更新原始文件。倘若一次仅修改一个文件，可直接在指令列中下达指令依序执行。如果配合修补文件的方式则能一次修补大批文件，这也是Linux系统核心的升级方法之一。
参　　数：
　-b或--backup 　备份每一个原始文件。 
　-B<备份字首字符串>或--prefix=<备份字首字符串> 　设置文件备份时，附加在文件名称前面的字首字符串，该字符串可以是路径名称。 
　-c或--context 　把修补数据解译成关联性的差异。 
　-d<工作目录>或--directory=<工作目录> 　设置工作目录。 
　-D<标示符号>或--ifdef=<标示符号> 　用指定的符号把改变的地方标示出来。 
　-e或--ed 　把修补数据解译成ed指令可用的叙述文件。 
　-E或--remove-empty-files 　若修补过后输出的文件其内容是一片空白，则移除该文件。 
　-f或--force 　此参数的效果和指定"-t"参数类似，但会假设修补数据的版本为新　版本。 
　-F<监别列数>或--fuzz<监别列数> 　设置监别列数的最大值。 
　-g<控制数值>或--get=<控制数值> 　设置以RSC或SCCS控制修补作业。 
　-i<修补文件>或--input=<修补文件> 　读取指定的修补问家你。 
　-l或--ignore-whitespace 　忽略修补数据与输入数据的跳格，空格字符。 
　-n或--normal 　把修补数据解译成一般性的差异。 
　-N或--forward 　忽略修补的数据较原始文件的版本更旧，或该版本的修补数据已使　用过。 
　-o<输出文件>或--output=<输出文件> 　设置输出文件的名称，修补过的文件会以该名称存放。 
　-p<剥离层级>或--strip=<剥离层级> 　设置欲剥离几层路径名称。 
　-f<拒绝文件>或--reject-file=<拒绝文件> 　设置保存拒绝修补相关信息的文件名称，预设的文件名称为.rej。 
　-R或--reverse 　假设修补数据是由新旧文件交换位置而产生。 
　-s或--quiet或--silent 　不显示指令执行过程，除非发生错误。 
　-t或--batch 　自动略过错误，不询问任何问题。 
　-T或--set-time 　此参数的效果和指定"-Z"参数类似，但以本地时间为主。 
　-u或--unified 　把修补数据解译成一致化的差异。 
　-v或--version 　显示版本信息。 
　-V<备份方式>或--version-control=<备份方式> 　用"-b"参数备份目标文件后，备份文件的字尾会被加上一个备份字符串，这个字符串不仅可用"-z"参数变更，当使用"-V"参数指定不同备份方式时，也会产生不同字尾的备份字符串。
　-Y<备份字首字符串>或--basename-prefix=--<备份字首字符串> 　设置文件备份时，附加在文件基本名称开头的字首字符串。
　-z<备份字尾字符串>或--suffix=<备份字尾字符串> 　此参数的效果和指定"-B"参数类似，差别在于修补作业使用的路径与文件名若为src/linux/fs/super.c，加上"backup/"字符串后，文件super.c会备份于/src/linux/fs/backup目录里。 
　-Z或--set-utc 　把修补过的文件更改，存取时间设为UTC。 
　--backup-if-mismatch 　在修补数据不完全吻合，且没有刻意指定要备份文件时，才备份文件。 
　--binary 　以二进制模式读写数据，而不通过标准输出设备。 
　--help 　在线帮助。 
　--nobackup-if-mismatch 　在修补数据不完全吻合，且没有刻意指定要备份文件时，不要备份文件。 
　--verbose 　详细显示指令的执行过程。
### rcp(remote copy)：远端复制文件或目录。
语　　法：rcp [-pr][源文件或目录][目标文件或目录] 或 rcp [-pr][源文件或目录...][目标文件] 
补充说明：rcp指令用在远端复制文件或目录，如同时指定两个以上的文件或目录，且最后的目的地是一个已经存在的目录，则它灰把前面指定的所有文件或目录复制到该目录中。
参　　数：
  -p 　保留源文件或目录的属性，包括拥有者，所属群组，权限与时间。 
  -r 　递归处理，将指定目录下的文件与子目录一并处理。
### rhmask：产生与还原加密文件。
语　　法：rhmask [加密文件][输出文件] 或 rhmask [-d][加密文件][源文件][输出文件]
补充说明：执行rhmask指令可制作加密过的文件，方便用户在公开的网络上传输该文件，而不至于被任意盗用。
参　　数：
  -d 　产生加密过的文件。
### rm(remove)：删除文件或目录。
语　　法：rm [-dfirv][--help][--version][文件或目录...]
补充说明：执行rm指令可删除文件或目录，如欲删除目录必须加上参数"-r"，否则预设仅会删除文件。 
参　　数：
　-d或--directory 　直接把欲删除的目录的硬连接数据删成0，删除该目录。 
　-f或--force 　强制删除文件或目录。 
　-i或--interactive 　删除既有文件或目录之前先询问用户。 
　-r或-R或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v或--verbose 　显示指令执行过程。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### slocate(secure locate)：查找文件或目录。
语　　法：slocate [-u][--help][--version][-d <目录>][查找的文件]
补充说明：slocate本身具有一个数据库，里面存放了系统中文件与目录的相关信息。
参　　数：
　-d<目录>或--database=<目录> 　指定数据库所在的目录。 
　-u 　更新slocate数据库。 
　--help 　显示帮助。 
　--version 　显示版本信息。
### split：切割文件。
语　　法：split [--help][--version][-<行数>][-b <字节>][-C <字节>][-l <行数>][要切割的文件][输出文件名]
补充说明：split可将文件切成较小的文件，预设每1000行会切成一个小文件。
参　　数：
　-<行数>或-l<行数> 　指定每多少行就要切成一个小文件。 
　-b<字节> 　指定每多少字就要切成一个小文件。 
　-C<字节> 　与-b参数类似，但切割时尽量维持每行的完整性。 
　--help 　显示帮助。 
　--version 　显示版本信息。 
　[输出文件名] 　设置切割后文件的前置文件名，split会自动在前置文件名后再加上编号。
### tee：读取标准输入的数据，并将其内容输出成文件。
语　　法：tee [-ai][--help][--version][文件...]
补充说明：tee指令会从标准输入设备读取数据，将其内容输出到标准输出设备，同时保存成文件。
参　　数：
　-a或--append 　附加到既有文件的后面，而非覆盖它． 
　-i-i或--ignore-interrupts 　忽略中断信号。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### tmpwatch(temporary watch)：删除暂存文件。
语　　法：tmpwatch [-afqv][--test][超期时间][目录...]
补充说明：执行tmpwatch指令可删除不必要的暂存文件，您可以设置文件超期时间，单位以小时计算。
参　　数：
　-a或--all 　删除任何类型的文件。 
　-f或--force 　强制删除文件或目录，其效果类似rm指令的"-f"参数。 
　-q或--quiet 　不显示指令执行过程。 
　-v或--verbose 　详细显示指令执行过程。 
　-test 　仅作测试，并不真的删除文件或目录。
### touch：改变文件或目录时间。
语　　法：touch [-acfm][-d <日期时间>][-r <参考文件或目录>][-t <日期时间>][--help]　　 [--version][文件或目录...] 或 touch [-acfm][--help][--version][日期时间][文件或目录...] 
补充说明：使用touch指令可更改文件或目录的日期时间，包括存取时间和更改时间。
参　　数：
　-a或--time=atime或--time=access或--time=use 　只更改存取时间。 
　-c或--no-create 　不建立任何文件。 
　-d<时间日期> 　使用指定的日期时间，而非现在的时间。 
　-f 　此参数将忽略不予处理，仅负责解决BSD版本touch指令的兼容性问题。 
　-m或--time=mtime或--time=modify 　只更改变动时间。 
　-r<参考文件或目录> 　把指定文件或目录的日期时间，统统设成和参考文件或目录的日期时间相同。 
　-t<日期时间> 　使用指定的日期时间，而非现在的时间。 
　--help 　在线帮助。 
　--version 　显示版本信息。
### umask：指定在建立文件时预设的权限掩码。
语　　法：umask [-S][权限掩码]
补充说明：umask可用来设定[权限掩码]。[权限掩码]是由3个八进制的数字所组成，将现有的存取权限减掉权限掩码后，即可产生建立文件时预设的权限。
参　　数：
  -S 　以文字的方式来表示权限掩码。
### whereis：查找文件。
语　　法：whereis [-bfmsu][-B <目录>...][-M <目录>...][-S <目录>...][文件...]
补充说明：whereis指令会在特定目录中查找符合条件的文件。这些文件的烈性应属于原始代码，二进制文件，或是帮助文件。
参　　数：
　-b 　只查找二进制文件。 
　-B<目录> 　只在设置的目录下查找二进制文件。 
　-f 　不显示文件名前的路径名称。 
　-m 　只查找说明文件。 
　-M<目录> 　只在设置的目录下查找说明文件。 
　-s 　只查找原始代码文件。 
　-S<目录> 　只在设置的目录下查找原始代码文件。 
　-u 　查找不包含指定类型的文件。
### which：查找文件。
语　　法：which [文件...]
补充说明：which指令会在环境变量$PATH设置的目录里查找符合条件的文件。
参　　数：
　-n<文件名长度> 　指定文件名长度，指定的长度必须大于或等于所有文件中最长的文件名。 
　-p<文件名长度> 　与-n参数相同，但此处的<文件名长度>包括了文件的路径。 
　-w 　指定输出时栏位的宽度。 
　-V 　显示版本信息。
cat
使用权限：所有使用者 
使用方式：cat [-AbeEnstTuv] [--help] [--version] fileName 
说明：把档案串连接后传到基本输出（萤幕或加 > fileName 到另一个档案）
参数：
-n 或 --number 由 1 开始对所有输出的行数编号
-b 或 --number-nonblank 和 -n 相似，只不过对于空白行不编号
-s 或 --squeeze-blank 当遇到有连续两行以上的空白行，就代换为一行的空白行
-v 或 --show-nonprinting
范例：
cat -n textfile1 > textfile2 把 textfile1 的档案内容加上行号后输入 textfile2 这个档案里
cat -b textfile1 textfile2 >> textfile3 把 textfile1 和 textfile2 的档案内容加上行号（空白行不加）之后将内容附加到 textfile3 里。
### chattr(change attribute)：改变文件属性。
语　　法：chattr [-RV][-v<版本编号>][+/-/=<属性>][文件或目录...]
补充说明：这项指令可改变存放在ext2文件系统上的文件或目录属性，这些属性共有以下8种模式： 
　a：让文件或目录仅供附加用途。 
　b：不更新文件或目录的最后存取时间。 
　c：将文件或目录压缩后存放。 
　d：将文件或目录排除在倾倒操作之外。 
　i：不得任意更动文件或目录。 
　s：保密性删除文件或目录。 
　S：即时更新文件或目录。 
　u：预防以外删除。 
 参　　数：
　-R  递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v<版本编号>  设置文件或目录版本。 
　-V  显示指令执行过程。 
　+<属性>  开启文件或目录的该项属性。 
　-<属性>  关闭文件或目录的该项属性。 
　=<属性>  指定文件或目录的该项属性。
### chgrp(change group)：变更文件或目录的所属群组。
语　　法：chgrp [-cfhRv][--help][--version][所属群组][文件或目录...] 或 chgrp [-cfhRv][--help][--reference=<参考文件或目录>][--version][文件或目录...] 
补充说明：在UNIX系统家族里，文件或目录权限的掌控以拥有者及所属群组来管理。您可以使用chgrp指令去变更文件与目录的所属群组，设置方式采用群组名称或群组识别码皆可。 
参　　数：
　-c或--changes  效果类似"-v"参数，但仅回报更改的部分。 
　-f或--quiet或--silent 　不显示错误信息。 
　-h或--no-dereference 　只对符号连接的文件作修改，而不更动其他任何相关文件。 
　-R或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v或--verbose 　显示指令执行过程。 
　--help 　在线帮助。 
　--reference=<参考文件或目录> 　把指定文件或目录的所属群组全部设成和参考文件或目录的所属群组相同。 
　--version 　显示版本信息。
### chmod(change mode)：变更文件或目录的权限。
语　　法：chmod [-cfRv][--help][--version][<权限范围>+/-/=<权限设置...>][文件或目录...] 或 chmod [-cfRv][--help][--version][数字代号][文件或目录...] 或 chmod [-cfRv][--help][--reference=<参考文件或目录>][--version][文件或目录...] 
补充说明：在UNIX系统家族里，文件或目录权限的控制分别以读取，写入，执行3种一般权限来区分，另有3种特殊权限可供运用，再搭配拥有者与所属群组管理权限范围。您可以使用chmod指令去变更文件与目录的权限，设置方式采用文字或数字代号皆可。符号连接的权限无法变更，如果您对符号连接修改权限，其改变会作用在被连接的原始文件。权限范围的表示法如下： 
　u：User，即文件或目录的拥有者。 
　g：Group，即文件或目录的所属群组。 
　o：Other，除了文件或目录拥有者或所属群组之外，其他用户皆属于这个范围。 
　a：All，即全部的用户，包含拥有者，所属群组以及其他用户。 
 　有关权限代号的部分，列表于下： 
　r：读取权限，数字代号为"4"。 
　w：写入权限，数字代号为"2"。 
　x：执行或切换权限，数字代号为"1"。 
　-：不具任何权限，数字代号为"0"。 
　s：特殊?b>功能说明：变更文件或目录的权限。
 参　　数：
　-c或--changes 　效果类似"-v"参数，但仅回报更改的部分。 
　-f或--quiet或--silent 　不显示错误信息。 
　-R或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v或--verbose 　显示指令执行过程。 
　--help 　在线帮助。 
　--reference=<参考文件或目录> 　把指定文件或目录的权限全部设成和参考文件或目录的权限相同 
　--version 　显示版本信息。 
　<权限范围>+<权限设置> 　开启权限范围的文件或目录的该项权限设置。 
　<权限范围>-<权限设置> 　关闭权限范围的文件或目录的该项权限设置。 
　<权限范围>=<权限设置> 　指定权限范围的文件或目录的该项权限设置。
 如chmod a+x …
### chown(change owner)：变更文件或目录的拥有者或所属群组。
语　　法：chown [-cfhRv][--dereference][--help][--version][拥有者.<所属群组>][文件或目录..] 或chown [-chfRv][--dereference][--help][--version][.所属群组][文件或目录... ...] 或chown [-cfhRv][--dereference][--help][--reference=<参考文件或目录>][--version][文件或目录...] 
补充说明：在UNIX系统家族里，文件或目录权限的掌控以拥有者及所属群组来管理。您可以使用chown指令去变更文件与目录的拥有者或所属群组，设置方式采用用户名称或用户识别码皆可，设置群组则用群组名称或群组识别码。 
参　　数：
　-c或--changes 　效果类似"-v"参数，但仅回报更改的部分。 
　-f或--quite或--silent 　不显示错误信息。 
　-h或--no-dereference 　之对符号连接的文件作修改，而不更动其他任何相关文件。 
　-R或--recursive 　递归处理，将指定目录下的所有文件及子目录一并处理。 
　-v或--version 　显示指令执行过程。 
　--dereference 　效果和"-h"参数相同。 
　--help 　在线帮助。 
　--reference=<参考文件或目录> 　把指定文件或目录的拥有者与所属群组全部设成和参考文件或目　录的拥有者与所属群组相同。 
　--version 　显示版本信息。
### cksum(check sum)：检查文件的CRC是否正确。 
语　　法：cksum [--help][--version][文件...]  
补充说明：CRC是一种排错检查方式，该演算法的标准由CCITT所指定，至少可检测到99.998%的已知错误。指定文件交由cksum演算，它会回报计算结果，供用户核对文件是否正确无误。若不指定任何文件名称或是所给予的文件名为"-"，则cksum指令会从标准输入设备读取数据。 
参　　数：
　--help 　在线帮助。 
　--version 　显示版本信息。
### cmp(compare)：比较两个文件是否有差异。
语　　法：cmp [-clsv][-i <字符数目>][--help][第一个文件][第二个文件]
补充说明：当相互比较的两个文件完全一样时，则该指令不会显示任何信息。若发现有所差异，预设会标示出第一个不同之处的字符和列数编号。若不指定任何文件名称或是所给予的文件名为"-"，则cmp指令会从标准输入设备读取数据。
参　　数：
　-c或--print-chars 　除了标明差异处的十进制字码之外，一并显示该字符所对应字符。 
　-i<字符数目>或--ignore-initial=<字符数目> 　指定一个数目。 
　-l或--verbose 　标示出所有不一样的地方。 
　-s或--quiet或--silent 　不显示错误信息。 
　-v或--version 　显示版本信息。 
　--help 　在线帮助。
### cp(copy)：复制文件或目录。 
语　　法：cp [-abdfilpPrRsuvx][-S <备份字尾字符串>][-V <备份方式>][--help][--spares=<使用时机>][--version][源文件或目录][目标文件或目录] [目的目录] 
补充说明：cp指令用在复制文件或目录，如同时指定两个以上的文件或目录，且最后的目的地是一个已经存在的目录，则它会把前面指定的所有文件或目录复制到该目录中。若同时指定多个文件或目录，而最后的目的地并非是一个已存在的目录，则会出现错误信息。
参　　数：
　-a或--archive 　此参数的效果和同时指定"-dpR"参数相同。 
　-b或--backup 　删除，覆盖目标文件之前的备份，备份文件会在字尾加上一个备份字符串。 
　-d或--no-dereference 　当复制符号连接时，把目标文件或目录也建立为符号连接，并指向与源文件或目录连接的原始文件或目录。 
　-f或--force 　强行复制文件或目录，不论目标文件或目录是否已存在。 
　-i或--interactive 　覆盖既有文件之前先询问用户。 
　-l或--link 　对源文件建立硬连接，而非复制文件。 
　-p或--preserve 　保留源文件或目录的属性。 
　-P或--parents 　保留源文件或目录的路径。 
　-r 　递归处理，将指定目录下的文件与子目录一并处理。 
　-R或--recursive 　递归处理，将指定目录下的所有文件与子目录一并处理。 
　-s或--symbolic-link 　对源文件建立符号连接，而非复制文件。 
　-S<备份字尾字符串>或--suffix=<备份字尾字符串> 　用"-b"参数备份目标文件后，备份文件的字尾会被加上一个备份字符串，预设的备份字尾字符串是符号"~"。 
　-u或--update 　使用这项参数后只会在源文件的更改时间较目标文件更新时或是　名称相互对应的目标文件并不存在，才复制文件。 
　-v或--verbose 　显示指令执行过程。 
　-V<备份方式>或--version-control=<备份方式> 　用"-b"参数备份目标文件后，备份文件的字尾会被加上一个备份字符串，这字符串不仅可用"-S"参数变更，当使用"-V"参数指定不同备份方式时，也会产生不同字尾的备份字串。  
　-x或--one-file-system 　复制的文件或目录存放的文件系统，必须与cp指令执行时所处的文件系统相同，否则不予复制。 
　--help 　在线帮助。 
　--sparse=<使用时机> 　设置保存稀疏文件的时机。 
　--version 　显示版本信息。
cut
使用权限：所有使用者
用法：cut -cnum1-num2 filename
说明：显示每行从开头算起 num1 到 num2 的文字。
范例：
        shell>> cat example
        test2
        this is test1
        shell>> cut -c0-6 example       ## print 开头算起前 6 个字元
        test2
        this i 
### indent：调整C原始代码文件的格式。
 语　　法：indent [参数][源文件] 或 indent [参数][源文件][-o 目标文件]
补充说明：indent可辨识C的原始代码文件，并加以格式化，以方便程序设计师阅读。 
参　　数：
　-bad或--blank-lines-after-declarations  　在声明区段或加上空白行。 
　-bap或--blank-lines-after-procedures 　在程序或加上空白行。 
　-bbb或--blank-lines-after-block-comments 　在注释区段后加上空白行。 
　-bc或--blank-lines-after-commas  　在声明区段中，若出现逗号即换行。 
　-bl或--braces-after-if-line 　if(或是else,for等等)与后面执行区段的"{"不同行，且"}"自成一行。 
　-bli<缩排格数>或--brace-indent<缩排格数> 　设置{ }缩排的格数。 
　-br或--braces-on-if-line 　if(或是else,for等等)与后面执行跛段的"{"不同行，且"}"自成一行。 
　-bs或--blank-before-sizeof 　在sizeof之后空一格。 
　-c<栏数>或--comment-indentation<栏数> 　将注释置于程序码右侧指定的栏位。 
　-cd<栏数>或--declaration-comment-column<栏数> 　将注释置于声明右侧指定的栏位。 
　-cdb或--comment-delimiters-on-blank-lines 　注释符号自成一行。 
　-ce或--cuddle-else 　将else置于"}"(if执行区段的结尾)之后。 
　-ci<缩排格数>或--continuation-indentation<缩排格数> 　叙述过长而换行时，指定换行后缩排的格数。 
　-cli<缩排格数>或--case-indentation-<缩排格数> 　使用case时，switch缩排的格数。 
　-cp<栏数>或-else-endif-column<栏数> 　将注释置于else与elseif叙述右侧定的栏位。 
　-cs或--space-after-cast 　在cast之后空一格。 
　-d<缩排格数>或-line-comments-indentation<缩排格数> 　针对不是放在程序码右侧的注释，设置其缩排格数。 
　-di<栏数>或--declaration-indentation<栏数> 　将声明区段的变量置于指定的栏位。 
　-fc1或--format-first-column-comments 　针对放在每行最前端的注释，设置其格式。 
　-fca或--format-all-comments 　设置所有注释的格式。 
　-gnu或--gnu-style 　指定使用GNU的格式，此为预设值。 
　-i<格数>或--indent-level<格数> 　设置缩排的格数。 
　-ip<格数>或--parameter-indentation<格数> 　设置参数的缩排格数。 
　-kr或--k-and-r-style 　指定使用Kernighan&Ritchie的格式。 
　-lp或--continue-at-parentheses 　叙述过长而换行，且叙述中包含了括弧时，将括弧中的每行起始栏位内容垂直对其排列。 
　-nbad或--no-blank-lines-after-declarations 　在声明区段后不要加上空白行。 
　-nbap或--no-blank-lines-after-procedures 　在程序后不要加上空白行。 
　-nbbb或--no-blank-lines-after-block-comments 　在注释区段后不要加上空白行。 
　-nbc或--no-blank-lines-after-commas 　在声明区段中，即使出现逗号，仍旧不要换行。 
　-ncdb或--no-comment-delimiters-on-blank-lines 　注释符号不要自成一行。 
　-nce或--dont-cuddle-else 　不要将else置于"}"之后。 
　-ncs或--no-space-after-casts 　不要在cast之后空一格。 
　-nfc1或--dont-format-first-column-comments 　不要格式化放在每行最前端的注释。 
　-nfca或--dont-format-comments 　不要格式化任何的注释。 
　-nip或--no-parameter-indentation 　参数不要缩排。 
　-nlp或--dont-line-up-parentheses 　叙述过长而换行，且叙述中包含了括弧时，不用将括弧中的每行起始栏位垂直对其排列。 
　-npcs或--no-space-after-function-call-names 　在调用的函数名称之后，不要加上空格。 
　-npro或--ignore-profile 　不要读取indent的配置文件.indent.pro。 
　-npsl或--dont-break-procedure-type 　程序类型与程序名称放在同一行。 
　-nsc或--dont-star-comments 　注解左侧不要加上星号(*)。 
　-nsob或--leave-optional-semicolon 　不用处理多余的空白行。 
　-nss或--dont-space-special-semicolon  　若for或while区段仅有一行时，在分号前不加上空格。 
　-nv或--no-verbosity 　不显示详细的信息。 
　-orig或--original 　使用Berkeley的格式。 
　-pcs或--space-after-procedure-calls 　在调用的函数名称与"{"之间加上空格。 
　-psl或--procnames-start-lines 　程序类型置于程序名称的前一行。 
　-sc或--start-left-side-of-comments 　在每行注释左侧加上星号(*)。 
　-sob或--swallow-optional-blank-lines 　删除多余的空白行。 
　-ss或--space-special-semicolon 　若for或swile区段今有一行时，在分号前加上空格。 
　-st或--standard-output 　将结果显示在标准输出设备。 
　-T 　数据类型名称缩排。 
　-ts<格数>或--tab-size<格数> 　设置tab的长度。 
　-v或--verbose 　执行时显示详细的信息。 
　-version 　显示版本信息。

 
## linux文件传输命令
### bye：中断FTP连线并结束程序。
语　　法：bye
补充说明：在ftp模式下，输入bye即可中断目前的连线作业，并结束ftp的执行。
### ftp(file transfer protocol)：设置文件系统相关功能。
语　　法：ftp [-dignv][主机名称或IP地址]
补充说明：FTP是ARPANet的标准文件传输协议，该网络就是现今Internet的前身。
参　　数：
  -d   详细显示指令执行过程，便于排错或分析程序执行的情形。
  -i   关闭互动模式，不询问任何问题。
  -g   关闭本地主机文件名称支持特殊字符的扩充特性。
  -n   不使用自动登陆。
  -v   显示指令执行过程。
### ftpcount：显示目前以FTP登入的用户人数。
语　　法：ftpcount
补充说明：执行这项指令可得知目前用FTP登入系统的人数以及FTP登入人数的上限。
### ftpshut(ftp shutdown)：在指定的时间关闭FTP服务器。
语　　法：ftpshut [-d<分钟>][-l<分钟>][关闭时间]["警告信息"]
补充说明：本指令提供系统管理者在设置的时间关闭FTP服务器，且能在关闭之前发出警告信息通知用户。关闭时间若设置后为"none"，则会马上关闭服务器。如果采用"+30"的方式来设置表示服务器在30分钟之后关闭。依次类推，假设使用"1130"的格式则代表服务器会在每日的11时30分关闭，时间格式为24小时制。FTP服务器关闭后，在/etc目录下会产生一个名称为shutmsg的文件，把它删除后即可再度启动FTP服务器的功能。
参　　数：
  -d<分钟>   切断所有FTP连线时间。
  -l<分钟>   停止接受FTP登入的时间。
### ftpwho：显示目前所有以FTP登入的用户信息。
语　　法：ftpwho
补充说明：执行这项指令可得知目前用FTP登入系统的用户有那些人，以及他们正在进行的操作。
### ncftp(nc file transfer protocol)：传输文件。
语　　法：ncftp [主机或IP地址]
补充说明：FTP让用户得以下载存放于服务器主机的文件，也能将文件上传到远端主机放置。NcFTP是文字模式FTP程序的佼佼者，它具备多样特色， 包括显示传输速率，下载进度，自动续传，标住书签，可通过防火墙和代理服务器等。
### tftp(trivial file transfer protocol)：传输文件。
语　　法：tftp [主机名称或IP地址]
补充说明：FTP让用户得以下载存放于远端主机的文件，也能将文件上传到远端主机放置。tftp是简单的文字模式ftp程序，它所使用的指令和FTP类似。
### uucico：UUCP文件传输服务程序。
语　　法：uucico [-cCDefqvwz][-i<类型>][-I<文件>][-p<连接端口号码>][-][-rl][-s<主机>][-S<主机>][-u<用户>][-x<类型>][--help]
补充说明：uucico是用来处理uucp或uux送到队列的文件传输工具。uucico有两种工作模式：主动模式和附属模式。当在主动模式下时，uucico会调用远端主机；在附属模式下时，uucico则接受远端主机的调用。
参　　数：
  -c或--quiet   当不执行任何工作时，不要更改记录文件的内容及更新目前的状态。
  -C或--ifwork   当有工作要执行时，才调用-s或-S参数所指定主机。
  -D或--nodetach   不要与控制终端机离线。
  -e或--loop   在附属模式下执行，并且出现要求登入的提示画面。
  -f或--force   当执行错误时，不等待任何时间即重新调用主机。
  -i<类型>或--stdin<类型>   当使用到标准输入设备时，指定连接端口的类型。
  -I<文件>--config<文件>   指定使用的配置文件。
  -l或--prompt   出现要求登入的提示画面。
  -p<连接端口号码>或-port<连接端口号码>   指定连接端口号码。
  -q或--quiet   不要启动uuxqt服务程序。
  -r0或--slave   以附属模式启动。
  -s<主机>或--system<主机>   调用指定的主机。
  -u<用户>或--login<用户>   指定登入的用户帐号，而不允许输入任意的登入帐号。
  -v或--version   显示版本信息，并且结束程序。
  -w或--wait   在主动模式下，当执行调用动作时，则出现要求登入的提示画面。
  -x<类型>或-X<类型>或outgoing-debug<类型>   启动指定的排错模式。
  -z或--try-next   当执行不成功时，尝试下一个选择而不结束程序。
  --help   显示帮助，并且结束程序。
### uucp：在Unix系统之间传送文件。
语　　法：uucp [-cCdfjmrRtvW][-g<等级>][-I<配置文件>][-n<用户>][-x<类型>][--help][...来源][目的]
补充说明：UUCP为Unix系统之间，通过序列线来连线的协议。uucp使用UUCP协议，主要的功能为传送文件。
参　　数：
  -c或--nocopy   不用将文件复制到缓冲区。
  -C或--copy   将文件复制到缓冲区。
  -d或--directiories   在传送文件时，自动在[目的]建立必要的目录。
  -f或--nodirectiories   在传送文件时，若需要在[目的]建立目录，则放弃执行该作业。
  -g<等级>或--grade<等级>   指定文件传送作业的优先顺序。
  -I<配置文件>或--config<配置文件>   指定uucp配置文件。
  -j或--jobid   显示作业编号。
  -m或--mail   作业结束后，以电子邮件报告作业是否顺利完成。
  -n<用户>或--notify<用户>   作业结束后，以电子邮件向指定的用户报告作业是否顺利完成。
  -r或--nouucico   不要立即启动uucico服务程序，仅将作业送到队列中，待稍后再执行。
  -R或--recursive   若[来源]为目录，则将整个目录包含子目录复制到[目的]。
  -t或--uuto   将最后一个参数视为"主机名!用户"。
  -v或--version   显示版本信息。
  -W或--noexpand   不要将目前所在的目录加入路径。
  -x<类型>或--debug<类型>启动指定的排错模式。
  --help   显示帮助。
  [源...]   指定源文件或路径。
  [目的]   指定目标文件或路径。
### uupick：处理传送进来的文件。
语　　法：uupick [-v][-I<配置文件>][-s<主机>][-x<层级>][--help]
补充说明：当其他主机通过UUCP将文件传送进来时，可利用uupick指令取出这些文件。
参　　数：
  -I<配置文件>或--config<配置文件>   指定配置文件。
  -s<主机>或--system<主机>   处理由指定主机传送过来的文件。
  -v或--version   显示版本信息。
  --help   显示帮助。
### uuto：将文件传送到远端的UUCP主机。
语　　法：uuto [文件][目的]
补充说明：uuto为script文件，它实际上会执行uucp，用来将文件传送到远端UUCP主机，并在完成工作后，以邮件通知远端主机上的用户。
参　　数：
  相关参数请参考uucp指令。
 
## linux磁盘管理命令
### cd(change directory)：切换目录。
语　　法：cd [目的目录]
补充说明：cd指令可让用户在不同的目录间切换，但该用户必须拥有足够的权限进入目的目录。
### df(disk free)：显示磁盘的相关信息。
语　　法：df [-ahHiklmPT][--block-size=<区块大小>][-t <文件系统类型>][-x <文件系统类型>][--help][--no-sync][--sync][--version][文件或设备]
补充说明：df可显示磁盘的文件系统与使用情形。
参　　数：
  -a或--all   包含全部的文件系统。
  --block-size=<区块大小>   以指定的区块大小来显示区块数目。
  -h或--human-readable   以可读性较高的方式来显示信息。
  -H或--si   与-h参数相同，但在计算时是以1000 Bytes为换算单位而非1024 Bytes。
  -i或--inodes   显示inode的信息。
  -k或--kilobytes   指定区块大小为1024字节。
  -l或--local   仅显示本地端的文件系统。
  -m或--megabytes   指定区块大小为1048576字节。
  --no-sync   在取得磁盘使用信息前，不要执行sync指令，此为预设值。
  -P或--portability   使用POSIX的输出格式。
  --sync   在取得磁盘使用信息前，先执行sync指令。
  -t<文件系统类型>或--type=<文件系统类型>   仅显示指定文件系统类型的磁盘信息。
  -T或--print-type   显示文件系统的类型。
  -x<文件系统类型>或--exclude-type=<文件系统类型>   不要显示指定文件系统类型的磁盘信息。
  --help   显示帮助。
  --version   显示版本信息。
  [文件或设备]   指定磁盘设备。
### dirs：显示目录记录。
语　　法：dirs [+/-n -l]
补充说明：显示目录堆叠中的记录。
参　　数：
  +n   显示从左边算起第n笔的目录。
  -n   显示从右边算起第n笔的目录。
  -l   显示目录完整的记录。
### du(disk usage)：显示目录或文件的大小。
语　　法：du [-abcDhHklmsSx][-L <符号连接>][-X <文件>][--block-size][--exclude=<目录或文件>][--max-depth=<目录层数>][--help][--version][目录或文件]
补充说明：du会显示指定的目录或文件所占用的磁盘空间。
参　　数：
  -a或-all   显示目录中个别文件的大小。
  -b或-bytes   显示目录或文件大小时，以byte为单位。
  -c或--total   除了显示个别目录或文件的大小外，同时也显示所有目录或文件的总和。
  -D或--dereference-args   显示指定符号连接的源文件大小。
  -h或--human-readable   以K，M，G为单位，提高信息的可读性。
  -H或--si   与-h参数相同，但是K，M，G是以1000为换算单位。
  -k或--kilobytes   以1024 bytes为单位。
  -l或--count-links   重复计算硬件连接的文件。
  -L<符号连接>或--dereference<符号连接>   显示选项中所指定符号连接的源文件大小。
  -m或--megabytes   以1MB为单位。
  -s或--summarize   仅显示总计。
  -S或--separate-dirs   显示个别目录的大小时，并不含其子目录的大小。
  -x或--one-file-xystem   以一开始处理时的文件系统为准，若遇上其它不同的文件系统目录则略过。
  -X<文件>或--exclude-from=<文件>   在<文件>指定目录或文件。
  --exclude=<目录或文件>   略过指定的目录或文件。
  --max-depth=<目录层数>   超过指定层数的目录后，予以忽略。
  --help   显示帮助。
  --version   显示版本信息。
### edquota(edit quota)：编辑用户或群组的quota。
语　　法：edquota [-p <源用户名称>][-ug][用户或群组名称...] 或 edquota [-ug] -t
补充说明：edquota预设会使用vi来编辑使用者或群组的quota设置。
参　　数：
  -u   设置用户的quota，这是预设的参数。
  -g   设置群组的quota。
  -p<源用户名称>   将源用户的quota设置套用至其他用户或群组。
  -t   设置宽限期限。
### eject：退出抽取式设备。
语　　法：eject [-dfhnqrstv][-a <开关>][-c <光驱编号>][设备]
补充说明：若设备已挂入，则eject会先将该设备卸除再退出。
参　　数：
  [设备]   设备可以是驱动程序名称，也可以是挂入点。
  -a<开关>或--auto<开关>   控制设备的自动退出功能。
  -c<光驱编号>或--changerslut<光驱编号>   选择光驱柜中的光驱。
  -d或--default   显示预设的设备，而不是实际执行动作。
  -f或--floppy   退出抽取式磁盘。
  -h或--help   显示帮助。
  -n或--noop   显示指定的设备。
  -q或--tape   退出磁带。
  -r或--cdrom   退出光盘。
  -s或--scsi   以SCSI指令来退出设备。
  -t或--trayclose   关闭光盘的托盘。
  -v或--verbose   执行时，显示详细的说明。
### lndir(link directory)：连接目录内容。
语　　法：lndir [-ignorelinks][-silent][源目录][目的目录]
补充说明：执行lndir指令，可一口气把源目录底下的文件和子目录统统建立起相互对应的符号连接。
参　　数：
  -ignorelinks   直接建立符号连接的符号连接。
  -silent   不显示指令执行过程。
### ls(list)：列出目录内容。
语　　法：ls [-1aAbBcCdDfFgGhHiklLmnNopqQrRsStuUvxX][-I <范本样式>][-T <跳格字数>][-w <每列字符数>][--block-size=<区块大小>][--color=<使用时机>][--format=<列表格式>][--full-time][--help][--indicator-style=<标注样式>][--quoting-style=<引号样式>][--show-control-chars][--sort=<排序方式>][--time=<时间戳记>][--version][文件或目录...]
补充说明：执行ls指令可列出目录的内容，包括文件和子目录的名称。
参　　数：
  -1   每列仅显示一个文件或目录名称。
  -a或--all   下所有文件和目录。
  -A或--almost-all   显示所有文件和目录，但不显示现行目录和上层目录。
  -b或--escape   显示脱离字符。
  -B或--ignore-backups   忽略备份文件和目录。
  -c   以更改时间排序，显示文件和目录。
  -C   以又上至下，从左到右的直行方式显示文件和目录名称。
  -d或--directory   显示目录名称而非其内容。
  -D或--dired   用Emacs的模式产生文件和目录列表。
  -f   此参数的效果和同时指定"aU"参数相同，并关闭"lst"参数的效果。
  -F或--classify   在执行文件，目录，Socket，符号连接，管道名称后面，各自加上"*","/","=","@","|"号。
  -g   次参数将忽略不予处理。
  -G或--no-group   不显示群组名称。
  -h或--human-readable   用"K","M","G"来显示文件和目录的大小。
  -H或--si   此参数的效果和指定"-h"参数类似，但计算单位是1000Bytes而非1024Bytes。
  -i或--inode   显示文件和目录的inode编号。
  -I<范本样式>或--ignore=<范本样式>   不显示符合范本样式的文件或目录名称。
  -k或--kilobytes   此参数的效果和指定"block-size=1024"参数相同。
  -l   使用详细格式列表。
  -L或--dereference   如遇到性质为符号连接的文件或目录，直接列出该连接所指向的原始文件或目录。
  -m   用","号区隔每个文件和目录的名称。
  -n或--numeric-uid-gid   以用户识别码和群组识别码替代其名称。
  -N或--literal   直接列出文件和目录名称，包括控制字符。
  -o   此参数的效果和指定"-l" 参数类似，但不列出群组名称或识别码。
  -p或--file-type   此参数的效果和指定"-F"参数类似，但不会在执行文件名称后面加上"*"号。
  -q或--hide-control-chars   用"?"号取代控制字符，列出文件和目录名称。
  -Q或--quote-name   把文件和目录名称以""号标示起来。
  -r或--reverse   反向排序。
  -R或--recursive   递归处理，将指定目录下的所有文件及子目录一并处理。
  -s或--size   显示文件和目录的大小，以区块为单位。
  -S   用文件和目录的大小排序。
  -t   用文件和目录的更改时间排序。
  -T<跳格字符>或--tabsize=<跳格字数>   设置跳格字符所对应的空白字符数。
  -u   以最后存取时间排序，显示文件和目录。
  -U   列出文件和目录名称时不予排序。
  -v   文件和目录的名称列表以版本进行排序。
  -w<每列字符数>或--width=<每列字符数>   设置每列的最大字符数。
  -x   以从左到右，由上至下的横列方式显示文件和目录名称。
  -X   以文件和目录的最后一个扩展名排序。
  --block-size=<区块大小>   指定存放文件的区块大小。
  --color=<列表格式>   培植文件和目录的列表格式。
  --full-time   列出完整的日期与时间。
  --help   在线帮助。
  --indicator-style=<标注样式>   在文件和目录等名称后面加上标注，易于辨识该名称所属的类型。
  --quoting-syte=<引号样式>   把文件和目录名称以指定的引号样式标示起来。
  --show-control-chars   在文件和目录列表时，使用控制字符。
  --sort=<排序方式>   配置文件和目录列表的排序方式。
  --time=<时间戳记>   用指定的时间戳记取代更改时间。
  --version   显示版本信息。
### mcd：在MS-DOS文件系统中切换工作目录
语　　法：mcd [目录名称]
补充说明：mcd为mtools工具指令，可在MS-DOS文件系统中切换工作目录。若不加任何参数，则显示目前所在的磁盘与工作目录。
### mdeltree：删除MS-DOS目录
语　　法：mdeltree [目录...]
补充说明：mdel为mtools工具指令，模拟MS-DOS的deltree指令，可删除MS-DOS文件系统中的目录及目录下所有子目录与文件。
### mdu：显示MS-DOS目录所占用的磁盘空间。
语　　法：mdu [-as][目录]
补充说明：mdu为mstools工具指令，可显示MS-DOS文件系统中目录所占用的磁盘空间。
参　　数：
  -a   显示每个文件及整个目录所占用的空间。
  -s   仅显示整个目录所占用的空间。
### mkdir(make directories)：建立目录
语　　法：mkdir [-p][--help][--version][-m <目录属性>][目录名称]
补充说明：mkdir可建立目录并同时设置目录的权限。
参　　数：
  -m<目录属性>或--mode<目录属性>   建立目录时同时设置目录的权限。
  -p或--parents   若所要建立目录的上层目录目前尚未建立，则会一并建立上层目录。
  --help   显示帮助。
  --verbose   执行时显示详细的信息。
  --version   显示版本信息。
### mlabel：显示或设置MS-DOS磁盘驱动器的标签名称。
语　　法：mlabel [-csvV][驱动器代号][标签名称]
补充说明：mlabel为mtools工具指令，模拟MS-DOS的label指令，可显示或设置MS-DOS磁盘驱动器的标签名称。
参　　数：
  -c   清楚标签名称
  -s   显示标签名称
  -v   执行时显示详细的信息。
  -V   显示版本信息。
### mmd：在MS-DOS文件系统中建立目录。
语　　法：mmd [目录...]
补充说明：mmd为mtools工具指令，模拟MS-DOS的md指令，可在MS-DOS的文件系统中建立目录。
### mmount：挂入MS-DOS文件系统。
语　　法：mmount [驱动器代号][mount参数]
补充说明：mmount为mtools工具指令，可根据[mount参数]中的设置，将磁盘内容挂入到Linux目录中。
参　　数：
  [mount参数]的用法请参考mount指令。
### mrd：删除MS-DOS文件系统中的目录。
语　　法：mrd [目录...]
补充说明：mrd为mtools工具指令，模拟MS-DOS的rd指令，可删除MS-DOS的目录。
### mzip：Zip/Jaz磁盘驱动器控制指令。
语　　法：mzip [-efpqruwx]
补充说明：mzip为mtools工具指令，可设置Zip或Jaz磁盘驱动区的保护模式以及执行退出磁盘的动作。
参　　数：
  -e   退出磁盘。
  -f   与-e参数一并使用，不管是否已经挂入磁盘中的文件系统，一律强制退出磁盘。
  -p   设置磁盘的写入密码。
  -q   显示目前的状态。
  -r   将磁盘设为防写状态。
  -u   退出磁盘以前，暂时解除磁盘的保护状态。
  -w   将磁盘设为可写入状态。
  -x   设置磁盘的密码。
### pwd(print working directory)：显示工作目录。
语　　法：pwd [--help][--version]
补充说明：执行pwd指令可立刻得知您目前所在的工作目录的绝对路径名称。
参　　数：
  --help   在线帮助。
  --version   显示版本信息。
### quota：显示磁盘已使用的空间与限制。
语　　法：quota [-quvV][用户名称...] 或 quota [-gqvV][群组名称...]
补充说明：执行quota指令，可查询磁盘空间的限制，并得知已使用多少空间。
参　　数：
  -g   列出群组的磁盘空间限制。
  -q   简明列表，只列出超过限制的部分。
  -u   列出用户的磁盘空间限制。
  -v   显示该用户或群组，在所有挂入系统的存储设备的空间限制。
  -V   显示版本信息。
### quotacheck：检查磁盘的使用空间与限制。
语　　法：quotacheck [-adgRuv][文件系统...]
补充说明：执行quotacheck指令，扫描挂入系统的分区，并在各分区的文件系统根目录下产生quota.user和quota.group文件，设置用户和群组的磁盘空间限制。
参　　数：
  -a   扫描在/etc/fstab文件里，有加入quota设置的分区。
  -d   详细显示指令执行过程，便于排错或了解程序执行的情形。
  -g   扫描磁盘空间时，计算每个群组识别码所占用的目录和文件数目。
  -R   排除根目录所在的分区。
  -u   扫描磁盘空间时，计算每个用户识别码所占用的目录和文件数目。
  -v   显示指令执行过程。
### quotaoff：关闭磁盘空间限制。
语　　法：quotaoff [-aguv][文件系统...]
补充说明：执行quotaoff指令可关闭用户和群组的磁盘空间限制。
参　　数：
  -a   关闭在/etc/fstab文件里，有加入quota设置的分区的空间限制。
  -g   关闭群组的磁盘空间限制。
  -u   关闭用户的磁盘空间限制。
  -v   显示指令执行过程。
### quotaon：开启磁盘空间限制。
语　　法：quotaon [-aguv][文件系统...]
补充说明：执行quotaon指令可开启用户和群组的才磅秒年空间限制，各分区的文件系统根目录必须有quota.user和quota.group配置文件。
参　　数：
  -a   开启在/ect/fstab文件里，有加入quota设置的分区的空间限制。
  -g   开启群组的磁盘空间限制。
  -u   开启用户的磁盘空间限制。
  -v   显示指令指令执行过程。
### repquota(report quota)：检查磁盘空间限制的状态。
语　　法：repquota [-aguv][文件系统...]
补充说明：执行repquota指令，可报告磁盘空间限制的状况，清楚得知每位用户或每个群组已使用多少空间。
参　　数：
  -a   列出在/etc/fstab文件里，有加入quota设置的分区的使用状况，包括用户和群组。
  -g   列出所有群组的磁盘空间限制。
  -u   列出所有用户的磁盘空间限制。
  -v   显示该用户或群组的所有空间限制。
### rmdir(remove directory)：删除目录。
语　　法：rmdir [-p][--help][--ignore-fail-on-non-empty][--verbose][--version][目录...]
补充说明：当有空目录要删除时，可使用rmdir指令。
参　　数：
  -p或--parents   删除指定目录后，若该目录的上层目录已变成空目录，则将其一并删除。
  --help   在线帮助。
  --ignore-fail-on-non-empty   忽略非空目录的错误信息。
  --verbose   显示指令执行过程。
  --version   显示版本信息。
### rmt(remote magnetic tape)：远端磁带传输协议模块。
语　　法：rmt
补充说明：通过rmt指令，用户可通过IPC连线，远端操控磁带机的倾倒和还原操作。
### stat(status)：显示inode内容。
语　　法：stat [文件或目录]
补充说明：stat以文字的格式来显示inode的内容。
### Tree：以树状图列出目录的内容。
语　　法：tree [-aACdDfFgilnNpqstux][-I <范本样式>][-P <范本样式>][目录...]
补充说明：执行tree指令，它会列出指定目录下的所有文件，包括子目录里的文件。
参　　数：
  -a   显示所有文件和目录。
  -A   使用ASNI绘图字符显示树状图而非以ASCII字符组合。
  -C   在文件和目录清单加上色彩，便于区分各种类型。
  -d   显示目录名称而非内容。
  -D   列出文件或目录的更改时间。
  -f   在每个文件或目录之前，显示完整的相对路径名称。
  -F   在执行文件，目录，Socket，符号连接，管道名称名称，各自加上"*","/","=","@","|"号。
  -g   列出文件或目录的所属群组名称，没有对应的名称时，则显示群组识别码。
  -i   不以阶梯状列出文件或目录名称。
  -I<范本样式>   不显示符合范本样式的文件或目录名称。
  -l   如遇到性质为符号连接的目录，直接列出该连接所指向的原始目录。
  -n   不在文件和目录清单加上色彩。
  -N   直接列出文件和目录名称，包括控制字符。
  -p   列出权限标示。
  -P<范本样式>   只显示符合范本样式的文件或目录名称。
  -q   用"?"号取代控制字符，列出文件和目录名称。
  -s   列出文件或目录大小。
  -t   用文件和目录的更改时间排序。
  -u   列出文件或目录的拥有者名称，没有对应的名称时，则显示用户识别码。
  -x   将范围局限在现行的文件系统中，若指定目录下的某些子目录，其存放于另一个文件系统上，则将该子目录予以排除在寻找范围外。
### umount：卸除文件系统。
语　　法：umount [-ahnrvV][-t <文件系统类型>][文件系统]
补充说明：umount可卸除目前挂在Linux目录中的文件系统。
参　　数：
  -a   卸除/etc/mtab中记录的所有文件系统。
  -h   显示帮助。
  -n   卸除时不要将信息存入/etc/mtab文件中。
  -r   若无法成功卸除，则尝试以只读的方式重新挂入文件系统。
  -t<文件系统类型>   仅卸除选项中所指定的文件系统。
  -v   执行时显示详细的信息。
  -V   显示版本信息。
  [文件系统]   除了直接指定文件系统外，也可以用设备名称或挂入点来表示文件系统。

 
## linux磁盘维护命令
### badblocks：检查磁盘装置中损坏的区块。
语　　法：badblocks [-svw][-b <区块大小>][-o <输出文件>][磁盘装置][磁盘区块数][启始区块]
补充说明：执行指令时须指定所要检查的磁盘装置，及此装置的磁盘区块数。
参　　数：
  -b<区块大小>   指定磁盘的区块大小，单位为字节。
  -o<输出文件>   将检查的结果写入指定的输出文件。
  -s   在检查时显示进度。
  -v   执行时显示详细的信息。
  -w   在检查时，执行写入测试。
  [磁盘装置]   指定要检查的磁盘装置。
  [磁盘区块数]   指定磁盘装置的区块总数。
  [启始区块]   指定要从哪个区块开始检查。
### cfdisk：磁盘分区。
语　　法：cfdisk [-avz][-c <柱面数目>-h <磁头数目>-s <盘区数目>][-P <r,s,t>][外围设备代号]
补充说明：cfdisk是用来磁盘分区的程序，它十分类似DOS的fdisk，具有互动式操作界面而非传统fdisk的问答式界面，您可以轻易地利用方向键来操控分区操作。
参　　数：
  -a   在程序里不用反白代表选取，而以箭头表示。
  -c<柱面数目>   忽略BIOS的数值，直接指定磁盘的柱面数目。
  -h<磁头数目>   忽略BIOS的数值，直接指定磁盘的磁头数目。
  -P<r,s,t>   显示分区表的内容，附加参数"r"会显示整个分区表的详细资料，附加参数"s"会依照磁区的顺序显示相关信息，附加参数"t"则会以磁头，磁区，柱面的方式来显示资料。
  -s<磁区数目>   忽略BIOS的数值，直接指定磁盘的磁区数目。
  -v   显示版本信息。
  -z   不读取现有的分区，直接当作没有分区的新磁盘使用。
### dd：读取，转换并输出数据。
语　　法：dd [bs=<字节数>][cbs=<字节数>][conv=<关键字>][count=<区块数>][ibs=<字节数>][if=<文件>][obs=<字节数>][of=<文件>][seek=<区块数>][skip=<区块数>][--help][--version]
补充说明：dd可从标准输入或文件读取数据，依指定的格式来转换数据，再输出到文件，设备或标准输出。
参　　数：
  bs=<字节数>   将ibs( 输入)与obs(输出)设成指定的字节数。
  cbs=<字节数>   转换时，每次只转换指定的字节数。
  conv=<关键字>   指定文件转换的方式。
  count=<区块数>   仅读取指定的区块数。
  ibs=<字节数>   每次读取的字节数。
  if=<文件>   从文件读取。
  obs=<字节数>   每次输出的字节数。
  of=<文件>   输出到文件。
  seek=<区块数>   一开始输出时，跳过指定的区块数。
  skip=<区块数>   一开始读取时，跳过指定的区块数。
  --help   帮助。
  --version   显示版本信息。
### e2fsck(ext2 file system check)：检查ext2文件系统的正确性。
语　　法：e2fsck [-acCdfFnprsStvVy][-b <superblock>][-B <区块大小>][-l <文件>][-L <文件>][设备名称]
补充说明：e2fsck执行后的传回值及代表意义如下。
  0     没有任何错误发生。
  1     文件系统发生错误，并且已经修正。
  2     文件系统发生错误，并且已经修正。
  4     文件系统发生错误，但没有修正。
  8     运作时发生错误。
  16    使用的语法发生错误。
  128   共享的函数库发生错误。
参　　数：
  -a   不询问使用者意见，便自动修复文件系统。
  -b<superblock>   指定superblock，而不使用预设的superblock。
  -B<区块大小>   指定区块的大小，单位为字节。
  -c   一并执行badblocks，以标示损坏的区块。
  -C   将检查过程的信息完整记录在file descriptor中，使得整个检查过程都能完整监控。
  -d   显示排错信息。
  -f   即使文件系统没有错误迹象，仍强制地检查正确性。
  -F   执行前先清除设备的缓冲区。
  -l<文件>   将文件中指定的区块加到损坏区块列表。
  -L<文件>   先清除损坏区块列表，再将文件中指定的区块加到损坏区块列表。因此损坏区块列表的区块跟文件中指定的区块是一样的。
  -n   以只读模式开启文件系统，并采取非互动方式执行，所有的问题对话均设置以"no"回答。
  -p   不询问使用者意见，便自动修复文件系统。
  -r   此参数只为了兼容性而存在，并无实际作用。
  -s   如果文件系统的字节顺序不适当，就交换字节顺序，否则不做任何动作。
  -S   不管文件系统的字节顺序，一律交换字节顺序。
  -t   显示时间信息。
  -v   执行时显示详细的信息。
  -V   显示版本信息。
  -y   采取非互动方式执行，所有的问题均设置以"yes"回答。
### ext2ed(ext2 file system editor)：ext2文件系统编辑程序。
语　　法：ext2ed
补充说明：ext2ed可直接处理硬盘分区上的数据，这指令只有Red Hat Linux才提供。
参　　数：
  一般指令
  setdevice[设备名称]   指定要处理的设备。
  disablewrite   将ext2ed设为只读的状态。
  enablewrite   将ext2ed设为可读写的状态。
  help[指令]   显示个别指令的帮助。
  next   移至下一个单位，单位会依目前所在的模式而异。
  prev   移至前一个单位，单位会依目前所在的模式而异。
  pgup   移至下一页。
  pgdn   移至上一页。
  set   修改目前的数据，参数会依目前所在的模式而异。
  writedata   在执行此指令之后，才会实际修改分区中的数据。
  ext2进入3种模式的指令
  super   进入main superblock,即Superblock模式。
  group<编号>   进入指定的group，即Group模式。
  cd<目录或文件>   在inode模式下，进入指定的目录或文件，即Inode模式。
  Superblock模式
  gocopy<备份编号>   进入指定的superblock备份。
  setactivecopy   将目前所在的superblock，复制到main superblock。
  Group模式
  blockbitmap   显示目前groupo的区块图。
  inode   进入目前group的第一个inode。
  inodebitmap   显示目前group的inode二进制码。
  Inode模式
  dir   进入目录模式。
  file   进入文件模式。
### fdisk：磁盘分区。
语　　法：fdisk [-b <分区大小>][-uv][外围设备代号] 或 fdisk [-l][-b <分区大小>][-uv][外围设备代号...] 或 fdisk [-s <分区编号>]
补充说明：fdisk是用来磁盘分区的程序，它采用传统的问答式界面，而非类似DOS fdisk的cfdisk互动式操作界面，因此在使用上较为不便，但功能却丝毫不打折扣。
参　　数：
  -b<分区大小>   指定每个分区的大小。
  -l   列出指定的外围设备的分区表状况。
  -s<分区编号>   将指定的分区大小输出到标准输出上，单位为区块。
  -u   搭配"-l"参数列表，会用分区数目取代柱面数目，来表示每个分区的起始地址。
  -v   显示版本信息。
### fsck.ext2(file system check-second filesystem)：检查文件系统并尝试修复错误。
语　　法：fsck.ext2 [-acdfFnprsStvVy][-b <分区第一个磁区地址>][-B <区块大小>][-C <反叙述器>][-I <inode缓冲区块数>][-l/L <损坏区块文件>][-P <处理inode大小>][外围设备代号]
补充说明：当ext2文件系统发生错误时，可用fsck.ext2指令尝试加以修复。
参　　数：
  -a   自动修复文件系统，不询问任何问题。
  -b<分区第一个磁区地址>   指定分区的第一个磁区的起始地址，也就是Super Block。
  -B<区块大小>   设置该分区每个区块的大小。
  -c   检查指定的文件系统内，是否存在有损坏的区块。
  -C<反叙述器>   指定反叙述器，fsck.ext2指令会把全部的执行过程，都交由其逆向叙述，便于排错或监控程序执行的情形。
  -d   详细显示指令执行过程，便于排错或分析程序执行的情形。
  -f   强制对该文件系统进行完整检查，纵然该文件系统在慨略检查下没有问题。
  -F   检查文件系统之前，先清理该保存设备块区内的数据。
  -I<inode缓冲区块数>   设置欲检查的文件系统，其inode缓冲区的区块数目。
  -l<损坏区块文件>   把文件中所列出的区块，视为损坏区块并将其标示出来，避免应用程序使用该区块。
  -L<损坏区块文件>   此参数的效果和指定"-l"参数类似，但在参考损坏区块文件标示损坏区块之前，会先将原来标示成损坏区块者统统清楚，即全部重新设置，而非仅是加入新的损坏区块标示。
  -n   把欲检查的文件系统设成只读，并关闭互动模式，否决所有询问的问题。
  -p   此参数的效果和指定"-a"参数相同。
  -P<处理inode大小>   设置fsck.ext2指令所能处理的inode大小为多少。
  -r   此参数将忽略不予处理，仅负责解决兼容性的问题。
  -s   检查文件系统时，交换每对字节的内容。
  -S   此参数的效果和指定"-s"参数类似，但不论该文件系统是否已是标准位顺序，一律交换每对字节的内容。
  -t   显示fsck.ext2指令的时序信息。
  -v   详细显示指令执行过程。
  -V   显示版本信息。
  -y   关闭互动模式，且同意所有询问的问题。
### fsck(file system check)：检查文件系统并尝试修复错误。
语　　法：fsck [-aANPrRsTV][-t <文件系统类型>][文件系统...]
补充说明：当文件系统发生错误四化，可用fsck指令尝试加以修复。
参　　数：
  -a   自动修复文件系统，不询问任何问题。
  -A   依照/etc/fstab配置文件的内容，检查文件内所列的全部文件系统。
  -N   不执行指令，仅列出实际执行会进行的动作。
  -P   当搭配"-A"参数使用时，则会同时检查所有的文件系统。
  -r   采用互动模式，在执行修复时询问问题，让用户得以确认并决定处理方式。
  -R   当搭配"-A"参数使用时，则会略过/目录的文件系统不予检查。
  -s   依序执行检查作业，而非同时执行。
  -t<文件系统类型>   指定要检查的文件系统类型。
  -T   执行fsck指令时，不显示标题信息。
  -V   显示指令执行过程。
### fsck.minix(file system check-minix filesystem)：检查文件系统并尝试修复错误。
语　　法：fsck.minix [-aflmrsv][外围设备代号]
补充说明：当minix文件系统发生错误时，可用fsck.minix指令尝试加以参考。
参　　数：
  -a   自动修复文件系统，不询问任何问题。
  -f   强制对该文件系统进行完整检查，纵然该文件系统在慨略检查下没有问题。
  -l   列出所有文件名称。
  -m   使用类似MINIX操作系统的警告信息。
  -r   采用互动模式，在执行修复时询问问题，让用户得以确认并决定处理方式。
  -s   显示该分区第一个磁区的相关信息。
  -v   显示指令执行过程。
### fsconf(file system configurator)：设置文件系统相关功能。
语　　法：fsconf [--check]
补充说明：fsconf是Red Hat Linux发行版专门用来调整Linux各项设置的程序。
参　　数：
  --chedk   检查特定文件的权限。
### hdparm(hard disk parameters)：显示与设定硬盘的参数。
语　　法：hdparm [-CfghiIqtTvyYZ][-a <快取分区>][-A <0或1>][-c <I/O模式>][-d <0或1>][-k <0或1>][-K <0或1>][-m <分区数>][-n <0或1>][-p <PIO模式>][-P <分区数>][-r <0或1>][-S <时间>][-u <0或1>][-W <0或1>][-X <传输模式>][设备]
补充说明：hdparm可检测，显示与设定IDE或SCSI硬盘的参数。
参　　数：
  -a<快取分区>   设定读取文件时，预先存入块区的分区数，若不加上<快取分区>选项，则显示目前的设定。
  -A<0或1>   启动或关闭读取文件时的快取功能。
  -c<I/O模式>   设定IDE32位I/O模式。
  -C   检测IDE硬盘的电源管理模式。
  -d<0或1>   设定磁盘的DMA模式。
  -f   将内存缓冲区的数据写入硬盘，并清楚缓冲区。
  -g   显示硬盘的磁轨，磁头，磁区等参数。
  -h   显示帮助。
  -i   显示硬盘的硬件规格信息，这些信息是在开机时由硬盘本身所提供。
  -I   直接读取硬盘所提供的硬件规格信息。
  -k<0或1>   重设硬盘时，保留-dmu参数的设定。
  -K<0或1>   重设硬盘时，保留-APSWXZ参数的设定。
  -m<磁区数>   设定硬盘多重分区存取的分区数。
  -n<0或1>   忽略硬盘写入时所发生的错误。
  -p<PIO模式>   设定硬盘的PIO模式。
  -P<磁区数>   设定硬盘内部快取的分区数。
  -q   在执行后续的参数时，不在屏幕上显示任何信息。
  -r<0或1>   设定硬盘的读写模式。
  -S<时间>   设定硬盘进入省电模式前的等待时间。
  -t   评估硬盘的读取效率。
  -T   平谷硬盘快取的读取效率。
  -u<0或1>   在硬盘存取时，允许其他中断要求同时执行。
  -v   显示硬盘的相关设定。
  -W<0或1>   设定硬盘的写入快取。
  -X<传输模式>   设定硬盘的传输模式。
  -y   使IDE硬盘进入省电模式。
  -Y   使IDE硬盘进入睡眠模式。
  -Z   关闭某些Seagate硬盘的自动省电功能。
### losetup(loop setup)：设置循环设备。
语　　法：losetup [-d][-e <加密方式>][-o <平移数目>][循环设备代号][文件]
补充说明：循环设备可把文件虚拟成区块设备，籍以模拟整个文件系统，让用户得以将其视为硬盘驱动器，光驱或软驱等设备，并挂入当作目录来使用。
参　　数：
  -d   卸除设备。
  -e<加密方式>   启动加密编码。
  -o<平移数目>   设置数据平移的数目。
### mbadblocks：检查MS-DOS文件系统的磁盘是否有损坏的磁区。
语　　法：mbadblocks [驱动器代号]
补充说明：mbadblocks为mtools工具指令，可用来扫描MS-DOS文件系统的磁盘驱动器，并标示出损坏的磁区。
### mformat：对MS-DOS文件系统的磁盘进行格式化。
语　　法：mformat [-1aCFIKX][-0 <数据传输率>][-2 <磁区数>][-A <数据传输率>][-B <开机区文件>][-c <丛集大小>][-h <磁头数>][-H <隐藏磁区数>][-l <磁盘标签>][-M <磁区大小>][-n <序号>][-r <根目录大小>][-s <磁区数>][-S <磁区大小>][-t <柱面数>][驱动器代号]
补充说明：mformat为mtools工具指令，模拟MS-DOS的format指令，可将指定的磁盘或硬盘分区格式化为MS-DOS文件系统。
参　　数：
  -0<数据传输率>   指定第0磁轨的数据传输率。
  -1   不使用2M格式。
  -2<磁区数>   指定在第0磁轨的第0个磁头的磁区数，也就是所谓的2M格式。
  -a   指定产生Atari格式的序号。
  -A<数据传输率>   指定第0磁轨以外所有磁轨的数据传输率。
  -B<开机区文件>   从指定的文件来建立开机区。
  -c<丛集大小>   指定丛集大小，单位为磁区。
  -C   建立磁盘MS-DOS文件系统的映像文件。
  -F   以FAT32来格式化磁盘。
  -H<隐藏磁区数>   指定隐藏磁区的数目。
  -I   指定FAT32的frVersion编号。
  -k   尽可能地保留原有的开机区。
  -l<磁盘标签>   指定磁盘标签名称。
  -M<磁区大小>   指定MS-DOS文件系统所使用的磁区大小，预设应该与硬体磁区大小相同，单位为字节。
  -n<序号>   指定序号。
  -r<根目录大小>   指定根目录的大小，单位为磁区。
  -s<磁区数>   指定每一磁轨所包含的磁区数目。
  -S<磁区大小>   指定硬件磁区大小。
  -t<柱面数>   指定柱面数目。
  -X   格式化XDF磁盘。
### mkbootdisk(make boot disk)：建立目前系统的启动盘。
语　　法：mkbootdisk [--noprompt][--verbose][--version][--device <设备>][--mkinitrdargs <参数>][kernel 版本]
补充说明：mkbootdisk可建立目前系统的启动盘。
参　　数：
  --device<设备>   指定设备。
  --mkinitrdargs<参数>   设置mkinitrd的参数。
  --noprompt   不会提示用户插入磁盘。
  --verbose   执行时显示详细的信息。
  --version   显示版本信息。
### mkdosfs(make Dos file system)：建立MS-DOS文件系统。
语　　法：mkdosfs [-cv][-f <FAT数目>][-F <FAT记录的单位>][-i <磁盘序号>][-l <文件名>][-m <信息文件>][-n <标签>][-r <根目录项目数>][-s <磁区数>][设备名称][区块数]
补充说明：mkdosfs可在Linux下，将磁盘格式化为MS-DOS文件系统的格式。
参　　数：
  -c   检查是否有损坏的区块。
  -f<FAT数目>   指定FAT的数目，目前支持1与2两种选项。
  -F<FAT记录的单位>   指定FAT记录的单位为12或16位。
  -i<磁盘序号>   设置文件系统的磁盘序号。
  -l<文件名>   从指定的文件中，读取文件系统中损坏区块的信息。
  -m<信息文件>   若以次文件系统开机，而却因操作系统没有正常安装，则会显示选项中所指定的信息文件内容。
  -n<标签>   设置文件系统的磁盘标签名称。
  -r<根目录项目数>   设置根目录最多能记载项目数。
  -s<磁区数>   指定每个丛集所包含的磁区数。
  -v   执行时显示详细的信息。
### mke2fs(make ext2 file system)：建立ext2文件系统。
语　　法：mke2fs [-cFMqrSvV][-b <区块大小>][-f <不连续区段大小>][-i <字节>][-N <inode数>][-l <文件>][-L <标签>][-m <百分比值>][-R=<区块数>][ 设备名称][区块数]
补充说明：mke2fs可建立Linux的ext2文件系统。
参　　数：
  -b<区块大小>   指定区块大小，单位为字节。
  -c   检查是否有损坏的区块。
  -f<不连续区段大小>   指定不连续区段的大小，单位为字节。
  -F   不管指定的设备为何，强制执行mke2fs。
  -i<字节>   指定"字节/inode"的比例。
  -N<inode数>   指定要建立的inode数目。
  -l<文件>   从指定的文件中，读取文件西中损坏区块的信息。
  -L<标签>   设置文件系统的标签名称。
  -m<百分比值>   指定给管理员保留区块的比例，预设为5%。
  -M   记录最后一次挂入的目录。
  -q   执行时不显示任何信息。
  -r   指定要建立的ext2文件系统版本。
  -R=<区块数>   设置磁盘阵列参数。
  -S   仅写入superblock与group descriptors，而不更改inode able inode bitmap以及block bitmap。
  -v   执行时显示详细信息。
  -V   显示版本信息。
### mkfs.ext2：与mke2fs相同。
### mkfs(make file system)：建立各种文件系统。
语　　法：mkfs [-vV][fs][-f <文件系统类型>][设备名称][区块数]
补充说明：mkfs本身并不执行建立文件系统的工作，而是去调用相关的程序来执行。
参　　数：
  fs   指定建立文件系统时的参数。
  -t<文件系统类型>   指定要建立何种文件系统。
  -v   显示版本信息与详细的使用方法。
  -V   显示简要的使用方法。
### mkfs.minix：建立Minix文件系统。
语　　法：mkfs.minix [-cv][-i <inode数目>][-l <文件>][-n <文件名长度>][设备名称][区块数]
补充说明：mkfs.minix可建立Minix文件系统。
参　　数：
  -c   检查是否有损坏的区块。
  -i<inode数目>   指定文件系统的inode总数。
  -l<文件>   从指定的文件中，读取文件系统中损坏区块的信息。
  -n<文件名长度>   指定文件名称长度的上限。
  -v   建立第2版的Minix文件系统。
### mkfs.msdos：与mkdosfs相同。
### mkinitrd(make initial ramdisk images)：建立要载入ramdisk的映像文件。
语　　法：mkinitrd [-fv][--omit-scsi-modules][--version][--preload=<模块名称>][--with=<模块名称>][映像文件][Kernel 版本]
补充说明：mkinitrd可建立映像文件，以供Linux开机时载入ramdisk。
参　　数：
  -f   若指定的映像问家名称与现有文件重复，则覆盖现有的文件。
　-v　　执行时显示详细的信息。
  --omit-scsi-modules   不要载入SCSI模块。
  --preload=<模块名称>   指定要载入的模块。
  --with=<模块名称>   指定要载入的模块。
  --version   显示版本信息。
### mkisofs(make iso file system)：建立ISO 9660映像文件。
语　　法：mkisofs [-adDfhJlLNrRTvz][-print-size][-quiet][-A <应用程序ID>][-abstract <摘要文件>][-b <开机映像文件>][-biblio <ISBN文件>][-c <开机文件名称>][-C <盘区编号，磁区编号>][-copyright <版权信息文件>][-hide <目录或文件名>][-hide-joliet <文件或目录名>][-log-file <记录文件>][-m <目录或文件名>][-M <开机映像文件>][-o <映像文件>][-p <数据处理人>][-P <光盘发行人>][-sysid <系统ID >][-V <光盘ID >][-volset <卷册集ID>][-volset-size <光盘总数>][-volset-seqno <卷册序号>][-x <目录>][目录或文件]
补充说明：mkisofs可将指定的目录与文件做成ISO 9660格式的映像文件，以供刻录光盘。
参　　数：
  -a或--all   mkisofs通常不处理备份文件。使用此参数可以把备份文件加到映像文件中。
  -A<应用程序ID>或-appid<应用程序ID>   指定光盘的应用程序ID。
  -abstract<摘要文件>   指定摘要文件的文件名。
  -b<开机映像文件>或-eltorito-boot<开机映像文件>   指定在制作可开机光盘时所需的开机映像文件。
  -biblio<ISBN文件>   指定ISBN文件的文件名，ISBN文件位于光盘根目录下，记录光盘的ISBN。
  -c<开机文件名称>   制作可开机光盘时，mkisofs会将开机映像文件中的全-eltorito-catalog<开机文件名称>全部内容作成一个文件。
  -C<盘区编号，盘区编号>   将许多节区合成一个映像文件时，必须使用此参数。
  -copyright<版权信息文件>   指定版权信息文件的文件名。
  -d或-omit-period   省略文件后的句号。
  -D或-disable-deep-relocation   ISO 9660最多只能处理8层的目录，超过8层的部分，RRIP会自动将它们设置成ISO 9660兼容的格式。使用-D参数可关闭此功能。
  -f或-follow-links   忽略符号连接。
  -h   显示帮助。
  -hide<目录或文件名>   使指定的目录或文件在ISO 9660或Rock RidgeExtensions的系统中隐藏。
  -hide-joliet<目录或文件名>   使指定的目录或文件在Joliet系统中隐藏。
  -J或-joliet   使用Joliet格式的目录与文件名称。
  -l或-full-iso9660-filenames   使用ISO 9660 32字符长度的文件名。
  -L或-allow-leading-dots   允许文件名的第一个字符为句号。
  -log-file<记录文件>   在执行过程中若有错误信息，预设会显示在屏幕上。
  -m<目录或文件名>或-exclude<目录或文件名>   指定的目录或文件名将不会房入映像文件中。
  -M<映像文件>或-prev-session<映像文件>   与指定的映像文件合并。
  -N或-omit-version-number   省略ISO 9660文件中的版本信息。
  -o<映像文件>或-output<映像文件>   指定映像文件的名称。
  -p<数据处理人>或-preparer<数据处理人>   记录光盘的数据处理人。
  -print-size   显示预估的文件系统大小。
  -quiet   执行时不显示任何信息。
  -r或-rational-rock   使用Rock Ridge Extensions，并开放全部文件的读取权限。
  -R或-rock   使用Rock Ridge Extensions。
  -sysid<系统ID>   指定光盘的系统ID。
  -T或-translation-table   建立文件名的转换表，适用于不支持Rock Ridge Extensions的系统上。
  -v或-verbose   执行时显示详细的信息。
  -V<光盘ID>或-volid<光盘ID>   指定光盘的卷册集ID。
  -volset-size<光盘总数>   指定卷册集所包含的光盘张数。
  -volset-seqno<卷册序号>   指定光盘片在卷册集中的编号。
  -x<目录>   指定的目录将不会放入映像文件中。
  -z   建立通透性压缩文件的SUSP记录，此记录目前只在Alpha机器上的Linux有效。
### mkswap：设置交换区(swap area)。
语　　法：mkswap [-cf][-v0][-v1][设备名称或文件][交换区大小]
补充说明：mkswap可将磁盘分区或文件设为Linux的交换区。
参　　数：
  -c   建立交换区前，先检查是否有损坏的区块。
  -f   在SPARC电脑上建立交换区时，要加上此参数。
  -v0   建立旧式交换区，此为预设值。
  -v1   建立新式交换区。
  [交换区大小]   指定交换区的大小，单位为1024字节。
### mpartition：建立或删除MS-DOS的分区。
语　　法：mpartition [-acdfIprv][-b <磁区数>][-h <磁头数>][l <磁区数>][-s <磁区数>][-t <柱面数>][驱动器代号]
补充说明：mpartition为mtools工具指令，可建立或删除磁盘分区。
参　　数：
  -a   将分区设置为可开机分区。
  -b<磁区数>   建立分区时，指定要从第几个磁区开始建立分区。
  -c   建立分区。
  -d   将分区设置为无法开机的分区。
  -f   强制地修改分区而不管检查时发生的错误信息。
  -h<磁头数>   建立分区时，指定分区的磁头数。
  -I   删除全部的分区。
  -l<磁区数>   建立分区时，指定分区的容量大小，单位为磁区数。
  -p   当要重新建立分区时，显示命令列。
  -r   删除分区。
  -s<磁区数>   建立分区时，指定每个磁轨的磁区数。
  -t<柱面数>   建立分区时，指定分区的柱面数。
  -v   与-p参数一并使用，若没有同时下达修改分区的命令，则显示目前分区的状态。
### sfdisk：硬盘恩区工具程序。
语　　法：sfdisk [-?Tvx][-d <硬盘>][-g <硬盘>][-l <硬盘>][-s <分区>][-V <硬盘>]
补充说明：sfdisk为硬盘分区工具程序，可显示分区的设置信息，并检查分区是否正常。
参　　数：
  -?或--help   显示帮助。
  -d<硬盘>   显示硬盘分区的设置。
  -g<硬盘>或--show-geometry<硬盘>   显示硬盘的CHS参数。
  -l<硬盘>   显示后硬盘分区的相关设置。
  -s<分区>   显示分区的大小，单位为区块。
  -T或--list-types   显示所有sfdisk能辨识的文件系统ID。
  -v或--version   显示版本信息。
  -V<硬盘>或--verify<硬盘>   检查硬盘分区是否正常。
  -x或--show-extend   显示扩展分区中的逻辑分区。
### swapoff：关闭系统交换区(swap area)。
语　　法：swapoff [设备]
补充说明：swapoff实际上为swapon的符号连接，可用来关闭系统的交换区。
### swapon：启动系统交换区(swap area)。
语　　法：swapon [-ahsV][-p <优先顺序>][设备]
补充说明：Linux系统的内存管理必须使用交换区来建立虚拟内存。
  -a   将/etc/fstab文件中所有设置为swap的设备，启动为交换区。
  -h   显示帮助。
  -p<优先顺序>   指定交换区的优先顺序。
  -s   显示交换区的使用状况。
  -V   显示版本信息。
### symlinks(symbolic links)：维护符号连接的工具程序。
语　　法：symlinks [-cdrstv][目录]
补充说明：symlinks可检查目录中的符号连接，并显示符号连接类型。以下为symlinks可判断的符号连接类型：
  absolute：符号连接使用了绝对路径。
  dangling：原始文件已经不存在。
  lengthy：符号连接的路径中包含了多余的"../"。
  messy：符号连接的路径中包含了多余的"/"。
  other_fs：原始文件位于其他文件系统中。
  relative：符号连接使用了相对路径。
参　　数：
  -c   将使用绝对路径的符号连接转换为相对路径。
  -d   移除dangling类型的符号连接。
  -r   检查目录下所有子目录中的符号连接。
  -s   检查lengthy类型的符号连接。
  -t   与-c一并使用时，会显示如何将绝对路径的符号连接转换为相对路径，但不会实际转换。
  -v   显示所有类型的符号连接。
### sync：将内存缓冲区内的数据写入磁盘。
语　　法：sync [--help][--version]
补充说明：在Linux系统中，当数据需要存入磁盘时，通常会先放到缓冲区内，等到适当的时刻再写入磁盘，如此可提高系统的执行效率。
参　　数：
  --help   显示帮助。
  --version   显示版本信息。

 
## linux网络通讯命令
### dip：IP拨号连接。
语　　法：dip [-aikltv][-m<MTU数目>][-p<协议>][拨号script文件]
补充说明：dip可控制调制解调器，以拨号IP的方式建立对外的双向连接。
参　　数：
  -a   询问用户名称与密码。
  -i   启动拨号服务器功能。
  -k   删除执行中的dip程序。
  -l   指定要删除的连线，必须配合-k参数一起使用。
  -m<MTU数目>   设置最大传输单位，预设值为296。
  -p<协议>   设置通信协议。
  -t   进入dip的指令模式。
  -v   执行时显示详细的信息。
### getty(get teletypewriter)：设置终端机模式，连线速率和管制线路。
语　　法：getty [-h][-d<组态配置文件>][-r<延迟秒数>][-t<超时秒数>][-w<等待字符串>][终端机编号][连线速率<终端机类型><管制线路>] 或 getty [-c<定义配置文件>]
补充说明：getty指令是UNIX之类操作系统启动时所必须的3个步骤之一。
参　　数：
  -c<定义配置文件>   指定定义配置文件，预设为/etc/gettydefs。
  -d<组态配置文件>   指定组态配置文件，预设为/etc/conf.getty。
  -h   当传输速率为0时就强制断线。
  -r<延迟秒数>   设置延迟时间。
  -t<超时秒数>   设置等待登入的时间。
  -w<等待字符串>   设置等待回应的字符串。
### mingetty：精简版的getty。
语　　法：mingetty [--long-hostname][--noclear][tty]
补充说明：mingetty适用于本机上的登入程序。
参　　数：
  --long-hostname   显示完整的主机名称。
  --noclear   在询问登入的用户名称之前不要清楚屏幕画面。
### ppp-off：关闭ppp连线。
语　　法：ppp-of
补充说明：这是Slackware发行版内附的程序，让用户切断PPP的网络连线。
### smbd(samba daemon)：Samba服务器程序。
语　　法：smbd [-aDhoP][-d<排错层级>][-i<范围>][-l<记录文件>][-O<连接槽选项>][-p<连接端口编号>][-s<配置文件>]
补充说明：smbd为Samba服务器程序，可分享文件与打印机等网络资源供Windows相关的用户端程序存取。
参　　数：
  -a   所有的连线记录都会加到记录文件中。
  -d<排错层级>   指定记录文件所记载事件的详细程度。
  -D   使用此参数时，smbd会以服务程序的方式在后台执行。
  -h   显示帮助。
  -i<范围>   指定NetBIOS名称的范围。
  -l<记录文件>   指定记录文件的名称。
  -o   每次启动时，会覆盖原有的记录文件。
  -O<连接槽选项>   设置连接槽选项。
  -p<连接端口编号>   设置连接端口编号。
  -P   仅用来测试smbd程序的正确性。
  -s<配置文件>   指定smbd的设置文件。
### telnet：远端登入。
语　　法：telnet [-8acdEfFKLrx][-b<主机别名>][-e<脱离字符>][-k<域名>][-l<用户名称>][-n<记录文件>][-S<服务类型>][-X<认证形态>][主机名称或IP地址<通信端口>]
补充说明：执行telnet指令开启终端机阶段作业，并登入远端主机。
参　　数：
  -8   允许使用8位字符资料，包括输入与输出。
  -a   尝试自动登入远端系统。
  -b<主机别名>   使用别名指定远端主机名称。
  -c   不读取用户专属目录里的.telnetrc文件。
  -d   启动排错模式。
  -e<脱离字符>   设置脱离字符。
  -E   滤除脱离字符。
  -f   此参数的效果和指定"-F"参数相同。
  -F   使用Kerberos V5认证时，加上此参数可把本地主机的认证数据上传到远端主机。
  -k<域名>   使用Kerberos认证时，加上此参数让远端主机采用指定的领域名，而非该主机的域名。
  -K   不自动登入远端主机。
  -l<用户名称>   指定要登入远端主机的用户名称。
  -L   允许输出8位字符资料。
  -n<记录文件>   指定文件记录相关信息。
  -r   使用类似rlogin指令的用户界面。
  -S<服务类型>   设置telnet连线所需的IP TOS信息。
  -x   假设主机有支持数据加密的功能，就使用它。
  -X<认证形态>   关闭指定的认证形态。
### uulog：显示UUCP记录文件。
语　　法：uulog [-DFISv][-<行数>][-f<主机>][-I<配置文件>][-n<行数>][-s<主机>][-u<用户>][-X<层级>][--help]
补充说明：uulog可用来显示UUCP记录文件中记录。
参　　数：
  -D或--debuglog   显示排错记录。
  -f<主机>或--follow<主机>   与-F参数类似，但仅显示与指定主机相关的记录。
  -I<配置文件>或--config<配置文件>   指定程序的配置文件。
  -<行数>,-n<行数>或--lines<行数>   显示记录文件中，从最后算起指定行数的数值。
  -s<主机>   仅显示记录文件中，与指定文件相关的记录。
  -S或--statslog   显示统计记录。
  -u<用户>或--suer<用户>   仅显示记录文件中，与指定用户相关的记录。
  -v或--version   显示版本信息。
  -X<层级>或--debug<层级>   设定排错层级。
  --help   显示帮助。
### uustat：显示UUCP目前的状况。
语　　法：uustat [-aeiKmMNpqQRv][-B<行数>][-c<指令>][-C<指令>][-I<配置文件>][-k<工作>][-o<小时>][-r<工作>][-s<主机>][-S<主机>][-u<用户>][-U<用户>][-W<附注>][-y<小时>][-x<层级>][--help]
补充说明：执行uucp与uux指令后，会先将工作送到队列，再由uucico来执行工作。uustat可显示，删除或启动队列中等待执行的工作。
参　　数：
  -a或-all   显示全部的UUCP工作。
  -B<行数>或--mail-lines<行数>   与-M或-N参数一并使用，用来指定邮件中要包含多少行的信息。
  -c<指令>或--command<指令>   显示与<指令>有关的工作。
  -C<指令>或--not-command<指令>   显示与<指令>无关的工作。
  -e或--executions   仅显示待执行的工作。
  -i或--prompt   针对队列中的每项工作，询问使用是否要删除工作。
  -I<配置文件>或--config<配置文件>   指定配置文件。
  -k<工作>或--kill<工作>   删除指定的工作。
  -m或--status   删除全部的工作。
  -M或-mail   将状态信息邮寄给UUCP管理员。
  -N或--notify   将状态信息分别邮寄给提出该项工作的用户。
  -o<小时>或--older-than<小时>   显示超过指定时数的工作。
  -p或--ps   显示负责UUCP锁定的程序。
  -q或--list   显示每台远端主机上所要执行工作的状态。
  -Q或--no-list   不显示工作。
  -r<工作>或--rejuvenate<工作>   重新启动指定的工作。
  -R或--rejuvenate-all   重新启动全部的工作。
  -s<主机>或--system<主机>   显示与<主机>有关的工作。
  -S<主机>或--not-system<主机>   显示与<主机>无关的工作。
  -v或--version   显示版本信息。
  -u<用户>或--user<用户>   显示与<用户>有关的工作。
  -U<用户>或--not-user<用户>   显示与<用户>无关的工作。
  -W<附注>或--comment<附注>   要放在邮件信息中的附注。
  -y<小时>或--younger-than<小时>   显示低于指定时数的工作。
  -x<层级>或--debug<层级>   指定排错层级。
  --help   显示帮助。
### uux：在远端的UUCP主机上执行指令。
语　　法：uux [-bcCIjlnrvz][-a<地址>][-g<等级>][-s<文件>][-x<层级>][--help][指令]
补充说明：uux可在远端的UUCP主机上执行指令或是执行本机上的指令，但在执行时会使用远端电脑的文件。
参　　数：
  -或-p或--stdin   直接从键盘读取要执行的指令。
  -a<地址>或--requestor<地址>   执行邮件地址，以便寄送状态信息。
  -b或--erturn-stdin   在屏幕上显示状态信息。
  -c或--nocopy   不用将文件复制到缓冲区。
  -C或--copy   将文件复制到缓冲区。
  -g<等级>或--grade<等级>   指定文件传送作业的优先顺序。
  -I或--config file   指定uux配置文件。
  -j或--jobid   显示作业编号。
  -l或--link   将本机上的文件连接到缓冲区。
  -n或--notification=no   无论发生任何状态，都不寄邮件通知用户。
  -r或--nouucico   不要立即启动uucico服务程序，仅将作业送到队列中，然后再执行。
  -s<文件>或--status<文件>   将完成状态保存为指定的文件。
  -v或--version   显示版本信息。
  -x<层级>或--debug<层级>   指定排错层级。
  -z或--notification=error   若发生错误，则以邮件来通知用户。
  --help   显示帮助。
### cu(call up)：连接另一个系统主机。
语　　法：cu [dehnotv][-a<通信端口>][-c<电话号码>][-E<脱离字符>][-I<设置文件>][-l<外围设备代号>][-s<连线速率>][-x<排错模式>][-z<系统主机>][--help][-nostop][--parity=none][<系统主机>/<电话号码>]
补充说明：本指令可连接另一台主机，并采用类似拨号终端机的接口工作，也可执行简易的文件传输作业。
参　　数：
  -a<通信端口>或-p<通信端口>或--port<通信端口>   使用指定的通信端口进行连线。
  -c<电话号码>或--phone<电话号码>   拨打该电话号码。
  -d   进入排错模式。
  -e或--parity=even   使用双同位检查。
  -E<脱离字符>或--escape<脱离字符>   设置脱离字符。
  -h或--halfduple   使用半双工模式。
  -I<配置文件>或--config<配置文件>   指定要使用的配置文件。
  -l<外围设备代号>或--line<外围设备代号>   指定某项外围设备，作为连接的设备。
  -n或--prompt   拨号时等待用户输入电话号码。
  -o或--parity=odd   使用单同位检查。
  -s<连线速率>或--speed<连线速率>或--baud<连线速率>或-<连线速率>   设置连线的速率，单位以鲍率计算。
  -t或--maper   把CR字符置换成LF+CR字符。
  -v或--version   显示版本信息。
  -x<排错模式>或--debug<排错模式>   使用排错模式。
  -z<系统主机>或--system<系统主机>   连接该系统主机。
  --help   在线帮助。
  --nostop   关闭Xon/Xoff软件流量控制。
  --parity=none   不使用同位检查。
### dnsconf(dns configurator)：设置DNS服务器组态。
语　　法：dnsconf [--deldomain<域>][--delsecondary<域>][--newdomain<域>][--set<主机><IP>][--setcname<CNAME><主机>][--setmx<域><主机>][--setns<域><主机>][--unset<主机>]
补充说明：dnsconf实际上为linuxconf的符号连接，提供图形截面的操作方式，供管理员管理DNS服务器。
参　　数：
  --deldomain<域>   删除域。
  --delsecondary<域>   删除次级域。
  --newdomain<域>   新增域。
  --set<主机><IP>   新增主机记录。
  --setcname<CNAME><主机>   设置<CNAME>。
  --setmx<域><主机>   指定域的邮件主机。
  --setns<域><主机>   指定域的DNS服务器。
  --unset<主机>   删除DNS中某台主机的记录。
### efax：收发传真。
语　　法：efax [-sw][-a<AT指令>][-c<调制解调器属性>][-d<驱动程序>][-f<字体文件>][-g<指令>][-h<传真标题字符串>][-i<AT指令>][-j<AT指令>][-k<AT指令>][-l<识别码>][-o<选项>][-q<错误次数>][-r<文件名>][-v<信息类型>][-x<UUCP锁定文件>][-t<电话号码><传真文件>]
补充说明：支持Class 1与Class 2的调制解调器来收发传真。
参　　数：
  -a<AT指令>   以指定的AT指令来接电话。
  -c<调制解调器属性>   设置本机调制解调器的属性。
  -d<驱动程序>   指定调制解调器驱动程序。
  -f<字体文件>   使用指定的字体文件来建立传真标题。
  -g<指令>   若接到的电话为数据，则执行指定的指令。
  -h<传真标题字符串>   指定字符串为每页最前端的标题。
  -i<AT指令>   在调制解调器进入传真模式前，传送AT指令到调制解调器。
  -j<AT指令>   在调制解调器进入传真模式后，传送AT指令到调制解调器。
  -k<AT指令>   在调制解调器离开传真模式前，传送AT指令到调制解调器。
  -l<识别码>   设置本机调制解调器的识别码。
  -o<选项>   使用非标准调制解调器时设置相关选项。
  -q<错误次数>   接收传真时，当每页发生错误次数超过指定的数目时，要求对方重发。
  -r<文件名>   在接收传真时，将每页分别保存成文件。
  -v<信息类型>   选择要印出的信息类型。
  -w   不要接听电话，等待OK或CONNECT的信号。
  -x<UUCP锁定文件>   使用UUCP格式的锁定文件来锁定调制解调器。
  -t<电话号码><传真文件>   以<电话号码>中的号码来拨号，并将<传真文件>传真出去。
### httpd(http daemon)：Apache HTTP服务器程序。
语　　法：httpd [-hlLStvVX][-c<httpd指令>][-C<httpd指令>][-d<服务器根目录>][-D<设定文件参数>][-f<设定文件>]
补充说明：httpd为Apache HTTP服务器程序。直接执行程序可启动服务器的服务。
参　　数：
  -c<httpd指令>   在读取配置文件前，先执行选项中的指令。
  -C<httpd指令>   在读取配置文件后，再执行选项中的指令。
  -d<服务器根目录>   指定服务器的根目录。
  -D<设定文件参数>   指定要传入配置文件的参数。
  -f<设定文件>   指定配置文件。
  -h   显示帮助。
  -l   显示服务器编译时所包含的模块。
  -L   显示httpd指令的说明。
  -S   显示配置文件中的设定。
  -t   测试配置文件的语法是否正确。
  -v   显示版本信息。
  -V   显示版本信息以及建立环境。
  -X   以单一程序的方式来启动服务器。
### ifconfig：显示或设置网络设备。
语　　法：ifconfig [网络设备][down up -allmulti -arp -promisc][add<地址>][del<地址>][<hw<网络设备类型><硬件地址>][io_addr<I/O地址>][irq<IRQ地址>][media<网络媒介类型>][mem_start<内存地址>][metric<数目>][mtu<字节>][netmask<子网掩码>][tunnel<地址>][-broadcast<地址>][-pointopoint<地址>][IP地址]
补充说明：ifconfig可设置网络设备的状态，或是显示目前的设置。
参　　数：
  add<地址>   设置网络设备IPv6的IP地址。
  del<地址>   删除网络设备IPv6的IP地址。
  down   关闭指定的网络设备。
  <hw<网络设备类型><硬件地址>   设置网络设备的类型与硬件地址。
  io_addr<I/O地址>   设置网络设备的I/O地址。
  irq<IRQ地址>   设置网络设备的IRQ。
  media<网络媒介类型>   设置网络设备的媒介类型。
  mem_start<内存地址>   设置网络设备在主内存所占用的起始地址。
  metric<数目>   指定在计算数据包的转送次数时，所要加上的数目。
  mtu<字节>   设置网络设备的MTU。
  netmask<子网掩码>   设置网络设备的子网掩码。
  tunnel<地址>   建立IPv4与IPv6之间的隧道通信地址。
  up   启动指定的网络设备。
  -broadcast<地址>   将要送往指定地址的数据包当成广播数据包来处理。
  -pointopoint<地址>   与指定地址的网络设备建立直接连线，此模式具有保密功能。
  -promisc   关闭或启动指定网络设备的promiscuous模式。
  [IP地址]   指定网络设备的IP地址。
  [网络设备]   指定网络设备的名称。
### mesg：设置终端机的写入权限。
语　　法：mesg [ny]
补充说明：将mesg设置y时，其他用户可利用write指令将信息直接显示在您的屏幕上。
参　　数：
  n   不允许气筒用户将信息直接显示在你的屏幕上。
  y   允许气筒用户将信息直接显示在你的屏幕上。
### minicom：调制解调器通信程序。
语　　法：minicom [-8lmMostz][-a<on或0ff>][-c<on或off>][-C<取文件>][-d<编号>][-p<模拟终端机>][-S<script文件>][配置文件]
补充说明：minicom是一个相当受欢迎的PPP拨号连线程序。
参　　数：
  -8   不要修改任何8位编码的字符。
  -a<on或0ff>   设置终端机属性。
  -c<on或off>   设置彩色模式。
  -C<取文件>   指定取文件，并在启动时开启取功能。
  -d<编号>   启动或直接拨号。
  -l   不会将所有的字符都转成ASCII码。
  -m   以Alt或Meta键作为指令键。
  -M   与-m参数类似。
  -o   不要初始化调制解调器。
  -p   <模拟终端机>   使用模拟终端机。
  -s   开启程序设置画面。
  -S<script文件>   在启动时，执行指定的script文件。
  -t   设置终端机的类型。
  -z   在终端机上显示状态列。
  [配置文件]   指定minicom配置文件。
### nc：设置路由器。
语　　法：nc [-hlnruz][-g<网关...>][-G<指向器数目>][-i<延迟秒数>][-o<输出文件>][-p<通信端口>][-s<来源位址>][-v...][-w<超时秒数>][主机名称][通信端口...]
补充说明：执行本指令可设置路由器的相关参数。
参　　数：
  -g<网关>   设置路由器跃程通信网关，最丢哦可设置8个。
  -G<指向器数目>   设置来源路由指向器，其数值为4的倍数。
  -h   在线帮助。
  -i<延迟秒数>   设置时间间隔，以便传送信息及扫描通信端口。
  -l   使用监听模式，管控传入的资料。
  -n   直接使用IP地址，而不通过域名服务器。
  -o<输出文件>   指定文件名称，把往来传输的数据以16进制字码倾倒成该文件保存。
  -p<通信端口>   设置本地主机使用的通信端口。
  -r   乱数指定本地与远端主机的通信端口。
  -s<来源位址>   设置本地主机送出数据包的IP地址。
  -u   使用UDP传输协议。
  -v   显示指令执行过程。
  -w<超时秒数>   设置等待连线的时间。
  -z   使用0输入/输出模式，只在扫描通信端口时使用。
### netconf：设置各项网络功能。
语　　法：netconf
补充说明：netconf是Red Hat Linux发行版专门用来调整Linux各项设置的程序。

### netconfig：设置网络环境。
语　　法：netconfig
补充说明：这是Slackware发行版内附程序，它具有互动式的问答界面，让用户轻易完成网络环境的设置。
### netstat：显示网络状态。
语　　法：netstat [-acCeFghilMnNoprstuvVwx][-A<网络类型>][--ip]
补充说明：利用netstat指令可让你得知整个Linux系统的网络情况。
参　　数：
  -a或--all   显示所有连线中的Socket。
  -A<网络类型>或--<网络类型>   列出该网络类型连线中的相关地址。
  -c或--continuous   持续列出网络状态。
  -C或--cache   显示路由器配置的快取信息。
  -e或--extend   显示网络其他相关信息。
  -F或--fib   显示FIB。
  -g或--groups   显示多重广播功能群组组员名单。
  -h或--help   在线帮助。
  -i或--interfaces   显示网络界面信息表单。
  -l或--listening   显示监控中的服务器的Socket。
  -M或--masquerade   显示伪装的网络连线。
  -n或--numeric   直接使用IP地址，而不通过域名服务器。
  -N或--netlink或--symbolic   显示网络硬件外围设备的符号连接名称。
  -o或--timers   显示计时器。
  -p或--programs   显示正在使用Socket的程序识别码和程序名称。
  -r或--route   显示Routing Table。
  -s或--statistice   显示网络工作信息统计表。
  -t或--tcp   显示TCP传输协议的连线状况。
  -u或--udp   显示UDP传输协议的连线状况。
  -v或--verbose   显示指令执行过程。
  -V或--version   显示版本信息。
  -w或--raw   显示RAW传输协议的连线状况。
  -x或--unix   此参数的效果和指定"-A unix"参数相同。
  --ip或--inet   此参数的效果和指定"-A inet"参数相同。
### ping：检测主机。
语　　法：ping [-dfnqrRv][-c<完成次数>][-i<间隔秒数>][-I<网络界面>][-l<前置载入>][-p<范本样式>][-s<数据包大小>][-t<存活数值>][主机名称或IP地址]
补充说明：执行ping指令会使用ICMP传输协议，发出要求回应的信息，若远端主机的网络功能没有问题，就会回应该信息，因而得知该主机运作正常。
参　　数：
  -d   使用Socket的SO_DEBUG功能。
  -c<完成次数>   设置完成要求回应的次数。
  -f   极限检测。
  -i<间隔秒数>   指定收发信息的间隔时间。
  -I<网络界面>   使用指定的网络界面送出数据包。
  -l<前置载入>   设置在送出要求信息之前，先行发出的数据包。
  -n   只输出数值。
  -p<范本样式>   设置填满数据包的范本样式。
  -q   不显示指令执行过程，开头和结尾的相关信息除外。
  -r   忽略普通的Routing Table，直接将数据包送到远端主机上。
  -R   记录路由过程。
  -s<数据包大小>   设置数据包的大小。
  -t<存活数值>   设置存活数值TTL的大小。
  -v   详细显示指令的执行过程。
### pppstats(point to point protocol status)：显示PPP连线状态。
语　　法：pppstats [-adrv][-c<执行次数>][-w<间隔秒数>][网络界面]
补充说明：利用pppstats指令可让你得知PPP连接网络的相关信息。
参　　数：
  -a   显示绝对统计值。
  -c<执行次数>   设置回报状况的次数。
  -d   显示相对统计值。
  -r   显示数据包压缩比率的统计值。
  -v   显示VJTCP文件头的压缩效率统计值。
  -w<间隔秒数>   设置显示统计信息的间隔时间。
### samba：Samba服务器控制。
语　　法：samba [start][stop][status][restart]
补充说明：samba为script文件，可启动，停止Samba服务器或回报目前的状态。
参　　数：
  start   启动Samba服务器的服务。
  stop    停止Samba服务器的服务。
  status   显示Samba服务器目前的状态。
  restart   重新启动Samba服务器。
### setserial：设置或显示串口的相关信息。
语　　法：setserial [-abgGqvVz][设备][串口参数]
补充说明：setserial可用来设置串口或显示目前的设置。
参　　数：
  -a   显示详细信息。
  -b   显示摘要信息。
  -g   显示串口的相关信息。
  -G   以指令列表的格式来显示信息。
  -q   执行时显示较少的信息。
  -v   执行时显示较多的信息。
  -V   显示版本信息。
  -z   设置前，先将所有的标记归零。
### shapecfg(shaper configuration)：管制网络设备的流量。
语　　法：shapecfg attach [流量管制器][网络设备] 或  shapecfg speed [流量管制器][带宽]
补充说明：自Linux-2.15开始，便支持流量管制的功能。
参　　数：
  attach   将流量管制器与实际的网络设备结合。
  speed   设置流量管制器的对外传输带宽。
### smbd(samba daemon)：Samba服务器程序。
语　　法：smbd [-aDhoP][-d<排错层级>][-i<范围>][-l<记录文件>][-O<连接槽选项>][-p<连接端口编号>][-s<配置文件>]
补充说明：smbd为Samba服务器程序，可分享文件与打印机等网络资源供Windows相关的用户端程序存取。
参　　数：
  -a   所有的连线记录都会加到记录文件中。
  -d<排错层级>   指定记录文件所记载事件的详细程度。
  -D   使用此参数时，smbd会以服务程序的方式在后台执行。
  -h   显示帮助。
  -i<范围>   指定NetBIOS名称的范围。
  -l<记录文件>   指定记录文件的名称。
  -o   每次启动时，会覆盖原有的记录文件。
  -O<连接槽选项>   设置连接槽选项。
  -p<连接端口编号>   设置连接端口编号。
  -P   仅用来测试smbd程序的正确性。
  -s<配置文件>   指定smbd的设置文件。
### statserial(status ofserial port)：显示串口状态。
语　　法：statserial [-dnx][串口设备名称]
补充说明：statserial可显示各个接脚的状态，常用来判断串口是否正常。
参　　数：
  -d   以10进制数字来表示串口的状态。
  -n   仅显示一次串口的状态后即结束程序。
  -x   与-n参数类似，但是以16进制来表示。
### talk：与其他用户交谈。
语　　法：talk [用户名称][终端机编号]
补充说明：通过talk指令，你可以和另一个用户线上交谈。
### Tcpdump：倾倒网络传输数据。
语　　法：tcpdump [-adeflnNOpqStvx][-c<数据包数目>][-dd][-ddd][-F<表达文件>][-i<网络界面>][-r<数据包文件>][-s<数据包大小>][-tt][-T<数据包类型>][-vv][-w<数据包文件>][输出数据栏位]
补充说明：执行tcpdump指令可列出经过指定网络界面的数据包文件头，在Linux操作系统中，你必须是系统管理员。
参　　数：
  -a   尝试将网络和广播地址转换成名称。
  -c<数据包数目>   收到指定的数据包数目后，就停止进行倾倒操作。
  -d   把编译过的数据包编码转换成可阅读的格式，并倾倒到标准输出。
  -dd   把编译过的数据包编码转换成C语言的格式，并倾倒到标准输出。
  -ddd   把编译过的数据包编码转换成十进制数字的格式，并倾倒到标准输出。
  -e   在每列倾倒资料上显示连接层级的文件头。
  -f   用数字显示网际网络地址。
  -F<表达文件>   指定内含表达方式的文件。
  -i<网络界面>   使用指定的网络截面送出数据包。
  -l   使用标准输出列的缓冲区。
  -n   不把主机的网络地址转换成名字。
  -N   不列出域名。
  -O   不将数据包编码最佳化。
  -p   不让网络界面进入混杂模式。
  -q   快速输出，仅列出少数的传输协议信息。
  -r<数据包文件>   从指定的文件读取数据包数据。
  -s<数据包大小>   设置每个数据包的大小。
  -S   用绝对而非相对数值列出TCP关联数。
  -t   在每列倾倒资料上不显示时间戳记。
  -tt   在每列倾倒资料上显示未经格式化的时间戳记。
  -T<数据包类型>   强制将表达方式所指定的数据包转译成设置的数据包类型。
  -v   详细显示指令执行过程。
  -vv   更详细显示指令执行过程。
  -x   用十六进制字码列出数据包资料。
  -w<数据包文件>   把数据包数据写入指定的文件。
### testparm(test parameter)：测试Samba的设置是否正确无误。
语　　法：testparm [-s][配置文件][<主机名称><IP地址>]
补充说明：执行testparm指令可以简单测试Samba的配置文件，假如测试结果无误，Samba常驻服务就能正确载入该设置值，但并不保证其后的操作如预期般一切正常。
参　　数：
  -s   不显示提示符号等待用户按下Enter键，就直接列出Samba服务定义信息。
### traceroute：显示数据包到主机间的路径。
语　　法：traceroute [-dFlnrvx][-f<存活数值>][-g<网关>...][-i<网络界面>][-m<存活数值>][-p<通信端口>][-s<来源地址>][-t<服务类型>][-w<超时秒数>][主机名称或IP地址][数据包大小]
补充说明：traceroute指令让你追踪网络数据包的路由途径，预设数据包大小是40Bytes，用户可另行设置。
参　　数：
  -d   使用Socket层级的排错功能。
  -f<存活数值>   设置第一个检测数据包的存活数值TTL的大小。
  -F   设置勿离断位。
  -g<网关>   设置来源路由网关，最多可设置8个。
  -i<网络界面>   使用指定的网络界面送出数据包。
  -I   使用ICMP回应取代UDP资料信息。
  -m<存活数值>   设置检测数据包的最大存活数值TTL的大小。
  -n   直接使用IP地址而非主机名称。
  -p<通信端口>   设置UDP传输协议的通信端口。
  -r   忽略普通的Routing Table，直接将数据包送到远端主机上。
  -s<来源地址>   设置本地主机送出数据包的IP地址。
  -t<服务类型>   设置检测数据包的TOS数值。
  -v   详细显示指令的执行过程。
  -w<超时秒数>   设置等待远端主机回报的时间。
  -x   开启或关闭数据包的正确性检验。
### tty(teletypewriter)：显示终端机连接标准输入设备的文件名称。
语　　法：tty [-s][--help][--version]
补充说明：在Linux操作系统中，所有外围设备都有其名称与代号，这些名称代号以特殊文件的类型存放于/dev目录下。你可以执行tty指令查询目前使用的终端机的文件名称。
参　　数：
  -s或--silent或--quiet   不显示任何信息，只回传状态代码。
  --help   在线帮助。
  --version   显示版本信息。
### uuname：显示全部的UUCP远端主机。
语　　法：uuname [-alv][-I<配置文件>][--help]
补充说明：uuname可显示UUCP远端主机。
参　　数：
  -a或--aliases   显示别名。
  -I<配置文件>或--config<配置文件>   指定程序的配置文件。
  -l或--local   显示本机名称。
  -v或--version   显示版本信息。
  --help   显示帮助。
### wall(write all)：传送信息。
语　　法：wall [公告信息]
补充说明：通过wall指令可将信息发送给每位同意接收公众信息的终端机用户，若不给予其信息内容，则wall指令会从标准输入设备读取数据，然后再把所得到的数据传送给所有终端机用户。
### write：传送信息。
语　　法：write [用户名称][终端机编号]
补充说明：通过write指令可传递信息给另一位登入系统的用户，当输入完毕后，键入EOF表示信息结束，write指令就会将信息传给对方。如果接收信息的用户不只登入本地主机一次，你可以指定接收信息的终端机编号。
### ytalk：与其他用户交谈。
语　　法：ytalk [-isxY][-h<主机名称IP地址>][用户名称...]
补充说明：通过ytalk指令，你可以和其他用户线上交谈，如果想和其他主机的用户交谈，在用户名称后加上其主机名称或IP地址即可。
参　　数：
  -h<主机名称IP地址>   指定交谈对象所在的远端主机。
  -i   用提醒声响代替显示信息。
  -s   在指令提示符号先开启ytalk交谈窗。
  -x   关闭图形界面。
  -Y   所有必须回应yes或no的问题，都必须用大写英文字母"Y"或"N"回答。
### arpwatch(ARP watcher)：监听网络上ARP的记录。
语　　法：arpwatch [-d][-f<记录文件>][-i<接口>][-r<记录文件>]
补充说明：ARP(Address Resolution Protocol)是用来解析IP与网络装置硬件地址的协议。arpwatch可监听区域网络中的ARP数据包并记录，同时将监听到的变化通过E-mail来报告。
参　　数：
  -d   启动排错模式。
  -f<记录文件>   设置存储ARP记录的文件，预设为/var/arpwatch/arp.dat。
  -i<接口>   指定监听ARP的接口，预设的接口为eth0。
  -r<记录文件>   从指定的文件中读取ARP记录，而不是从网络上监听。
### apachectl(Apache control interface)：可用来控制Apache HTTP服务器的程序。
语　　法：apachectl [configtest][fullstatus][graceful][help][restart][start][status][stop]
补充说明：apachectl是slackware内附Apache HTTP服务器的script文件，可供管理员控制服务器，但在其他Linux的Apache HTTP服务器不一定有这个文件。
参　　数：
  configtest   检查设置文件中的语法是否正确。
  fullstatus   显示服务器完整的状态信息。
  graceful   重新启动Apache服务器，但不会中断原有的连接。
  help   显示帮助信息。
  restart   重新启动Apache服务器。
  start   启动Apache服务器。
  status   显示服务器摘要的状态信息。
  stop   停止Apache服务器。
### smbclient(samba client)：可存取SMB/CIFS服务器的用户端程序。
语　　法：smbclient [网络资源][密码][-EhLN][-B<IP地址>][-d<排错层级>][-i<范围>][-I<IP地址>][-l<记录文件>][-M<NetBIOS名称>][-n<NetBIOS名称>][-O<连接槽选项>][-p<TCP连接端口>][-R<名称解析顺序>][-s<目录>][-t<服务器字码>][-T<tar选项>][-U<用户名称>][-W<工作群组>]
补充说明：SMB与CIFS为服务器通信协议，常用于Windows95/98/NT等系统。smbclient可让Linux系统存取Windows系统所分享的资源。
参　　数：
  [网络资源]  [网络资源]的格式为//服务器名称/资源分享名称。
  [密码]   输入存取网络资源所需的密码。
  -B<IP地址>   传送广播数据包时所用的IP地址。
  -d<排错层级>   指定记录文件所记载事件的详细程度。
  -E   将信息送到标准错误输出设备。
  -h   显示帮助。
  -i<范围>   设置NetBIOS名称范围。
  -I<IP地址>   指定服务器的IP地址。
  -l<记录文件>   指定记录文件的名称。
  -L   显示服务器端所分享出来的所有资源。
  -M<NetBIOS名称>   可利用WinPopup协议，将信息送给选项中所指定的主机。
  -n<NetBIOS名称>   指定用户端所要使用的NetBIOS名称。
  -N   不用询问密码。
  -O<连接槽选项>   设置用户端TCP连接槽的选项。
  -p<TCP连接端口>   指定服务器端TCP连接端口编号。
  -R<名称解析顺序>   设置NetBIOS名称解析的顺序。
  -s<目录>   指定smb.conf所在的目录。
  -t<服务器字码>   设置用何种字符码来解析服务器端的文件名称。
  -T<tar选项>   备份服务器端分享的全部文件，并打包成tar格式的文件。
  -U<用户名称>   指定用户名称。
  -W<工作群组>   指定工作群组名称。
### pppsetup：设置PPP连线。
语　　法：pppsetup
补充说明：这是Slackware发行版内附程序，它具有互动式的问答界面，让用户轻易完成PPP的连线设置。
 
## linux电子邮件与新闻组命令
### archive：新闻组文件备存程序。
语　　法：archive [-fmr][-a<备存目录>][-i<索引文件>][源文件]
补充说明：archive会读取新闻组的文件，然后加以备存。archive备存时，会保存[源文件]中文件的目录阶层。
参　　数：
  -a<备存目录>   指定备存的目录。
  -f   将阶层目录名称的阶层关系取消，以句点来取代斜线。
  -i<索引文件>   制作文件索引。
  -m   将文件复制到目的地，建立文件连接后，以连接取代原始的文件。
  -r   archive预设的标准错误输出为/var/errlog。-r参数可取消此功能。
  [源文件]   在源文件中指定要备份的文件。
### ctlinnd(control the internet news daemon)：设置INN新闻组服务器。
语　　法：ctlinnd [-hs][-t<超时秒数>][控制命令<设置数值>...]
补充说明：通过ctlinnd指令，可以对新闻组服务器这项常驻服务下达控制命令，直接设置相关数值。
参　　数：
  -h   在线帮助。
  -s   不显示指令执行过程。
  -t<超时秒数>   设置等待服务器回报的时间，单位以秒计算，超过此一时间则为超时。
### elm：E-mail管理程序。
语　　法：elm [-s<邮件主题>][收信人别名或地址]<[要输入的文件] 或
                 elm [-hmvz][-f<目录>][-i<文件名>]
补充说明：elm为全屏文字模式的电子邮件管理程序。与mail相比，elm提供较易操作的界面，但仍不及pine来得友善。
参　　数：
  -s<邮件主题>   指定新邮件的邮件主题。
  -f<目录>   开启程序时，读取指定的目录。
  -h   显示帮助。
  -i<文件名>   将文件内容插入送出的邮件中。
  -m   进入elm后，不显示指令说明。
  -v   显示elm的版本信息。
  -z   若收件信箱没有邮件，则不启动elm程序。
### getlist：下载新闻组清单。
语　　法：getlist [-h<新闻组服务器>][-p<通信端口>][群组类型<范本样式><群组名称>]
补充说明：getlist指令会从新闻服务器下载新闻组清单，并把清单输出到标准输出设备。
参　　数：
  -h<新闻组服务器>   指定欲连接的新闻组服务器。
  -p<通信端口>   设置连接新闻组服务器的通讯端口。
### inncheck(inn check)：检查inn相关配置文件的语法是否正确。
语　　法：inncheck [-afv][-pedantic][-perm][-noperm][文件]
补充说明：inncheck会检查inn配置文件的语法是否正确，但不会实际修改这些配置文件。
参　　数：
  -a   不管是否出现[文件]参数，仍然检查全部的配置文件。
  -f   若发现配置文件有误，则列出可修正错误的指令。
  -v   执行时显示详细的信息。
  -pedantic   以较为严格的方式来检查。
  -perm   检查文件权限是否有误。
  --noperm   不检查文件的权限设置。
  [文件]   若不指定此参数，则inncheck会检查所有预设的文件。
### mail：E-mail管理程序。
语　　法：mail [-iInNv][-b<地址>][-c<地址>][-f<邮件文件>][-s<邮件主题>][-u<用户帐号>][收信人地址]
补充说明：mail是一个文字模式的邮件管理程序，操作的界面不像elm或pine那么容易使用，但功能尚称完整。
参　　数：
  -b<地址>   指定密件副本的收信人地址。
  -c<地址>   指定副本的收信人地址。
  -f<邮件文件>   读取指定邮件文件中的邮件。
  -i   不显示终端发出的信息。
  -I   使用互动模式。
  -n   程序使用时，不使用mail.rc文件中的设置。
  -N   阅读邮件时，不显示邮件的标题。
  -s<邮件主题>   指定邮件的主题。
  -u<用户帐号>   读取指定用户的邮件。
  -v   执行时，显示详细的信息。
### mailconf：Sendmail设置程序。
语　　法：mailconf [--addvdom<虚拟网络>][--delvdom<虚拟网络>][--help]
补充说明：mailconf实际为linuxconf的符号连接，用来设置Sendmail相关的选项。执行mailconf时，不管目前的环境为指令列模式或是X Window，都会出现图形界面的设置画面，让您轻松地完成设置。所有指令列的参数都可以在图形界面中设置。建议不用加上任何参数，直接进入mailconf的图形界面来设置。
### mailq(mail queue)：显示待寄邮件的清单。
语　　法：mailq [-q]
补充说明：mailq可列出待寄邮件的清单，包括邮件ID，邮件大小，邮件保存时间，寄信人，收信人，以及邮件无法寄出的原因，提供管理员参考的信息。mailq实际上是执行sendmail -bp指令。
参　　数：
  -v   显示较详细的信息。
### messages：查看信箱中的邮件数。
语　　法：messages [邮件目录或文件]
补充说明：messages实际上是一个script文件，可用来检查信箱中的邮件数。若不指定[邮件目录或文件]参数，则显示目前收件信箱中的邮件数。
### metamail：负责处理非文字E-mail的程序。
语　　法：metafile [-bBdehpPqrRTwxyz][-c<邮件内容类型>][-f<寄信人地址>][-m<邮件管理程序名称>][-s<邮件主题>][文件名]
补充说明：当E-mail管理程序遇到非ISO-8859-1标准文字的邮件内容，而不知道要如何处理的时候，即调用metamail。
参　　数：
  -b   邮件只有内文，没有表头。
  -B   在背景显示信息。
  -c<邮件内容类型>   不要读取邮件表头的"Content-type"栏位，而去读取<邮件内容类型>中的设置，以决定要调用哪个程序来处理邮件信息。
  -d   metamail在调用其他程序显示信息前，不需向用户确认。
  -e   删除内文前端的空白行。
  -f<寄信人地址>   不要读取邮件表头中的寄信人地址，而去读取选项中的设置。
  -h   使用metamail的邮件管理程序名称。
  -m<邮件管理程序名称>   记录调用metamail的邮件管理程序名称。
  -p   将显示内容送给more指令处理之后，再送到输出设备。
  -P   与-p参数功能类似，但会在信息结束显示"Press RETURN to go on"字符串。
  -q   metamail在调用其他程序来显示信息前，不用在屏幕上显示相关的信息。
  -r   metamail预设无法让root用户执行。
  -R   显示信息前先清除屏幕。
  -s<邮件主题>   不要读取邮件表头中的邮件主题，而去读取选项中的设置。
  -T   若要在一个新的终端机模拟视窗中重新启动metamail，则metamail必须以递归模式来调用程序本身，这时便需要使用-T参数。
  -w   不要读取mailcap中的设置，而直接处理邮件，并将数据存入文件中。
  -x   不要以终端机模式来执行，即不使用互动模式。
  -y   抽出邮件内文中MIME格式的部分。
  -z   程序结束时将输入的文件删除。
### mutt：E-mail管理程序。
语　　法：mutt [-hnpRvxz][-a<文件>][-b<地址>][-c<地址>][-f<邮件文件>][-F<配置文件>][-H<邮件草稿>][-i<文件>][-m<类型>][-s<主题>][邮件地址]
补充说明：mutt是一个文字模式的邮件管理程序，提供了全屏幕的操作界面。
参　　数：
  -a<文件>   在邮件中加上附加文件。
  -b<地址>   指定密件副本的收信人地址。
  -c<地址>   指定副本的收信人地址。
  -f<邮件文件>   指定要载入的邮件文件。
  -F<配置文件>   指定mutt程序的设置文件，而不读取预设的.muttrc文件。
  -h   显示帮助。
  -H<邮件草稿>   将指定的邮件草稿送出。
  -i<文件>   将指定文件插入邮件内文中。
  -m<类型>   指定预设的邮件信箱类型。
  -n   不要去读取程序培植文件(/etc/Muttrc)。
  -p   在mutt中编辑完邮件后，而不想将邮件立即送出，可将该邮件暂缓寄出。
  -R   以只读的方式开启邮件文件。
  -s<主题>   指定邮件的主题。
  -v   显示mutt的版本信息以及当初编译此文件时所给予的参数。
  -x   模拟mailx的编辑方式。
  -z   与-f参数一并使用时，若邮件文件中没有邮件即不启动mutt。
### nntpget：下载新闻组文章。
语　　法：nntpget [-ov][-d<组类型清单>][-f<文件或目录>][-n<组名称清单>][-t<时间日期>][-u<文件或目录>][新闻组服务器]
补充说明：执行nntpget指令将从指定的新闻组服务器下载文章。nntpget指令会从标准输入设备读取信息识别码，而下载的新闻文章则呈现到标准输出设备。
参　　数：
  -d<组类型清单>   指定新闻组类型为distribution的清单。
  -f<文件或目录>   下载其更改时间较指定文件或目录的更改时间更接近现在的新闻文章。
  -n<组名称清单>   指定新闻组名称的清单。
  -o   下载本地主机所没有的新闻文章。
  -t<时间日期>   指定时间与日期，较该日期时间更新的新闻文章都会被下载。
  -u<文件或目录>   此参数的效果和指定"-f"参数类似，但它在成功下载文章之后会去更改指定文件或木的更改时间，便于下次再作为参考。
  -v   当指令执行时，把每篇文章的信息识别码都送到标准输出。本参数需配合参数"-o"使用。
### pine：收发电子邮件，浏览新闻组。
语　　法：pine [-ahikorz][-attach<附件>][-attach_and_delete<附件>][-attachlist<附件清单>][-c<邮件编号>][-conf][-create_lu<地址薄><排序法>][-f<收件箱>][-F<文件>][-I<暂存快捷键>][-n<邮件编号>][-nr][-p<环境文件>][-pinerc<输出文件>][-P<配置文件>][-sort<排序法></reverse>][-url<URL>][-<功能选项>=<设置值>][用户名称或电子邮件地址...]
补充说明：pine是个以显示导向为住的邮件处理程序，包括一般电子邮件和新闻组的文章。pine也支持MIME格式，可让用户将MIME的对象保存成文件，甚至指定外挂程序直接开启对象。
参　　数：
  -a   在UWIN中使用匿名模式。
  -attach<附件>   把指定的文件附加在信件里，一并夹带寄出。
  -attach_and_delete   此参数的效果和指定"-attach"参数类似，但在信件寄出后<附件>会将指定的文件删除。
  -attachlist<附件清单>   设置一份文件清单，把该清单内的文件全部附加到信件中，夹带寄出。
  -c<邮件编号>   直接开启指定编号的邮件。
  -conf   产生pine.conf的范例文件，并输出到标准输出设备。
  -create_lu<地址薄><排序法>   将地址薄以指定的排序法排序，产生.addressbook.lu索引文件。
  -f<收件箱>   开启指定的接收信件箱，而非预设的INBOX收件箱。
  -F<文件>   直接浏览指定文件的内容。
  -h   在线帮助。
  -i   直接进入收件箱。
  -I<暂存快捷键>   设置进入pine就要执行的命令，该命令以快捷键的按键指定。
  -k   支持键盘上的F1,F2...等功能，而非使用一般快捷键。
  -n<邮件编号>   设置现在的邮件编号。
  -nr   使用UWIN的特殊模式。
  -o   设置收件箱为只读。
  -p<环境文件>   另外指定pine的环境文件。
  -pinerc<输出文件>   读取.pinerc环境文件的内容并输出成文件保存起来，您必须指定输出文件的名称。
  -r   展示模式。
  -sort<排序法>   把收件箱内的信件依指定的方式排序。
  -url<URL>   直接开启给予的URL。
  -z   让pine可被Ctrl+z中断，暂存在后台作业里。
  -<功能选项>=<设置值>   暂时指定各项功能的设置值。
### slrn：新闻组阅读程序。
语　　法：slrn [-Cdnp][-create][-help][-f<新闻组文件>][-h<服务器>][--version]
补充说明：slrn是一个很容易使用的新闻组阅读程序，提供全屏幕的图形操作界面，让用户能轻松地阅读与管理新闻组。
参　　数：
  -C   使用彩色模式。
  -d   下载组清单时，也一并下载每个组的说明。
  -f<新闻组文件>   指定要存储组内容的文件，预设的文件名为.newsrc。
  -h<服务器>   指定新闻组服务器。
  -n   不要检查是否有新的新闻组。
  -p   设置新闻组服务器的连接端口编号，预设为119。
  -create   第一次连上某台新闻组服务器时，必须使用此参数，以便下载服务器中的组清单，并保存在指定的newsrc文件中。
  -help   显示帮助。
  --version   显示版本信息。
 
## linux其他命令
### reconfig：转换配置文件。
语　　法：reconfig [Xconfig] XF86Config
补充说明：reconfig指令能将XFree86 3.1及以前的版本的配置文件，转成新的格式。现今的配置文件XF86Config，存放在/etc或/etc/X11目录下，它的格式包含了更多的信息，因此在转换之后，您还须自行编辑设置文件，输入不足的部分。
### startx(start X Window)：启动X Window。
语　　法：startx [程序][--服务器设置]
补充说明：startx为启动X Window的script文件，实际上启动X Window的程序为xinit。
参　　数：
  [程序]   此处所指定的[程序]为X Window所要执行程序的配置文件。
  [--服务器设置]   startx预设会去读取.xerverrc中的配置，决定启动X Window的方式。
### xconfigurator：设置XFree86。
语　　法：Xconfigurator [--card<显卡型号>][--expert][--help][--hsync"<水平扫描频率>"][-kickstart][--monitor<显示器型号>][--server<显示服务器>][--vsync"<垂直扫描频率>"]
补充说明：这是Red Hat公司遵循GPL规则所开发的程序，它具有互动式操作界面，您可以轻易地利用方向键和空格键等，设置显卡和显示器，并选择欲使用的画面分辨率与色彩深度，最后保存配置文件XF86Config。
参　　数：
  --card<显卡型号>   设置显卡的型号，可用"--help"参数查询型号列表。
  --expert   启用专家模式，增加用户得以自行选择的项目而非由程序完全自动检测。
  --help   在线帮助。
  --hsync"<水平扫描频率>"   设置水平扫描频率。
  --kickstart   启用全自动模式。
  --monitor<显示器型号>   设置显示器型号，可用"--help"参数查询型号列表。
  --server<显示服务器>   设置显示服务器。
  --vsync"<垂直扫描频率>"   设置垂直扫描频率。
### XF86Setup：设置XFee86。
语　　法：XF86Setup [-display<主机名称或IP地址>:<显示器编号>][-nodialog][-sync]
补充说明：这是Linux系统用来设置XFee86的程序，它会进入图形模式，通过互动操作界面，让用户轻松完成XFee86环境的设置。
参　　数：
  -display<主机名称或IP地址>:<显示器编号>   指定显示XF86Setup设置画面的主机名称以及欲在该主机上的哪个显示器显示画面。
  -nodialog   执行XF86Setup时，不显示交谈窗，直接以文字模式提出询问。
  --sync   和X Server进行数据同步传输。
### xlsatoms：列出X Server定义的成分。
语　　法：xlsatoms [-display<显示器编号>][-format<输出格式>][-name<成分名称>][-range<列表范围>]
补充说明：执行xlsatoms指令会列出X Server内部所有定义的成分，每个成分都有其编号，您可利用参数设置列表范围，或直接指定欲查询的成分名称。
参　　数：
  -display<显示器编号>   指定X Server连接的显示器编号，该编号由"0"开始计算， 依序递增。
  -format<输出格式>   设置成分清单的列表格式，您可使用控制字符改变显示样式。
  -name<成分名称>   列出指定的成分。
  -range<列表范围>   设置成分清单的列表范围。 
### xlsclients：列出显示器中的客户端应用程序。
语　　法：xlsclients [-al][-display<显示器编号>][-m<最大指令长度>]
补充说明：执行xlsclients指令可列出某个显示器中，正在执行的客户端应用程序信息。
参　　数：
  -a   列出所有显示器的客户端应用程序信息。
  -display<显示器编号>   指定X Server连接的显示器编号，该编号由"0"开始计算，依序递增。
  -l   使用详细格式列表。
  -m<最大指令长度>   设置显示指令信息的最大长度，单位以字符计算。
### xlsfonts：列出X Server使用的字体。
语　　法：xlsfonts [-1Clmou][-display<主机名称或IP地址>:<显示器编号>][-fn<范本样式>][-ll][-lll][-n<显示栏位数>][-w<每列字符数>]
补充说明：执行xlsfonts指令会显示目前X Server可使用的字体，也能使用范本样式仅列出的符合条件的字体。
参　　数：
  -1   仅用单一栏位列出字体清单。
  -C   使用多个栏位列出字体清单。
  -display<主机名称或IP地址>:<显示器编号>   指定显示XF86Setup设置画面的主机名称以及欲放在该主机机上的哪个显示器画面。
  -fn<范本样式>   只列出符合范本样式指定的条件的字体。
  -l   除字体名称外，同时列出字体的属性。
  -ll   此参数的效果和指定"l"参数类似，但显示更详细的信息。
  -lll   此参数的效果和指定"ll"参数类似，但显示更详细的信息。
  -m   配合参数"-l"使用时，一并列出字体大小的上下限。
  -n<显示栏位数>   设置每列显示的栏位数。
  -o   以OpenFont的形式列出字体清单。
  -u   列出字体清单时不依照其名称排序。
  -w<每列字符数>   设置每列的最大字符数。
yes
1.命令说明
输出回应的字符串。yes回应“y”字符，后者指定字符。 
2.路径
/usr/bin/yes
3.语法
yes [--help] [--version] string
4.参数
--help 显示帮助字符。 --version 显示版本信息 
5.相关命令
无
6.举例说明
[root@bixuan/]#yes "OK" //回应“OK”字符串。
命令执行结果如下：
OK
OK
OK 
...
Ctrl+c //按Ctrl+c键中断 


