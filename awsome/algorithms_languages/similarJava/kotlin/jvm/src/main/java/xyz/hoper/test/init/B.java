package xyz.hoper.test.init;

public class B extends A {

	public B() {
		System.out.println("B的构造方法 i=" + i);
	}

	int i = 1000;

	public void method() {
		System.out.println("test.init.B 的 method i = " + i);
	}

	public static void main(String[] args) {
		new B();
	}
}
