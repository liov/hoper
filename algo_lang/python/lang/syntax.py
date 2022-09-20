import datetime
import functools
from io import StringIO
import json


def fib(v):
    if v < 2:
        return 1
    else:
        return fib(v - 1) + fib(v - 2)


def foo(l):
    for v in l:
        print(v)


def bar(d):
    if 'foo' in d:
        print(d['foo'])


# 命名关键字参数
def foo1(*vs, f=sum, **d):
    print(f(vs))
    if 'foo' in d:
        print(d['foo'])


# 限制关键字参数的名字
def foo2(*, f):
    print(f)


# python有函数定义先后之分
def log(text):
    def decorator(func):
        def wrapper(*args, **kw):
            print('%s %s():' % (text, func.__name__))
            return func(*args, **kw)

        return wrapper

    return decorator


@log('test')
def f1():
    return lambda x: x * x


@log('time')
def now():
    print(datetime.datetime.now())


int2 = functools.partial(int, base=2)

if __name__ == '__main__':
    print(fib(10))
    # generator
    foo(map(f1(), [x * x for x in range(10)]))
    foo1(1, 2, 3, 4, 5, 6, 7, 8, foo=1)
    now()
    print(int2('10010010'))
    try:
        print('try...')
        r = 10 / int('a')
        print('result:', r)
    except ValueError as e:
        print('ValueError:', e)
    except ZeroDivisionError as e:
        print('ZeroDivisionError:', e)
    finally:
        print('finally...')
        print('END')
    with open('./complex.py', 'r') as f:
        print(f.read())
    f = StringIO('Hello!\nHi!\nGoodbye!')
    while True:
        s = f.readline()
        if s == '':
            break
        print(s.strip())
    d = dict(name='Bob', age=20, score=88)
    print(json.dumps(d))


# 因为python的多线程与本地线程一对一，但是因为GIL，同一时间只有一个线程在跑
# 如果想在 CPython 中使用真正的并行代码，只有多进程，但开销大。
# 3.7和3.8的提案貌似都没有真正解决并发问题，貌似是可以低调率并行了
