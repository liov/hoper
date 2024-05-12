#include <stdio.h>  
#include <string.h>  
#include <stdlib.h>  
#include <sys/inotify.h>  
#include <unistd.h>  
  
#define EVENT_NUM 12  
char *event_str[EVENT_NUM] =  
{  
"IN_ACCESS",  
"IN_MODIFY",        //文件修改  
"IN_ATTRIB",  
"IN_CLOSE_WRITE",  
"IN_CLOSE_NOWRITE",  
"IN_OPEN",  
"IN_MOVED_FROM",    //文件移动from  
"IN_MOVED_TO",      //文件移动to  
"IN_CREATE",        //文件创建  
"IN_DELETE",        //文件删除  
"IN_DELETE_SELF",  
"IN_MOVE_SELF"  
};  

// TODO:转成abi调用
int main(int argc, char *argv[])  
{  
    int fd;  
    int wd;  
    int len;  
    int nread;  
    char buf[BUFSIZ];  
    struct inotify_event *event;  
    int i;  
  
    // 判断输入参数  
    if (argc < 2) {  
        fprintf(stderr, "%s path\n", argv[0]);  
        return -1;  
    }  
  
    // 初始化  
    fd = inotify_init();  
    if (fd < 0) {  
        fprintf(stderr, "inotify_init failed\n");  
        return -1;  
    }  
  
    /* 增加监听事件 
     * 监听所有事件：IN_ALL_EVENTS 
     * 监听文件是否被创建,删除,移动：IN_CREATE|IN_DELETE|IN_MOVED_FROM|IN_MOVED_TO 
     */  
    wd = inotify_add_watch(fd, argv[1], IN_CREATE|IN_DELETE|IN_MOVED_FROM|IN_MOVED_TO);  
    if(wd < 0) {  
        fprintf(stderr, "inotify_add_watch %s failed\n", argv[1]);  
        return -1;  
    }  
  
    buf[sizeof(buf) - 1] = 0;  
    while( (len = read(fd, buf, sizeof(buf) - 1)) > 0 ) {  
        nread = 0;  
        while(len> 0) {  
            event = (struct inotify_event *)&buf[nread];  
            for(i=0; i<EVENT_NUM; i++) {  
                if((event->mask >> i) & 1) {  
                    if(event->len > 0){
                        fprintf(stdout, "%s --- %s%d\n", event->name, event_str[i]);
						if(i==8){
								char command[] = "/grame/sh ";  
								strcat(command,event->name);  
								system(command);  
							}
					}						
                    else  
                        fprintf(stdout, "%s --- %s\n", " ", event_str[i]);  
				}				
            }  
            nread = nread + sizeof(struct inotify_event) + event->len;  
            len = len - sizeof(struct inotify_event) - event->len;  
        }  
    }  
  
    return 0;  
}  