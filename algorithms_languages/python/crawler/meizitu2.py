# -*- coding: utf-8-*-
__author__ = 'lihailin'
__email__ = '1501210931@qq.com'
__date__ = '2018年5月10日'

import os
import random
import time

import requests
from lxml import etree

import log
import logging

log.initLogConf()
logg = logging.getLogger(__file__)
import crawBase
import mongoDb

'''
妹子图网站的反爬虫机制, Referer
'''


class CrawMzitu(crawBase.CrawMzBase):

    def __init__(self):
        super(CrawMzitu, self).__init__('mzitu')
        self.startUri = 'http://www.mzitu.com/all/'
        self.baseUri = 'http://www.mzitu.com'
        self.mzituInfo = mongoDb.MongoDbOpreater('crawMeizi', 'mzitu')
        self.num = 1

    def getFirstPage(self, url):
        '''
        获取mzitu网站的第一页里的所有连接
        :param url:
        :return:
        '''
        self.headers = {
            'User-Agent': random.choice(self.user),
            'Referer': self.baseUri
        }  # 更新头
        r = requests.get(url, headers=self.headers, timeout=5).content  # 不要代理
        # print(r)
        sel = etree.HTML(r)
        subUris = sel.xpath('//p/a/@href')
        for url in subUris:
            print(url)
            self.getSecondPage(url)

    def getSecondPage(self, url):
        '''
        获取mzitu网站二级页面中的所有图片链接,并下载图片
        :return:
        '''
        self.headers = {'User-Agent': random.choice(self.user),
                        'Referer': self.baseUri
                        }  # 更新头
        r = requests.get(url, headers=self.headers, timeout=5).content  # 不要代理
        sel = etree.HTML(r)
        # 图片集的第一张图片url
        picUriFirst = sel.xpath('//p/a/img/@src')[0]
        total = sel.xpath("//div[@class='pagenavi']/a/span")[-2].text
        # 该题下的图片数
        picDesc = sel.xpath("//h2")[0].text
        # 图片集合描述,可用作文件夹名
        total = int(total)
        urlList = self.genPicUris(picUriFirst, total)
        # print(urlList)
        self.insertDb(urlList, picDesc)

    def genPicUris(self, url, total):
        '''
        url拼接生成的图片列表
        :param url:
        :param total: int, 一共有的图片数
        :param picDesc: str, 图片描述
        :return:
        '''
        urlBase = url[:-6]
        urlList = []
        for i in range(1, total + 1):
            if i < 10:
                t = urlBase + '0' + str(i) + '.jpg'
            else:
                t = urlBase + str(i) + '.jpg'
            urlList.append(t)
        # print(urlList)
        return urlList

    def insertDb(self, urlList, picDesc):
        '''
        图片保存和数据入库
        :param urlList:
        :return:
        '''
        picDesc = picDesc.replace(' ', '')
        # 创建文件夹保存图片
        dic = '%s/%s' % (self.picDictory, picDesc)
        if not os.path.exists(dic):
            os.system('mkdir ' + dic)

        for url in urlList:
            # 下载图片
            name = url.split('/')[-1]
            f = os.path.join(self.picDictory, picDesc, name)
            self.baseHeaders['Referer'] = self.baseUri
            self.logCrawPic(url, {}, f)
        # 数据入库
        rd = dict()
        rd['id'] = self.num
        insertTime = time.strftime('%Y/%m/%d %H:%M:%S', time.localtime(time.time()))
        rd['insertTime'] = insertTime
        rd['desc'] = picDesc
        rd['localPath'] = dic
        rd['urls'] = urlList
        self.num += 1
        self.mzituInfo.insert(rd)

    def crawlRun(self):
        self.getFirstPage(self.startUri)


if __name__ == '__main__':
    crawMzitu = CrawMzitu()
    crawMzitu.crawlRun()
