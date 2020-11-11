import requests

header = {
    'User-Agent': 'Python',
    'Cookie': ''
}

url = ''

ids = []

for eid in ids:
    response = requests.get(url, params={'empId': eid}, headers=header)
    print(response.url)
    response.encoding = 'utf-8'
    data = response.json()
    print(data)
    if data['status'] != 0:
        print('出错了', eid)
