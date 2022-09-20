#![feature(str_as_mut_ptr)]

/*在多任务操作系统中的每一个进程都运行在一个属于它自己的内存沙盘中。这个沙盘就是虚拟地址空间（virtual address space），在32位模式下它总是一个4GB的内存地址块。
这些虚拟地址通过页表（page table）映射到物理内存，页表由操作系统维护并被处理器引用。每一个进程拥有一套属于它自己的页表，但是还有一个隐情。
只要虚拟地址被使能，那么它就会作用于这台机器上运行的所有软件，包括内核本身。因此一部分虚拟地址必须保留给内核使用：

这并不意味着内核使用了那么多的物理内存，仅表示它可支配这么大的地址空间，可根据内核需要，将其映射到物理内存。
内核空间在页表中拥有较高的特权级（ring 2或以下），因此只要用户态的程序试图访问这些页，就会导致一个页错误（page fault）。
在Linux中，内核空间是持续存在的，并且在所有进程中都映射到同样的物理内存。内核代码和数据总是可寻址的，随时准备处理中断和系统调用。
与此相反，用户模式地址空间的映射随进程切换的发生而不断变化：

地址空间的随机排布方式逐渐流行起来。Linux通过对栈、内存映射段、堆的起始地址加上随机的偏移量来打乱布局。

通过不断向栈中压入的数据，超出其容量就有会耗尽栈所对应的内存区域。
这将触发一个页故障（page fault），并被Linux的expand_stack()处理，它会调用acct_stack_growth()来检查是否还有合适的地方用于栈的增长。
如果栈的大小低于RLIMIT_STACK（通常是8MB），那么一般情况下栈会被加长，程序继续愉快的运行，感觉不到发生了什么事情。这是一种将栈扩展至所需大小的常规机制。
然而，如果达到了最大的栈空间大小，就会栈溢出（stack overflow），程序收到一个段错误（Segmentation Fault）。
当映射了的栈区域扩展到所需的大小后，它就不会再收缩回去，即使栈不那么满了。

动态栈增长是唯一一种访问未映射内存区域（图中白色区域）而被允许的情形。
其它任何对未映射内存区域的访问都会触发页故障，从而导致段错误。一些被映射的区域是只读的，因此企图写这些区域也会导致段错误。

在栈的下方，是我们的内存映射段。此处，内核将文件的内容直接映射到内存。
任何应用程序都可以通过Linux的mmap()系统调用（实现）或Windows的CreateFileMapping() / MapViewOfFile()请求这种映射。
内存映射是一种方便高效的文件I/O方式，所以它被用于加载动态库。创建一个不对应于任何文件的匿名内存映射也是可能的，此方法用于存放程序的数据。
在Linux中，如果你通过malloc()请求一大块内存，C运行库将会创建这样一个匿名映射而不是使用堆内存。‘大块’意味着比MMAP_THRESHOLD还大，缺省是128KB，可以通过mallopt()调整。

说到堆，它是接下来的一块地址空间。与栈一样，堆用于运行时内存分配；但不同点是，堆用于存储那些生存期与函数调用无关的数据。
大部分语言都提供了堆管理功能。因此，满足内存请求就成了语言运行时库及内核共同的任务。
在C语言中，堆分配的接口是malloc()系列函数，而在具有垃圾收集功能的语言（如C#）中，此接口是new关键字。

如果堆中有足够的空间来满足内存请求，它就可以被语言运行时库处理而不需要内核参与。
否则，堆会被扩大，通过brk()系统调用（实现）来分配请求所需的内存块。堆管理是很复杂的，需要精细的算法，应付我们程序中杂乱的分配模式，优化速度和内存使用效率。处理一个堆请求所需的时间会大幅度的变动。
实时系统通过特殊目的分配器来解决这个问题。堆也可能会变得零零碎碎。

最后，我们来看看最底部的内存段：BSS，数据段，代码段。在C语言中，BSS和数据段保存的都是静态（全局）变量的内容。
区别在于BSS保存的是未被初始化的静态变量内容，它们的值不是直接在程序的源代码中设定的。
BSS内存区域是匿名的：它不映射到任何文件。如果你写static int cntActiveUsers，则cntActiveUsers的内容就会保存在BSS中。

另一方面，数据段保存在源代码中已经初始化了的静态变量内容。这个内存区域不是匿名的。
它映射了一部分的程序二进制镜像，也就是源代码中指定了初始值的静态变量。
所以，如果你写static int cntWorkerBees = 10，则cntWorkerBees的内容就保存在数据段中了，而且初始值为10。
尽管数据段映射了一个文件，但它是一个私有内存映射，这意味着更改此处的内存不会影响到被映射的文件。
也必须如此，否则给全局变量赋值将会改动你硬盘上的二进制镜像，这是不可想象的。

1、栈区（stack）— 程序运行时由编译器自动分配，存放函数的参数值，局部变量的值等。其操作方式类似于数据结构中的栈。程序结束时由编译器自动释放。

2、堆区（heap） — 在内存开辟另一块存储区域。一般由程序员分配释放， 若程序员不释放，程序结束时可能由OS回收 。注意它与数据结构中的堆是两回事，分配方式倒是类似于链表，呵呵。用malloc, calloc, realloc等分配内存的函数分配得到的就是在堆上

3、全局区（静态区）（static）—编译器编译时即分配内存。全局变量和静态变量的存储是放在一块的。对于C语言初始化的全局变量和静态变量在一块区域， 未初始化的全局变量和未初始化的静态变量在相邻的另一块区域。而C++则没有这个区别 - 程序结束后由系统释放

4、文字常量区 —常量字符串就是放在这里的。 程序结束后由系统释放

5、程序代码区—存放函数体的二进制代码。

1.内存分配方式 
　　内存分配方式有三种：
　　[1]从静态存储区域分配。内存在程序编译的时候就已经分配好，这块内存在程序的整个运行期间都存在。例如全局变量，static变量。
　　[2]在栈上创建。在执行函数时，函数内局部变量的存储单元都可以在栈上创建，函数执行结束时这些存储单元自动被释放。栈内存分配运算内置于处理器的指令集中，效率很高，但是分配的内存容量有限。
　　[3]从堆上分配，亦称动态内存分配。程序在运行的时候用malloc或new申请任意多少的内存，程序员自己负责在何时用free或delete释放内存。动态内存的生存期由程序员决定，使用非常灵活，但如果在堆上分配了空间，就有责任回收它，否则运行的程序会出现内存泄漏，频繁地分配和释放不同大小的堆空间将会产生堆内碎块。
2.程序的内存空间 
　　一个程序将操作系统分配给其运行的内存块分为4个区域

　　1、栈区（stack）—　 由编译器自动分配释放 ，存放为运行函数而分配的局部变量、函数参数、返回数据、返回地址等。其操作方式类似于数据结构中的栈。
　　2、堆区（heap） —　 一般由程序员分配释放， 若程序员不释放，程序结束时可能由OS回收 。分配方式类似于链表。
　　3、全局区（静态区）（static）—存放全局变量、静态数据、常量。程序结束后由系统释放。
　　4、文字常量区 —常量字符串就是放在这里的。 程序结束后由系统释放。
　　5、程序代码区—存放函数体（类成员函数和全局函数）的二进制代码。

2.在C++中，内存分成5个区，他们分别是堆、栈、自由存储区、全局/静态存储区和常量存储区 
（1）.栈，就是那些由编译器在需要的时候分配，在不需要的时候自动清楚的变量的存储区。里面的变量通常是局部变量、函数参数等。 
（2）.自由存储区，就是那些由new分配的内存块，他们的释放编译器不去管，由我们的应用程序去控制，一般一个new就要对应一个delete。如果程序员没有释放掉，那么在程序结束后，操作系统会自动回收。 
（3）.堆，就是那些由malloc等分配的内存块，他和堆是十分相似的，不过它是用free来结束自己的生命的。 
（4）.全局/静态存储区，全局变量和静态变量被分配到同一块内存中，在以前的C语言中，全局变量又分为初始化的和未初始化的，在C++里面没有这个区分了，他们共同占用同一块内存区。 
（5）.常量存储区，这是一块比较特殊的存储区，他们里面存放的是常量，不允许修改（当然，你要通过非正当手段也可以修改） */


