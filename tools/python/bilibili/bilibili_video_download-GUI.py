# !/usr/bin/python
# -*- coding:utf-8 -*-
# time: 2019/07/02--08:12
__author__ = 'Henry'

from tools.python.bilibili.api import get_play_list, down_video
from tools.python.bilibili.util import combine_video, format_size

'''
项目: B站视频下载 - GUI版本
版本1: 加密API版,不需要加入cookie,直接即可下载1080p视频
20190422 - 增加多P视频单独下载其中一集的功能
20190702 - 增加视频多线程下载 速度大幅提升
20190711 - 增加GUI版本,可视化界面,操作更加友好
'''

import requests, time, hashlib, urllib.request, re, json
import imageio_ffmpeg
from moviepy.editor import *
import os, sys, threading



from tkinter import *
from tkinter import ttk
from tkinter import StringVar
root=Tk()
start_time = time.time()

# 将输出重定向到表格
def print(theText):
    msgbox.insert(END,theText+'\n')



# 下载视频
'''
 urllib.urlretrieve 的回调函数：
def callbackfunc(blocknum, blocksize, totalsize):
    @blocknum:  已经下载的数据块
    @blocksize: 数据块的大小
    @totalsize: 远程文件的大小
'''


def Schedule_cmd(blocknum, blocksize, totalsize):
    speed = (blocknum * blocksize) / (time.time() - start_time)
    # speed_str = " Speed: %.2f" % speed
    speed_str = " Speed: %s" % format_size(speed)
    recv_size = blocknum * blocksize

    # 设置下载进度条
    pervent = recv_size / totalsize
    percent_str = "%.2f%%" % (pervent * 100)
    download.coords(fill_line1,(0,0,pervent*465,23))
    root.update()
    pct.set(percent_str)



def Schedule(blocknum, blocksize, totalsize):
    speed = (blocknum * blocksize) / (time.time() - start_time)
    # speed_str = " Speed: %.2f" % speed
    speed_str = " Speed: %s" % format_size(speed)
    recv_size = blocknum * blocksize

    # 设置下载进度条
    f = sys.stdout
    pervent = recv_size / totalsize
    percent_str = "%.2f%%" % (pervent * 100)
    n = round(pervent * 50)
    s = ('#' * n).ljust(50, '-')
    print(percent_str.ljust(6, ' ') + '-' + speed_str)
    f.flush()
    time.sleep(2)
    # print('\r')


