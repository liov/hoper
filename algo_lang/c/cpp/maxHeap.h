#include <stdio.h>
#define MAX_ELEMENTS 20
#define MAX_SIZE 20
#define HEAP_FULL(n) (MAX_ELEMENTS - 1 == n)
#define HEAP_EMPTY(n) (!n)
typedef struct {
    int key;
}element;
element heap[MAX_ELEMENTS];

void insert_max_heap(element item ,int *n){
    if(HEAP_FULL(*n)){
      return;
    }
    int i = ++(*n);
    for(;(i != 1) && (item.key > heap[i/2].key);i = i / 2){/// i ≠ 1是因为数组的第一个元素并没有保存堆结点
      heap[i] = heap[i/2];/// 这里其实和递归操作类似，就是去找父结点
    }
    heap[i] = item;
}

element delete_max_heap(int *n){
  int parent, child;
  element temp, item;
  temp = heap[--*n];
  item = heap[1];
  parent = 1,child=2;
  for(;child <= *n; child = child * 2){
   if( (child < *n) && heap[child].key < heap[child+1].key){/// 这一步是为了看当前结点是左子结点大还是右子结点大，然后选择较大的那个子结点
        child++;
      }
      if(temp.key >= heap[child].key){
        break;
      }
      heap[parent] = heap[child];///这就是上图中第二步和第三步中黄色部分操作
      parent = child;/// 这其实就是一个递归操作，让parent指向当前子树的根结点
   }
  heap[parent] = temp;
  return item;
}

void create_max_heap(void){
        int total = (*heap).key;
        /// 求倒数第一个非叶子结点
        int child = 2,parent = 1;
        for (int node = total/2; node>0; node--) {
            parent = node;
            child = 2*node;
            int max_node = 2*node+1;
            element temp = *(heap+parent);
            for (; child <= total; child *= 2,max_node = 2*parent+1) {
                if (child+1 <= total && (*(heap+child)).key < (*(heap+child+1)).key) {
                    child++;
                }
                if (temp.key > (*(heap+child)).key) {
                    break;
                }
                *(heap+parent) = *(heap+child);
                parent = child;
            }
            *(heap+parent) = temp;
        }
    }

/**
 *
 * @param heap  最大堆；
 * @param items 输入的数据源
 * @return 1成功，0失败
 */
int create_binary_tree(element *heap,int items[MAX_ELEMENTS]){
    int total;
    if (!items) {
        return 0;
    }
    element *temp = heap;
    heap++;
    for (total = 1; *items;total++,(heap)++,items = items + 1) {
        element ele = {*items};
        element temp_key = {total};
        *temp = temp_key;
        *heap = ele;
    }
    return 1;
}

///函数调用
void main(){
    int items[MAX_ELEMENTS] = {79,66,43,83,30,87,38,55,91,72,49,9};
    element *position = heap;
    create_binary_tree(position, items);
    for (int i = 0; (*(heap+i)).key > 0; i++) {
    printf("binary tree element is %d\n",(*(heap + i)).key);
    }
    create_max_heap();
    for (int i = 0; (*(heap+i)).key > 0; i++) {
    printf("heap element is %d\n",(*(heap + i)).key);
    }
}

void __swap(element *lhs,element *rhs){
    element temp = *lhs;
    *lhs = *rhs;
    *rhs = temp;
}

int create_binarytree(element *heap, int items[MAX_SIZE], int n){
    if (n <= 0) return 0;
    for (int i = 0; i < n; i++,heap++) {
        element value = {items[i]};
        *heap = value;
    }
    return 1;
}

void adapt_maxheap(element *heap ,int node ,int n){
    int parent = node - 1 < 0 ? 0 : node - 1;
    int child = 2 * parent + 1;/// 因为没有哨兵，所以在数组中的关系由原来的：parent = 2 * child => parent = 2 * child + 1
    int max_node = max_node = 2*parent+2 < n - 1 ? 2*parent+2 : n - 1;
    element temp = *(heap + parent);
    for (;child <= max_node; parent = child,child = child * 2 + 1,max_node = 2*parent+2 < n - 1 ? 2*parent+2 : n - 1) {
        if ((heap + child)->key <= (heap + child + 1)->key && child + 1 < n) {
            child++;
        }
        if ((heap + child)->key < temp.key) {
            break;
        }
        *(heap + parent) = *(heap + child);
    }
    *(heap + parent) = temp;
}

int create_maxheap(element *heap ,int n){

    for (int node = n/2; node > 0; node--) {
        adapt_maxheap(heap, node, n);
    }
    return 1;
}

void heap_sort(element *heap ,int n){
    ///创建一个最大堆
    create_maxheap(heap, n);
    ///进行排序过程
    int i = n - 1;
    while (i >= 0) {
        __swap(heap+0, heap + i);/// 将第一个和最后一个进行交换
        adapt_maxheap(heap, 0, i--);///将总的元素个数减一，适配成最大堆，这里只需要对首元素进行最大堆的操作
    }
}