//静态字符不可修改，静态字符应该和静态数值放的不在一个静态内存区，后定义的静态数值的地址是要更小的,强行修改系统报错
//确确实实是存在一块不可修改的内存区的，即使是静态数值，不标mut，也不可以改
static mut HELLO_WORLD: &str = "JING_TAI_BIAN_LIANG__Hello, world!__JING_TAI_BIAN_LIANG";
//UBA:126,UBA的地址：0x47a020,UBA2的地址:0x47b747,HELLO_WORLD的地址：0x47b710
//0x47axxx可改，0x47bxxx不可改
static mut UBA: u8 = 125u8;
static UBA2: u8 = 125u8;

fn main() {
    unsafe {
        //静态值，放在代码区，不可修改
        //let i = 12;//可寻址，可修改
/*        let x = &12;
        let x_addr = x as *const i32 as *mut i32; //let x = &12;不可修改
        println!("i的指针：{:p},x的指针：{:p}",&x as *const &i32, x_addr);
        *x_addr = 13;
        println!("x的值：{}，i的值：{}", x,*x_addr);*/
        let a1 = "hello word";
        let a2 = "hello word";
        println!("相同字符串字面值是否一个地址：{}", a1.as_ptr() == a2.as_ptr());


        let ptr = HELLO_WORLD.as_ptr();
        println!("HELLO_WORLD的地址：{:p}", ptr);
        HELLO_WORLD = "World, hello!";
        /*        let mut_ptr = HELLO_WORLD.as_ptr() as *mut i32;
                *mut_ptr = 1196312906;*/
        println!("HELLO_WORLD的地址：{:p},值：{},首字母：{:?}", HELLO_WORLD.as_ptr(), HELLO_WORLD, *ptr as char);


        //let mut a = "a var of a";
        //let u8a = a.as_ptr() as *mut u8;
        //栈内str，可修改
        let a_bytes = [b'a', b' ', b'v', b'a', b'r', b' ', b'o', b'f', b' ', b'a'];

        let a = std::str::from_utf8_unchecked(&a_bytes);


        let addr_str_a = a as *const str as *mut u8;
        let addr_a = &a as *const &str;
        let ua = b"c var of c";

        let u = 72u8;
        let u_a = &u as *const u8 as *mut u8;


        println!("u8a:{:?}", addr_str_a);
        //*u8a = [b'c',b' ',b'v',b'a',b'r',b' ',b'o',b'f',b' ',b'c'];
        //*u8a=b'H';
        let uba_ptr = &UBA as *const u8 as *mut u8;
        *uba_ptr += 1;
        *u_a += 1;
        println!("u_a:{:?}", *u_a);
        //*u8a_asu8+=1;
        addr_str_a.copy_from(ua.as_ptr(), 10);
        println!("UBA:{:?},UBA的地址：{:p},UBA2的地址:{:p},a的值：{}", *uba_ptr, &UBA as *const u8, &UBA2 as *const u8, a);

        println!("a的地址：{:p},a的底层地址：{:p},a的值：{}", addr_a, addr_str_a, a);


        let mut b = a;
        let addr_b = &b as *const &str;
        println!("b的地址：{:p}", addr_b);
        b.replace("a", "b");
        let addr_b = &b as *const &str;

        println!("b的地址：{:p},b的值：{}", addr_b, *addr_b);


        let c = String::from("c var of c");
        let addr_c = &c as *const String;
        println!("c的地址：{:p}", addr_c);
        let d = c;
        let addr_d = &d as *const String;

        println!("d的地址：{:p},c的值：{}", addr_d, *addr_c);


        let mut e = [1, 1, 1, 1, 1, 1, 1, 1, 1];
        let f = &mut e[3..];
        f[0] = 0;
        println!("e的值：{:?}", e);


        let user = User { name: String::from("jyb"), phone: 196206 };
        let addr_user = &user as *const User;
        println!("user的地址：{:p}", addr_user);
        pass_by_value_ref_user(&user);
        pass_by_value_user(user);
        //奇怪，不同终端表现不一致
        println!("原来的user：{:?}", *addr_user);
    }
}


fn pass_by_value_ref_user(user: &User) {
    let addr_user = &user as *const &User;
    println!("借用函数参数的的地址：{:p}", addr_user);
}

fn pass_by_value_user(mut user: User) {
    let addr_user = &user as *const User;
    user.phone = 2;
    println!("所有权函数参数的的地址：{:p},user的值：{:?}", addr_user, user);
    //std::mem::forget(user);
}

#[derive(Debug)]
struct User {
    name: String,
    phone: i32,
}
