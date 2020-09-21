package xyz.hoper.test.innerclass;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/11 13:41
 * @description：匿名内部类
 * @modified By：
 */
public class Outer {

    private String string = "JYB";

    class Inner4 {
    }

    public class Inner1 {
        //内部类的成员
        private String string1 = "JYB1";

        Inner1() {
        }
    }

    static class Inner2 {
        String string2 = "JYB2";

        public Inner2() {
            System.out.println("静态内部类：" + string2);
        }
    }

    void inner3() {
        class Inner3 {
            String string3 = "JYB3";

            //电视宣传
            public void show() {
                System.out.println("内部局部类：" + string3);
            }
        }
    }

    public Inner getInner4(int age, String name) {
        return new Inner() {
            int age_;
            String name_;

            //构造代码块完成初始化工作
            {
                if (0 < age && age < 200) {
                    age_ = age;
                    name_ = name;
                }
            }

            public String getName() {
                return name_;
            }

            public int getAge() {
                return age_;
            }

            public String toString() {
                return this.age_ + this.name_;
            }

            public void work() {
                System.out.println("age:"+this.age_ +",name:"+this.name_);
            }
        };
    }

    public Outer getInner5(int age, String name) {
        return new Outer() {
            int age_;
            String name_;

            public String getName() {
                return name_;
            }

            public int getAge() {
                return age_;
            }

            public void work() {
                System.out.println("age:"+this.age_ +",name:"+this.name_);
            }
        };
    }

    private void work(){
        System.out.println(string);
    }

    private void outerTest(char ch) {

        Integer integer = 1;
        new Inner4() {
            void innerTest() {
                System.out.println(string);
                System.out.println(ch);
                System.out.println(integer);
            }
        }.innerTest();
    }

    public static void main(String[] args) {
        Outer outer = new Outer();
        Outer.Inner1 inner1 = outer.new Inner1();
        outer.outerTest('b');
        outer.getInner4(21, "chenssy").work();
        outer.getInner5(22, "chenssy").work();
    }
}

interface Inner {
    void work();
}
