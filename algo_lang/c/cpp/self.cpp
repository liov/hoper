#include <iostream>
#include<fstream>
using namespace std;

int main()
{
    ifstream infile("14打印14的源代码.cpp",ios::binary);
    char ch;
    while(infile.peek()!=EOF)
    {
        infile.read(&ch,sizeof(ch));
        cout<<ch;
    }
    cout<<endl;
    return 0;
}
