/* 所谓DLL注入就是将一个DLL放进某个进程的地址空间里，让它成为那个进程的一部分。要实现DLL注入，首先须要打开目标进程。
　　hRemoteProcess = OpenProcess(　PROCESS_CREATE_THREAD | //同意远程创建线程
　　PROCESS_VM_OPERATION　| //同意远程VM操作
　　PROCESS_VM_WRITE,　//同意远程VM写
　　FALSE, dwRemoteProcessId )
　　因为我们后面须要写入远程进程的内存地址空间并建立远程线程，所以须要申请足够的权限（PROCESS_CREATE_THREAD、VM_OPERATION、VM_WRITE）。
　　假设进程打不开，以后的操作就别想了。进程打开后，就能够建立远线程了，只是别急，先想想这个远线程的线程函数是什么？我们的目的是注入一个DLL。并且我们知道用LoadLibrary能够载入一个DLL到本进程的地址空间。于是，自然会想到假设能够在目标进程中调用LoadLibrary，不就能够把DLL载入到目标进程的地址空间了吗？对！就是这样。远线程就在这儿用了一次，建立的远线程的线程函数就是LoadLibrary，而參数就是要注入的DLL的文件名称。(这里须要自己想一想，注意到了吗，线程函数ThreadProc和LoadLibrary函数很相似，返回值，參数个数都一样) 另一个问题，LoadLibrary这个函数的地址在哪儿？或许你会说，这个简单，GetProcAddress就能够得出。于是代码就出来了。
　　char *pszLibFileRemote="my.dll";
　　PTHREAD_START_ROUTINE pfnStartAddr = (PTHREAD_START_ROUTINE)GetProcAddress(GetModuleHandle("Kernel32"), "LoadLibraryA");
　　CreateRemoteThread( hRemoteProcess, NULL, 0, pfnStartAddr, pszLibFileRemote, 0, NULL);
　　可是不正确！不要忘了，这是远线程，不是在你的进程里，而pszLibFileRemote指向的是你的进程里的数据，到了目标进程，这个指针都不知道指向哪儿去了，相同pfnStartAddr这个地址上的代码到了目标进程里也不知道是什么了，不知道是不是你想要的LoadLibraryA了。可是，问题总是能够解决的，Windows有些非常强大的API函数，他们能够在目标进程里分配内存，能够将你的进程中的数据复制到目标进程中。因此pszLibFileRemote的问题能够攻克了。
　　char *pszLibFileName="hack.dll";//注意，这个一定要是全路径文件名称，除非它在系统文件夹里；原因大家自己想想。
　　//计算DLL路径名须要的内存空间
　　int cb = (1 + lstrlenA(pszLibFileName)) * sizeof(char);
　　//使用VirtualAllocEx函数在远程进程的内存地址空间分配DLL文件名称缓冲区
　　pszLibFileRemote = (char *) VirtualAllocEx( hRemoteProcess, NULL, cb, MEM_COMMIT, PAGE_READWRITE);
　　//使用WriteProcessMemory函数将DLL的路径名拷贝到远程进程的内存空间
　　iReturnCode = WriteProcessMemory(hRemoteProcess, pszLibFileRemote, (PVOID) pszLibFileName, cb, NULL);
　　OK，如今目标进程也认识pszLibFileRemote了，可是pfnStartAddr好像不好办，我怎么可能知道LoadLibraryA在目标进程中的地址呢？事实上Windows为我们攻克了这个问题，LoadLibraryA这个函数是在Kernel32.dll这个核心DLL里的，而这个DLL非常特殊，无论对于哪个进程，Windows总是把它载入到同样的地址上去。因此你的进程中LoadLibraryA的地址和目标进程中LoadLibraryA的地址是同样的(事实上，这个DLL里的全部函数都是如此)。至此，DLL注入结束了。
 
  */


///#include "stdafx.h"
#include "windows.h"
#include "tchar.h"
#include "stdio.h"

BOOL InjectDll(DWORD dwPID, LPCTSTR szDllPath)
{
	HANDLE hProcess = NULL;
	HANDLE hThread = NULL;
	HMODULE hMod = NULL;
	LPVOID pRemoteBuf = NULL; //存储dll路径字符串的起始地址
	DWORD dwBufSize = (DWORD)(_tcslen(szDllPath)+1)*sizeof(TCHAR); // dll路径字符串的大小
	LPTHREAD_START_ROUTINE pThreadProc; // 存储LoadLibrary函数的地址


	// 使用dwPID获取目标进程句柄
	if(!(hProcess = OpenProcess(PROCESS_ALL_ACCESS,FALSE,dwPID)))
	{
		_tprintf(L"OpenProcess(%d) failed!!![%d]\n",dwPID,GetLastError());
		return FALSE;
	}

	// 在目标进程notepad.exe内存中分配szDLLName大小的内存
	pRemoteBuf = VirtualAllocEx(hProcess,NULL,dwBufSize,MEM_COMMIT,PAGE_READWRITE);

	// 将myhack.dll路径写入分配的内存
	WriteProcessMemory(hProcess,pRemoteBuf,(LPVOID)szDllPath,dwBufSize,NULL);

	// 获取LoadLibraryW() API的地址
	hMod = GetModuleHandle(L"kernel32.dll");
	pThreadProc = (LPTHREAD_START_ROUTINE)GetProcAddress(hMod,"LoadLibraryW");

	// 在notepad.exe中运行线程
	hThread = CreateRemoteThread(hProcess,NULL,0,pThreadProc,pRemoteBuf,0,NULL);

	WaitForSingleObject(hThread,INFINITE);
	CloseHandle(hThread);
	CloseHandle(hProcess);

	return TRUE;

}
// 提权函数
BOOL EnableDebugPriv() 
{ 
	HANDLE hToken;
	LUID sedebugnameValue; 
	TOKEN_PRIVILEGES tkp; 
 
	if ( ! OpenProcessToken( GetCurrentProcess(),TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY, &hToken ) ) 
	{
		printf("提权失败。");
		return FALSE; 
	}
 
	if ( ! LookupPrivilegeValue( NULL, SE_DEBUG_NAME, &sedebugnameValue ) ) 
	{ 
		CloseHandle( hToken ); 
		printf("提权失败。");
		return FALSE; 
	} 
	tkp.PrivilegeCount = 1; 
	tkp.Privileges[0].Luid = sedebugnameValue; 
	tkp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED; //
	if ( ! AdjustTokenPrivileges( hToken, FALSE, &tkp, sizeof tkp, NULL, NULL ) ) 
	{
		printf("提权失败。");
		CloseHandle( hToken );
	}
	else 
	{
		printf("提权成功！");
		return TRUE;
	}

}

int _tmain(int argc, _TCHAR* argv[])
{

	if (argc !=3)
	{
		_tprintf(L"USAGE: %s pid dll_path\n",argv[0]);
		return 1;
	}
	//inject dll

	EnableDebugPriv();

	if (InjectDll((DWORD)_tstol(argv[1]),argv[2]))
	
		_tprintf(L"InjectDll(\"%s\") success!!!\n",argv[2]);
	else
		_tprintf(L"InjectDll(\"%s\") failed!!!\n",argv[2]);
	
	

	/*DWORD dwPID = 0;
	LPCTSTR szDllPath = NULL;
	dwPID = 1776;
	szDllPath = L"C:\\work\\myhack.dll";
	InjectDll(dwPID,szDllPath);*/

	return 0;
}
