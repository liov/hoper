#include <stdio.h>
#include <Python.h>

static PyObject* Py_Add_Formular(PyObject *self, PyObject *args)
{
    PyObject *pObj = NULL;
    int nThread = 0, nCalcData = -1, isDebug = 0, isCalcHis = 0;
    if (!PyArg_ParseTuple(args, "Oiiii", &pObj, &nThread, &nCalcData, &isDebug, &isCalcHis))
    {
        PyErr_Print();
        return NULL;
    }
    printf("%d%d%d%d\r\n", nThread, nCalcData, isDebug, isCalcHis);
    return Py_True;
}

/*
PyMethodDef结构体有四个字段。

第一个是字符串，表示在Python文件中对应的方法的名称;

第二个是对应的C代码的函数名称;

第三个是一个标志位，表示该Python方法是否需要参数，METH_NOARGS表示不需要参数，METH_VARARGS表示需要参数;

第四个是一个字符串，它是该方法的__doc__属性，这个不是必须的，可以为NULL。

PyMethodDef结构体数组最后以{NULL,NULL,0,NULL}结尾。
*/
static PyMethodDef CMODMethods[]=
{
    {"Add_Formular", Py_Add_Formular, METH_VARARGS, "Add formular."},
    {NULL, NULL, 0, NULL}
};

/*
这个函数是用于模块初始化的，即是在第一次使用import语句导入模块时会执行。

其函数名必须为initmodule_name这样的格式，在这里我们的模块名为CMOD，所以函数名就是initCMOD。

在这个函数中又调用了PyInitModule函数，它执行了模块的初始化操作。

Py_InitModule函数传入了两个参数，第一个参数为字符串，表示模块的名称;第二个参数是一个Py_MethodDef的结构体数组，表示该模块都具有哪些方法。

因此在 initCMOD 方法之前还需要先定义 CMODMethods 数组。
*/
PyMODINIT_FUNC initCMOD(void)
{
    Py_InitModule("CMOD", CMODMethods);
}