--Lua 实现链表

node={}
list=node

--初始化，构建一个空表
function init()
    list.data=0
    list.next=nil
end

--向链表的尾部添加数据
function addRear(d)
    node.next={}  --建立一个节点，相当于malloc一个节点
    node=node.next
    node.next=nil
    node.data=d
    list.data=list.data+1
end

--向链表的头部添加数据
function addHead(d)
    newNode={} --建立一个节点，相当于malloc一个节点
    newNode.data=d
    newNode.next=list.next
    list.next=newNode
    list.data=list.data+1
end

--第i个位置插入数据d i>=1
function insert(i,d)
    if i<1 then
        print('插入的位置不合法')
        return
    end

    local j,k,l=i-1,0,list  --找到第i个位置
    while k~=j do
        k=k+1
        l=l.next
        if not l.next then break end
    end
    --if k~=j then print("插入位置不合法") return end

    --开始插入
    newNode={}
    newNode.next=l.next
    newNode.data=d
    l.next=newNode
    list.data=list.data+1
end

--打印链表的每一个元素
function display()
    local l=list.next
    while l do
        print(l.data.." ")
        l=l.next
    end
    print('\n-- display ok --')
end

--判断链表是否为空
function is_empty()
    return list.data==0
end

--删除第i个位置的数据 i>=1 返回删除数据的内容
function delete(i)
    if i<1 then
        print('删除位置不合法')
        return
    end

    local j,k,l=i-1,0,list
    while k~=j do
        k=k+1
        l=l.next
        if not l.next then break end
    end

    --开始删除
    d=l.next.data
    t=l.next.next
    l.next=nil
    l.next=t
    list.data=list.data-1
    return d
end

--清理链表，操作完成后，链表还在，只不过为空，相当刚开始的初始化状态
function clear()
    if not list then
        print('链表不存在')
    end

    while true do
        firstNode=list.next
        if not firstNode then  --表示链表已为空表了
            break
        end
        t=firstNode.next
        list.next=nil
        list.next=t
    end
    list.data=0
    print('-- clear ok --')
end

--销毁链表
function destory()
    clear() --先清除链表
    list=nil
end

--获取第i个元素i>1的值
function getData(i)
    if not list then
        print('链表不存在')
        return
    end
    if i<1 then
        print('位置不合法')
        return
    end

    local l=list.next  --指向第一个元素
    local k=1
    while l do
        if k==i then
            return l.data
        end
        l=l.next
        k=k+1
    end

    print('位置不合法')
end

--获取链表的长度
function getLen()
    if not list then
        print('链表不存在')
        return
    end
    return list.data
end

--主方法
function main()
    init()
    addRear(5)
    addRear(7)
    addRear(10)
    addHead(1)
    addHead(2)
    insert(2,4)
    display()


    print('输入你要删除的元素的位置:')
    pos=io.read('*number')
    ret=delete(pos)
    if not ret then
        print('删除失败')
    else
        print('你要删除的元素是:'..ret)
    end
    print('删除后的链表的内容为:')
    display()

    print('输入你想要得到的元素的位置:')
    pos=io.read('*number')
    print('第'..pos..'个元素内容是:'..getData(pos))
    print('链表的长度为:'..getLen())

    destory() --销毁链表
    print('-- main end --')
end

--程序的入口
main()