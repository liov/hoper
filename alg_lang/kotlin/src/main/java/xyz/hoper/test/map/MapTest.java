package xyz.hoper.test.map;

import java.util.HashMap;

/**
 * @author ：lbyi
 * @date ：Created in 2019/7/31
 * @description：hashmap test
 */

/*
二叉查找树
        1， 左子树上所有的节点的值均小于或等于他的根节点的值
        2， 右子数上所有的节点的值均大于或等于他的根节点的值
        3， 左右子树也一定分别为二叉排序树
红黑树定义和性质
        红黑树是一种含有红黑结点并能自平衡的二叉查找树。它必须满足下面性质：

        性质1：每个节点要么是黑色，要么是红色。
        性质2：根节点是黑色。
        性质3：每个叶子节点（NIL）是黑色。
        性质4：每个红色结点的两个子结点一定都是黑色。
        性质5：任意一结点到每个叶子结点的路径都包含数量相同的黑结点。

        从性质5又可以推出：

        性质5.1：如果一个结点存在黑子结点，那么该结点肯定有两个子结点
*/

public class MapTest {

    public static void main(String[] args){
        var map = new HashMap<Integer,String>();
        map.put(5,"5");
        System.out.println(map);
        int h;
        System.out.println((h = Integer.valueOf(0).hashCode()) ^ (h >>> 16));
    }
}
