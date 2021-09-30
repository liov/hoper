#include <iostream>
#include <vector>
using namespace std;

typedef int I;

template <class T>
class MaxHeap{
    private:
        T *heap;
        int HeapSize;
        int MaxSize;
    public:
        bool MaxHeapInsert (T &x)
        bool MaxHeapDelete (T &x)
};
MaxHeap H;

template <class T>
void MaxHeapInit (MaxHeap<T> &H){
	for(int i = H.HeapSize/2; i>=1; i--)
	{
		H.heap[0] = H.heap[i];
		int son = i*2;
		while(son <= H.HeapSize)
		{
			if(son < H.HeapSize && H.heap[son] < H.heap[son+1])
				son++;
			if(H.heap[0] >= H.heap[son])
				break;
			else
			{
				H.heap[son/2] = H.heap[son];
				son *= 2;
			}
		}
		H.heap[son/2] = H.heap[0];
	}
}

template <class T>
bool MaxHeapInsert (MaxHeap<T> &H, T &x) {
	if(H.HeapSize == H.MaxSize)
		return false;
	int i = ++H.HeapSize;
	while(i!=1 && x>H.heap[i/2])
	{
		H.heap[i] = H.heap[i/2];
		i = i/2;
	}
	H.heap[i] = x;
	return true;
}

template <class T>
bool MaxHeapDelete (MaxHeap<T> &H, T &x) {
	if(H.HeapSize == 0)
		return false;
	x = H.heap[1];
	H.heap[0] = H.heap[H.HeapSize--];
	int i = 1, son = i*2; 
 
	while(son <= H.HeapSize)
	{
		if(son <= H.HeapSize && H.heap[0] < H.heap[son+1])
			son++;
		if(H.heap[0] >= H.heap[son])
			break;
		H.heap[i] = H.heap[son];
		i = son;
		son  = son*2;
	}
	H.heap[i] = H.heap[0];
	return true;
}