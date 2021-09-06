/* 首先介绍一个函数VirtualProtectEx，它用来改变一个进程的虚拟地址中特定页里的某一区域的保护属性，这句话有些咬嘴，直接从MSDN中翻译过来的，简单来说就是改变某一进程中虚拟地址的保护属性，如果以前是只读的，那改变属性为PAGE_EXECUTE_READWRITE后，就可以更改这部分内存了。

具体看它的实现

BOOL WINAPI VirtualProtectEx(
  __in   HANDLE hProcess,
  __in   LPVOID lpAddress,
  __in   SIZE_T dwSize,
  __in   DWORD flNewProtect,
  __out  PDWORD lpflOldProtect
);

第一个参数是进程的句柄，这个句柄可以使用由CreateProcess()函数得到的PROCESS_INFORMATION结构中的hProcess成员，CreateProcess（）这个函数相信大家用的很多了，我就不详细介绍了。如果仅仅知道hProcess的ID，那么可以通过OpenProcess（）函数得到，OpenProcess（）有三个参数，最后一个参数就是进程的ID，它的返回值为进程的句柄。如果想得到线程的句柄，同样可以采用这两种方式，利用结构体PROCESS_INFORMATION重的hThread成员或使用函数OpenThread（）。

第二个参数lpAddress就是页中要改变保护属性的地址。第三个参数dwSize改变保护属性区域的大小，以字节为单位。通过这两个变量就可以确定要改变包括属性的区域了。

第四个参数fNewProtect是要改变的保护属性。包含PAGE_EXECUTE_READWRITE，PAGE_READONLY等，但先看第五个参数lpfOldProtect，是用来保存为改变之前的内存区域的保护属性，也就是说我们可以先备份之前的保护属性，然后将属性更改为PAGE_EXECUTE_READWRITE，然后对内存区域进行一些操作，然后再将内区域更改回来。

 

接下来看ReadProcessMemory（）函数

BOOL WINAPI ReadProcessMemory(
  __in   HANDLE hProcess,
  __in   LPCVOID lpBaseAddress,
  __out  LPVOID lpBuffer,
  __in   SIZE_T nSize,
  __out  SIZE_T *lpNumberOfBytesRead
);

第一个参数hProcess为进程的句柄，第二个参数lpBaseAddress为要读取内容的基地址，当然你事先要了解你要读取内容的基地址。第三个参数lpBuffer为接收所读取内容的基地址。第四个参数为要读取内容的大小，以字节为单位。最后一个参数lpNumberOfBytesRead为接收读取内容的buffer中收到的字节数。一般都设为NULL，这个该参数将被忽略掉。

使用ReadProcessMemory（）函数，可以获得该进程内存空间中的信息，或是用于监测进程的执行情况，或是将进程内的数据备份，然后调用writeProcessMemory（）进行修改，必要时再还原该进程的数据。

BOOL WINAPI WriteProcessMemory(
  __in   HANDLE hProcess,
  __in   LPVOID lpBaseAddress,
  __in   LPCVOID lpBuffer,
  __in   SIZE_T nSize,
  __out  SIZE_T *lpNumberOfBytesWritten
);

该函数与readProcessMemory的参数基本相同，只不过第三个参数lpBuffer为要写入进程的内容。

将内存的内容更改后，别忘了调用VirtualProtectEx（）恢复内存原来的保护属性。*/


#include <iostream>
#include <Windows.h>
#include <string>
using namespace std;

using namespace std;

void memoryWrite()
{
    HWND hWnd = FindWindow(0, TEXT("Window Name"));
  	if(hWnd == 0)
	{
    	MessageBox(0, L"Error cannot find window.", L"Error", MB_OK|MB_ICONERROR);
  	}
	else
	{
    	DWORD proccess_ID;
    	GetWindowThreadProcessId(hWnd, &proccess_ID);
    	HANDLE hProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, proccess_ID);
    	if(!hProcess)
		{
      		MessageBox(0, L"Could not open the process!", L"Error!", MB_OK|MB_ICONERROR);
    	}
		else
		{
		int newdata = 0;
		std::cout << "Pleas enter the new value: ";
		std::cin >> newdata;
     	DWORD newdatasize = sizeof(newdata);
      	if(WriteProcessMemory(hProcess, (LPVOID)0x235EC8, &newdata, newdatasize, NULL)) // newdatasize = 4 byte
		{
			MessageBox(NULL, L"WriteProcessMemory worked.", L"Success", MB_OK + MB_ICONINFORMATION);
		}
		else
		{
			MessageBox(NULL, L"Error cannot WriteProcessMemory!", L"Error", MB_OK + MB_ICONERROR);
		}
		CloseHandle(hProcess);
		}
  	}
}

int main(){
	string agree; // terms and conditions yes or no
	cout << "You must agree to our terms and conditions first (Y/N)? ";
	cin >> agree;
	if(agree =="y" || agree=="Y")
	{
		system("CLS"); // clear the screen
		memoryWrite(); // function
	}
	else if(agree =="n" || agree=="N")
	{
		MessageBox(HWND_DESKTOP,L"To use this software you must agree to our terms and conditions\nYou can read our terms and conditions online from this site below\n\nhttps://termsofuse.com\n",L"Error",MB_ICONINFORMATION| MB_OK);
	}
	else
	{
		cout << "Incorrect choice, pleas try again!" << endl;
	}
	return 0;
}