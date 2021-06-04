#!/usr/bin/env python
# coding=utf-8
import re
import requests
from bs4 import BeautifulSoup
import os


def downloadpic(url):
    headers = {
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
        'Accept-Encoding': 'gzip, deflate',
        'Accept-Language': 'zh-CN,zh;q=0.9',
        'Connection': 'keep-alive',
        'Cookie': 'UM_distinctid=160c072721f36a-049309acceadc2-e323462-144000-160c0727220f67; CNZZDATA3866066=cnzz_eid%3D1829424698-1494676185-%26ntime%3D1494676185; bdshare_firstime=1515057214243; Hm_lvt_9a737a8572f89206db6e9c301695b55a=1515057214,1515074260,1515159455; Hm_lpvt_9a737a8572f89206db6e9c301695b55a=1515159455',
        'Host': 'img1.mm131.me',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36',
        'Referer': 'http://www.mm131.com/'
    }

    # url='http://www.mm131.com/xinggan/3561.html'

    r = requests.get(url, headers={'Accept-Encoding': ''})
    # r.encoding="gb2312"
    r.encoding = r.apparent_encoding
    html = r.text

    # 套图标题
    title = BeautifulSoup(html, 'lxml').find("h5").get_text()
    # 获取页码
    page = BeautifulSoup(html, 'lxml').find("span", {"class": "page-ch"}).get_text()
    print(page)
    pattern = re.compile('\d*')
    page = pattern.findall(page)[1]
    # 创建以套图标题为题的文件夹
    path = "E:\\pic\\" + title + page + 'P'
    isExists = os.path.exists(path)
    if not isExists:
        os.makedirs(path)

    # 获取第一张图片地址
    a = re.search(r'img alt=.* src="(.*?)" /', html, re.S)
    print(a.group(1))
    pic = requests.get(a.group(1), headers=headers)
    # 下载图片
    f = open(path + '\\' + '1.jpg', "wb")
    f.write(pic.content)
    f.close

    # 下载第一张以后的图
    after = int(page) + 1
    for i in range(2, after):
        # 改变地址结构
        url0 = url[:-5]
        url1 = url0 + '_' + str(i) + '.html'
        # print url1
        try:
            html = requests.get(url1).text
            a = re.search(r'img alt=.* src="(.*?)" /', html, re.S)
            pic = requests.get(a.group(1), headers=headers)
            print(a.group(1))
            f = open(path + '\\' + str(i) + ".jpg", "wb")
            f.write(pic.content)
            f.close
        except:
            pass


if __name__ == '__main__':
    '''
    url = 'http://www.mm131.com/xinggan/'
    html = requests.get(url).text
    urls = BeautifulSoup(html, 'lxml').find('dl', {'class': 'list-left public-box'}).findAll('a', {'target': '_blank'})
    for url in urls:
        url = url['href']
        print url
        downloadpic(url)
    
    for i in range(2, 10):
        print("第" + str(i) + "页")
        url = 'http://www.mm131.com/xinggan/list_6_' + str(i) + '.html'
        html = requests.get(url).text
        urls = BeautifulSoup(html, 'lxml').find('dl', {'class': 'list-left public-box'}).findAll('a', {'target': '_blank'})
        for url in urls:
            url = url['href']
            print url
            downloadpic(url)
    '''
    '''4612 4364 4202 4032,3939,3831'''
    for i in range(4773, 4623, -1):
        url = 'http://www.mm131.com/xinggan/' + str(i) + '.html'
        downloadpic(url)
        i = i-1

