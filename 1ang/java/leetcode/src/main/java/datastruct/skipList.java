package datastruct;


import java.util.Random;


class SkipList<K extends Comparable<? super K>, V> {

  private Node<K, V> head;
  private Node<K, V> tail;
  private int size = 0;
  private int level = 0;

  private static final double PROBABILITY = 0.5;
  private final K HEAD_KEY;
  private final K TAIL_KEY;
  private static final Random random = new Random();

  SkipList(K headKey, K tailKey) {
      HEAD_KEY = headKey;
      TAIL_KEY = tailKey;
  }


  public Node<K,V> get(K key) {
    var p = findNode(key);
   if (p.key == key){
     return  p;
   }
    return null;
  }

  /**
   * put方法有一些需要注意的步骤：
   * 1.如果put的key值在跳跃表中存在，则进行修改操作；
   * 2.如果put的key值在跳跃表中不存在，则需要进行新增节点的操作，并且需要由random随机数决定新加入的节点的高度（最大level）；
   * 3.当新添加的节点高度达到跳跃表的最大level，需要添加一个空白层（除了-oo和+oo没有别的节点）
   *
   * @param k
   * @param v
   */
  public void put(K k, V v) {
    System.out.println("添加key:"+k.toString());
    Node<K,V> p =findNode(k); //这里不用get是因为下面可能用到这个节点
    if(p!=null){
      System.out.println("找到P:"+p);
      if (p.key == k) {
        p.data = v;
        return;
      }
    }

    Node<K,V> q = new Node<>(k, v);
    insertNode(p, q);
    var currentLevel = 0;
    while (random.nextDouble() > PROBABILITY) {
      if (currentLevel >= level) addEmptyLevel();
      p = head;
      //创建 q的镜像变量（只存储k，不存储v，因为查找的时候会自动找最底层数据）
      Node<K,V> z = new Node<>(k, null);
      insertNode(p, z);
      z.down = q;
      q.up = z;
      //别忘了把指针移到上一层。
      q = z;
      currentLevel++;
      System.out.println("添加后"+this);
    }
    size++;
  }


  private void insertNode(Node<K,V> p, Node<K,V> q) {
    q.next = p.next;
    q.pre = p;
    p.next .pre = q;
    p.next = q;
  }

  private void addEmptyLevel() {
    Node<K,V> p1 = new Node<>(HEAD_KEY, null);
    Node<K,V> p2 = new Node<>(TAIL_KEY, null);
    p1.next = p2;
    p1.down = head;
    p2.pre = p1;
    p2.down = tail;
    head.up = p1;
    tail.up = p2;
    head = p1;
    tail = p2;
    level++;
  }

  //首先查找到包含key值的节点，将节点从链表中移除，接着如果有更高level的节点，则repeat这个操作即可。
  public V remove(K k) {
    Node<K,V> p = get(k);
    var oldV = p.data;
    Node<K,V> q;
    while (p != null) {
      q = p.next ;
      q.pre = p.pre;
      p.pre.next = q;
      p = p.up;
    }
    return oldV;
  }

  private Node<K,V> findNode(K key) {
    Node<K,V> p =head;
    while (true) {
      if (p.next != null && p.next.key.compareTo(key)<=0){
        p = p.next;
      }
      if (p.down != null){
        p = p.down ;
      } else if (p.next != null && p.next.key.compareTo(key) > 0){
        break;
      }
    }
    return p;
  }

  public String toString() {
    var p = head;
    while (p.down != null) {
      p = p.down;
    }
      var sb = new StringBuilder();
    while (p.next != null) {
      sb.append(p);
      p = p.next;
    }
    return sb.toString();
  }

  public static void main(String[] args) {
    var srl = new SkipList<Integer,Integer>(Integer.MIN_VALUE,Integer.MAX_VALUE);
    srl.put(2,5);
    srl.put(6,7);
    srl.put(9,0);
    srl.put(5,3);
    srl.put(3,5);

    System.out.println(srl.findNode(5));;
  }
}


