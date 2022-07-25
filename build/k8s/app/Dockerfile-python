FROM jybl/face

RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple

RUN python3 -m pip install --upgrade pip

RUN pip3 install -r requirements.txt

WORKDIR /app

ADD ./app.py /app
ADD ./requirements.txt /app

CMD ["python", "app.py"]