import asyncio
import re
import requests
from bs4 import BeautifulSoup
from aiohttp import ClientSession
import time


async def download(url):
    async with ClientSession() as session:
        async with session.get(url) as response:
            print(response.status)
            response = await response.read()
            print(response)
            print('Hello World:%s' % time.time())

if __name__ == '__main__':
    tasks = []
    loop = asyncio.get_event_loop()
    task = asyncio.ensure_future(download(r"https://{{host}}/viewthread.php?tid=368995"))
    tasks.append(task)
    loop.run_until_complete(asyncio.wait(tasks))