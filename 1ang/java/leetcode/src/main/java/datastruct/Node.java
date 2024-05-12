package datastruct;

import java.util.List;

/**
 * @Description TODO
 * @Date 2024/3/24 21:39
 * @Created by lbyi
 */
public class Node <K extends Comparable<? super K>,V>{
  Node(K key,V data){
    this.key = key;
    this.data = data;
  }

  K key;
  V data ;
    int level ;
  List<Node<K,V>> forwards;

  Node<K,V> pre;
  Node<K,V> next;
  Node<K,V> up;
  Node<K,V> down; //上下左右四个节点，pre和up存在的意义在于 "升层"的时候需要查找相邻节点

  @Override
  public String toString() {
    return "($key,$value)";
  }
}
