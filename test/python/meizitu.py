# coding: utf-8
import  requests
from  bs4 import  BeautifulSoup
import  re
import os
from hashlib import md5
import multiprocessing
import time
import lxml
from requests.exceptions import RequestException
#http请求头
Hostreferer = {
    'User-Agent':'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36',
    'Referer':'http://www.mzitu.com'
}
Picreferer = {
    'User-Agent':'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36',
    'Referer':'http://i.meizitu.net'
}
#此请求头破解盗链
def header(referer):
    headers = {
        'User-Agent':'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36',
        'Referer':referer,
    }
    return headers
def get_one_page(url):
    try:
        response=requests.get(url,headers=Hostreferer)
        response.encoding = 'utf-8'
        if response.status_code==200:
            return response.text
        return None
    except:
        print("请求地址出错 "+url)
        return None
def get_one_page_picreferer(url):
    try:
        response=requests.get(url,headers=Picreferer)
        response.encoding = 'utf-8'
        if response.status_code==200:
            return response.text
        return None
    except:
        print("请求地址出错 "+url)
        return None

def parse_main_page(html):
    soup=BeautifulSoup(html,'lxml',exclude_encodings='utf-8')
    #print(soup.contents)
    items=soup.select('li')
    # print(items)
    pattern=re.compile('><a href="(.*?)"',re.S)
    pattern2=re.compile('<span class="title">.*?>(.*?)</a>',re.S)
    urlgroup={}
    i=0
    for item in items:
        content=str(item.contents)
        i=i+1
        url=re.findall(pattern,content)
        # print(url[0])
        nmaes=re.findall(pattern2,content)
        urlgroup.__setitem__(nmaes[0], url[0])
        # print(nmaes)
    # print(urlgroup)
    return  urlgroup

def parse_one_page(html):
    soup = BeautifulSoup(html,'html.parser', exclude_encodings='utf-8')
    # print(soup.contents)
    pattern = re.compile('src="(.*?)"/></a>', re.S)
    items=re.search(pattern,str(soup.contents))
    #取得每一图片的地址.
    # print(items[1])
    #获取总页数
    max_pages=soup.find(id="opic").previous_sibling.text
    # print(max_pages)
    # print(items[1])
    url_list=[]

    for i in  range(1,int(max_pages)+1):
        one_url = items[1]
        url_list.append(one_url[:-5]+str(i)+".jpg")
    # print(url_list)
    return url_list
    # save_one_pic(items[1])
def download_one_pic(url):
    try:
        response=requests.get(url,headers=header(url))
        # response.encoding = 'utf-8'
        if response.status_code==200:
            print("正在下载"+url)
            # print("header==",header(url))
            # print(response.content)
            return response.content
        return None
    except:
        print("请求地址出错 "+url)
        return None
#保存每一张图片.
def save_one_pic(filepath,content):
    image_content=download_one_pic(content)
    jpg_path='{0}/{1}.{2}'.format( filepath,md5(content.encode('utf-8')).hexdigest(),"jpg") # md5(content.encode('utf-8')).hexdigest()
    print("正在保存图片 " + filepath)
    try:
        if not os.path.exists(jpg_path):
            with open(jpg_path,'wb') as f:
                f.write(image_content)
                f.close()
    except:
        print("保存图片出错"+filepath)

def main_run(offset):
    print("开始进程:",offset)
    url=""
    for i in range(0,9):
        if offset==0 and i<2:
            url="http://www.mmjpg.com"
        else:
            url="http://www.mmjpg.com/home/"+str(offset*10+i)
        html=get_one_page(url)
        #取得每页包含的子页list
        items=parse_main_page(html)
        for item in items.keys():
            #子页的地址
            itemhtml=get_one_page_picreferer(items[item])
            #子页地址的图片地址
            image_urls=parse_one_page(itemhtml)
            for per_url in  image_urls:
                filepath= item
                if not  os.path.exists(filepath):
                    os.makedirs(filepath)
                save_one_pic(filepath,per_url)
# content=download_one_pic("http://img.mmjpg.com/2018/1223/1.jpg")
#  save_one_pic("1","http://img.mmjpg.com/2018/1223/1.jpg")
if __name__ == '__main__':
    multiprocessing.freeze_support()
    # main()
    # group=[x for x in  range(2,80)]

    # for i in group:
    #     main(i)
    #     print("正在下载第",i,"页的图片")
    #     time.sleep(2)
    mainStart = time.time()  # 记录主进程开始的时间
    p =multiprocessing.Pool(8) # 开辟进程池
    for i in range(8):
        p.apply_async(main_run, args=(i,))  # 每个进程都调用run_proc函数，
        # args表示给该函数传递的参数。
    print('Waiting for all subprocesses done ...')
    p.close()  # 关闭进程池
    p.join()  # 等待开辟的所有进程执行完后，主进程才继续往下执行
    print('All subprocesses done')
    mainEnd = time.time()  # 记录主进程结束时间
    print('All process ran %0.2f seconds.' % (mainEnd - mainStart))  # 主进程执行时间