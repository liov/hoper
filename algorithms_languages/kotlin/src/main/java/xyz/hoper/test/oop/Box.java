package xyz.hoper.test.oop;

public class Box {
    public Box(){
        System.out.println("创建box");
    }
    public Box Box(){
        return new Box();
    }
    public static void main(String[] args) {
        var box =new Box();
        box.Box();
    }
}
