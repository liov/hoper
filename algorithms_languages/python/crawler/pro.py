import asyncio
import random
import traceback

import aiofiles
import aiohttp
from bs4 import BeautifulSoup
import lxml

user_agent_list = ["Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 "
                   "Safari/537.36",
                   "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 "
                   "Safari/537.36",
                   "Mozilla/5.0 (Windows NT 10.0; …) Gecko/20100101 Firefox/61.0",
                   "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 "
                   "Safari/537.36",
                   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 "
                   "Safari/537.36",
                   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 "
                   "Safari/537.36",
                   "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
                   "Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10.5; en-US; rv:1.9.2.15) Gecko/20110303 Firefox/3.6.15",
                   'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                   'Chrome/63.0.3239.132 Safari/537.36'
                   ]

headers = {
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
    'Accept-Encoding': 'gzip, deflate',
    'Accept-Language': 'zh-CN,zh;q=0.9',
    'Connection': 'keep-alive',
}


async def fetch(session, url, params):
    async with session.get(url, params=params) as response:
        response = await response.text(encoding="utf-8")
        print(response)
        urls = parse_page(response)
        tasks = []
        headers['User-Agent'] = random.choice(user_agent_list)
        for url in urls:
            tasks.append(asyncio.create_task(download(session, url)))
        await asyncio.wait(tasks)


async def req(url, params):
    headers['User-Agent'] = random.choice(user_agent_list)
    async with aiohttp.ClientSession(headers=headers, timeout=aiohttp.ClientTimeout(total=5 * 120)) as session:
        tasks = []
        tasks.append(asyncio.create_task(fetch(session, url, params)))
        await asyncio.wait(tasks)


async def download(session, url):
    try:
        path = 'F:\\pic\\' + url.split(r'//', 2)[2]
        headers['User-Agent'] = random.choice(user_agent_list)
        async with session.get(url) as response:
            async with aiofiles.open(path, 'wb') as f:
                await f.write(await response.read())
                print("下载成功：", path)
    except Exception as e:
        print("下载失败：", url, e)
        await download(session, url)


def parse_page(html):
    soup = BeautifulSoup(html, 'lxml', exclude_encodings='utf-8')
    items = soup.find_all('img', src="images/common/none.gif")
    urls = []
    for item in items:
        urls.append(item.attrs['file'])
    return urls


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(req(r"https://f1113.wonderfulday27.live/viewthread.php", {"tid": 368995}))