def do_prepare(inputStart,inputQuality):
    # 清空进度条
    download.coords(fill_line1,(0,0,0,23))
    pct.set('0.00%')
    root.update()
    # 清空文本栏
    msgbox.delete('1.0','end')
    start_time = time.time()
    # 用户输入av号或者视频链接地址
    print('*' * 30 + 'B站视频下载小助手' + '*' * 30)
    start = inputStart
    if start.isdigit() == True:  # 如果输入的是av号
        # 获取cid的api, 传入aid即可
        start_url = 'https://api.bilibili.com/x/web-interface/view?aid=' + start
    else:
        # https://www.bilibili.com/video/av46958874/?spm_id_from=333.334.b_63686965665f7265636f6d6d656e64.16
        start_url = 'https://api.bilibili.com/x/web-interface/view?aid=' + re.search(r'/av(\d+)/*', start).group(1)

    # 视频质量
    # <accept_format><![CDATA[flv,flv720,flv480,flv360]]></accept_format>
    # <accept_description><![CDATA[高清 1080P,高清 720P,清晰 480P,流畅 360P]]></accept_description>
    # <accept_quality><![CDATA[80,64,32,16]]></accept_quality>
    quality = inputQuality
    # 获取视频的cid,title
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36'
    }
    html = requests.get(start_url, headers=headers).json()
    data = html['data']
    cid_list = []
    if '?p=' in start:
        # 单独下载分P视频中的一集
        p = re.search(r'\?p=(\d+)',start).group(1)
        cid_list.append(data['pages'][int(p) - 1])
    else:
        # 如果p不存在就是全集下载
        cid_list = data['pages']
    # print(cid_list)
    # 创建线程池
    threadpool = []
    title_list = []
    for item in cid_list:
        cid = str(item['cid'])
        title = item['part']
        title = re.sub(r'[\/\\:*?"<>|]', '', title)  # 替换为空的
        print('[下载视频的cid]:' + cid)
        print('[下载视频的标题]:' + title)
        title_list.append(title)
        page = str(item['page'])
        referer_url = start_url + "/?p=" + page
        video_list = get_play_list(referer_url, cid, quality)
        start_time = time.time()
        # down_video(video_list, title, start_url, page)
        # 定义线程
        th = threading.Thread(target=down_video, args=(video_list, title, referer_url, page))
        # 将线程加入线程池
        threadpool.append(th)

    # 开始线程
    for th in threadpool:
        th.start()
    # 等待所有线程运行完毕
    for th in threadpool:
        th.join()
    
    # 最后合并视频
    combine_video(title_list)

    end_time = time.time()  # 结束时间
    print('下载总耗时%.2f秒,约%.2f分钟' % (end_time - start_time, int(end_time - start_time) / 60))

    # 如果是windows系统，下载完成后打开下载目录
    currentVideoPath = os.path.join(sys.path[0], 'bilibili_video')  # 当前目录作为下载目录
    if (sys.platform.startswith('win')):
        os.startfile(currentVideoPath)



def thread_it(func, *args):
    '''将函数打包进线程'''
    # 创建
    t = threading.Thread(target=func, args=args) 
    # 守护 !!!
    t.setDaemon(True) 
    # 启动
    t.start()


if __name__ == "__main__":
    # 设置标题
    root.title('B站视频下载小助手-GUI')

    # 各项输入栏和选择框
    inputStart = Entry(root,bd=4,width=600)
    labelStart=Label(root,text="请输入您要下载的B站av号或者视频链接地址:") # 地址输入
    labelStart.pack(anchor="w")
    inputStart.pack()
    labelQual = Label(root,text="请选择您要下载视频的清晰度") # 清晰度选择
    labelQual.pack(anchor="w")
    inputQual = ttk.Combobox(root,state="readonly")
    # 可供选择的表
    inputQual['value']=('1080P','720p','480p','360p')
    # 对应的转换字典
    keyTrans=dict()
    keyTrans['1080P']='80'
    keyTrans['720p']='64'
    keyTrans['480p']='32'
    keyTrans['360p']='16'
    # 初始值为720p
    inputQual.current(1)
    inputQual.pack()
    confirm = Button(root,text="开始下载",command=lambda:thread_it(do_prepare,inputStart.get(), keyTrans[inputQual.get()] ))
    msgbox = Text(root)
    msgbox.insert('1.0',"对于单P视频:直接传入B站av号或者视频链接地址\n(eg: 49842011或者https://www.bilibili.com/video/av49842011)\n对于多P视频:\n1.下载全集:直接传入B站av号或者视频链接地址\n(eg: 49842011或者https://www.bilibili.com/video/av49842011)\n2.下载其中一集:传入那一集的视频链接地址\n(eg: https://www.bilibili.com/video/av19516333/?p=2)")
    msgbox.pack()
    download=Canvas(root,width=465,height=23,bg="white")
    # 进度条的设置
    labelDownload=Label(root,text="下载进度")
    labelDownload.pack(anchor="w")
    download.pack()
    fill_line1 = download.create_rectangle(0, 0, 0, 23, width=0, fill="green")
    pct=StringVar()
    pct.set('0.0%')
    pctLabel = Label(root,textvariable=pct)
    pctLabel.pack()
    root.geometry("600x800")
    confirm.pack()
    # GUI主循环
    root.mainloop()
    
