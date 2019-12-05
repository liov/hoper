针对IO，总是涉及到阻塞、非阻塞、异步、同步以及select/poll和epoll的一些描述，那么这些东西到底是什么，有什么差异？

一般来讲一个IO分为两个阶段：
等待数据到达
把数据从内核空间拷贝到用户空间

现在假设一个进程/线程A，试图进行一次IO操作。
A发出IO请求，两种情况：
  1）立即返回
  2）由于数据未准备好，需要等待，让出CPU给别的线程，自己sleep
  第一种情况就是非阻塞，A为了知道数据是否准备好，需要不停的询问，而在轮询的空歇期，理论上是可以干点别的活，例如喝喝茶、泡个妞。
  第二种情况就是阻塞，A除了等待就不能做任何事情。数据终于准备好了，A现在要把数据取回去，有几种做法：  1）A自己把数据从内核空间拷贝到用户空间。
  2）A创建一个新线程（或者直接使用内核线程），这个新线程把数据从内核空间拷贝到用户空间。
  第一种情况，所有的事情都是同一个线程做，叫做同步，有同步阻塞（BIO）、同步非阻塞（NIO）
  第二种情况，叫做异步，只有异步非阻塞（AIO）

同步阻塞：

  同一个线程在IO时一直阻塞，直到读取数据成功，把数据从核心空间拷贝到用户空间

同步非阻塞：

  同一个线程发起IO后，立即获得返回，后面定期轮询数据读取情况，发现数据读取成功，把数据从核心空间拷贝到用户空间

 异步非阻塞：

  一个线程发起IO后，立即返回，由另外的线程发现数据读取成功，把数据从核心空间拷贝到用户空间。

下面说一下多路复用：select/poll、epoll

select是几乎所有unix、linux都支持的一种多路IO方式，通过select函数发出IO请求后，线程阻塞，一直到数据准备完毕，然后才能把数据从核心空间拷贝到用户空间，所以select是同步阻塞方式。
  int select(int n, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout);
参数n表示监控的所有fd中最大值＋1
readfds、writefds和exceptfds分别表示可读、可写、异常的文件句柄，这个文件句柄中每一个bit表示一个文件fd，所以能够表示的最大文件数和fd_set的长度有关，
  假设fd_set的长度为1字节（即8bit），则可以表示8个可读文件、8个可写文件、8个异常文件句柄。下面以读文件为例：
  使用select的时候，先初始化FD_ZERO(fd_set *set)，把8bit全部置为0，readfds=00000000
  使用FD_SET(int fd, fd_set *set)来把文件fd设置到fd_set中，例如3个文件fd=2，fd=3，fd=5，则readfds=00010110
  然后使用select(6, readfds, 0, 0, 0)阻塞等待，若此时fd=2文件可读，则此时readfds=00000010（fd=5和fd=3对应的bit被清0）
  使用FD_ISSET(int fd, fd_set *set)函数来判断fd对应的bit是否为1，如果为1则可读。


poll对select的使用方法进行了一些改进，突破了最大文件数的限制，同时使用更加方便一些。
int poll(struct pollfd *ufds, unsigned int nfds, int timeout);
struct pollfd {
    int fd;           /* 对应的文件描述符 */
    short events;     /* 要监听的事件，例如POLLIN|POLLPRI */
    short revents;    /* 返回的事件，用于在poll返回时携带该fd上发生的事情，在poll调用时，该字段会自动被清空 */
};
通过poll函数发出IO请求后，线程阻塞，直到数据准备完毕，poll函数在pollfd中通过revents字段返回事件，然后线程把数据从核心空间拷贝到用户空间，
所以poll同样是同步阻塞方式，性能同select相比没有改进。


epoll是linux为了解决select/poll的性能问题而新搞出来的机制，基本的思路是：由专门的内核线程来不停地扫描fd列表，有结果后，把结果放到fd相关的链表中，
用户线程只需要定期从该fd对应的链表中读取事件就可以了。同时，为了节省把数据从核心空间拷贝到用户空间的消耗，采用了mmap的方式，允许程序在用户空间直接访问数据所在的内核空间，不需要把数据copy一份。
epoll一共有3个函数：
1.创建epoll文件描述符
int epoll_create(int size)；

2.把需要监听的文件fd和事件加入到epoll文件描述符，也可以对已有的fd进行修改和删除
文件fd保存在一个红黑树中，该fd的事件保存在一个链表中（每个fd一个事件链表），事件由内核线程负责填充，用户线程读取
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event)；
            typedef union epoll_data {
                void *ptr;
                int fd;
                __uint32_t u32;
                __uint64_t u64;
            } epoll_data_t;


            struct epoll_event {
                __uint32_t events;      /* Epoll events */
                epoll_data_t data;      /* User data variable */
            };

3.用户线程定期轮询epoll文件描述符上的事件，事件发生后，读取事件对应的epoll_data，该结构中包含了文件fd和数据地址，由于采用了mmap，程序可以直接读取数据。
int epoll_wait(int epfd, struct epoll_event * events, int maxevents, int timeout);

有人把epoll这种方式叫做同步非阻塞（NIO），因为用户线程需要不停地轮询，自己读取数据，看上去好像只有一个线程在做事情
也有人把这种方式叫做异步非阻塞（AIO），因为毕竟是内核线程负责扫描fd列表，并填充事件链表的
个人认为真正理想的异步非阻塞，应该是内核线程填充事件链表后，主动通知用户线程，或者调用应用程序事先注册的回调函数来处理数据，如果还需要用户线程不停的轮询来获取事件信息，就不是太完美了，所以也有不少人认为epoll是伪AIO，还是有道理的。


另外一个epoll的变化，是支持了边沿触发，以前select/poll中，每次遍历fd列表，发现fd可写、可读或异常后，就把bit置1（select）或返回对应事件（poll），
而在epoll中，同样支持这种方式，每次fd可写、可读或异常后，就写入事件到事件链表中，还支持只在事件发生变化时才写入事件链表，例如如果事件一直是可读，则只在第一次写入链表
业界把这两种方式分别叫做电平触发和边沿触发，像电信号（方波）一样，从高电平到低电平或低电平到高电平的“拐角”处的触发，叫做边沿触发，其他上下两个平面上的连续触发叫电平触发
epoll支持电平触发（Level Triggered）和边沿触发（Edge Triggered），默认为电平触发

转载：https://blog.csdn.net/lcx46/article/details/42006845 


