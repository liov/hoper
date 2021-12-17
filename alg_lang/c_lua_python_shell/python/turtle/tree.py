import random
import time

import turtle

n = 80.0

turtle.speed("fastest")
turtle.screensize(bg='seashell')
turtle.left(90)
turtle.forward(3 * n)
turtle.color("orange", "yellow")
turtle.begin_fill()
turtle.left(126)

for i in range(5):
    turtle.forward(n / 5)
    turtle.right(144)
    turtle.forward(n / 5)
    turtle.left(72)
turtle.end_fill()
turtle.right(126)

turtle.color("dark green")
turtle.backward(n * 4.8)


def tree(d, s):
    if d <= 0: return
    turtle.forward(s)
    tree(d - 1, s * .8)
    turtle.right(120)
    tree(d - 3, s * .5)
    turtle.right(120)
    tree(d - 3, s * .5)
    turtle.right(120)
    turtle.backward(s)


tree(15, n)
turtle.backward(n / 2)

for i in range(200):
    a = 200 - 400 * random.random()
    b = 10 - 20 * random.random()
    turtle.up()
    turtle.forward(b)
    turtle.left(90)
    turtle.forward(a)
    turtle.down()
    if random.randint(0, 1) == 0:
        turtle.color('tomato')
    else:
        turtle.color('wheat')
    turtle.circle(2)
    turtle.up()
    turtle.backward(a)
    turtle.right(90)
    turtle.backward(b)
time.sleep(60)
