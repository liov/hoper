const fn = f => n => n ? f(n - 1) * n : 1;
const Y = y => (x => y(x(x)))(x => y(n => x(x)(n)));
console.log(Y(fn)(5));

const func = (f => n => n ? f(f)(n-1)*n : 1)(f => n => n ? f(f)(n-1)*n : 1);

const func1 = n => {
    return (f => n => n ? f(f)(n-1)*n : 1)(f => n => n ? f(f)(n-1)*n : 1);
}

//开始化简 func：
const func2 = n => {
    const fa = f => n => n ? f(f)(n-1)*n : 1;
    return fa(fa);
};

//提取重复形式 f => n => n ? f(f)(n-1)*n : 1：
const func3 = n => {
    const fa = f => {
        return n => n ? f(f)(n-1)*n : 1;
    };
    return fa(fa);
};

//可以看出，其主要递归逻辑来自 f(f), 我们将这一部分解耦：
const func4 = n => {
    const fa = f => {
        const fb = n => f(f)(n);
        return n => n ? fb(n-1)*n : 1;
    };
    return fa(fa);
};

//可以看到 返回值 不再需要 fc 接收的参数 f, 将返回值表达式具名, 以便提取出 fc, 分离逻辑：
const func5 = n => {
    const fa = f => {
        const fb = n => f(f)(n);
        const fc = n => n ? fb(n-1)*n : 1;
        return fc;
    };
    return fa(fa);
};

//fc 还在依赖 fb, 将 fb 作为参数传入 fc, 解除 fc 对 fb 的依赖：
const func6 = n => {
    const fa = f => {
        const fb = n => f(f)(n);
        const fc = fb => n => n ? fb(n-1)*n : 1;
        return fc(fb);
    };
    return fa(fa);
};

//可以发现 fc 是计算逻辑部分，将 fc 提取出 fa：
const func7 = n => {
    const fa = fc => f => {
        const fb = n => f(f)(n);
        return fc(fb);
    };
    const fc = fb => n => n ? fb(n - 1) * n : 1;
    return fa(fc)(fa(fc));
};

//构造一个函数 fd, 化简返回值的形式：
const func8 = n => {
    const fa = fc => f => {
        const fb = n => f(f)(n);
        return fc(fb);
    };
    const fc = fb => n => n ? fb(n - 1) * n : 1;
    const fd = fa => fc => {
        return fa(fc)(fa(fc));
    };
    return fd(fa)(fc);
};

//将 fa 带入 fd, 得到递归逻辑部分：
const func9 = n => {
    const fc = fb => n => n ? fb(n - 1) * n : 1;
    const fd = fc => {
        const fa = fc => f => {
            const fb = n => f(f)(n);
            return fc(fb);
        };
        return fa(fc)(fa(fc));
    };
    return fd(fc);
};

//化简fd
const func10 = n => {
    const fc = fb => n => n ? fb(n - 1) * n : 1;
    const fd = fc => {
        const fa = f => {
            const fb = n => f(f)(n);
            return fc(fb);
        };
        return fa(fa);
    };
    return fd(fc);
};

//化简fd
const func11 = n => {
    const fc = fb => n => n ? fb(n - 1) * n : 1;
    const fd = fc => (f => fc(n => f(f)(n)))(f => fc(n => f(f)(n)));
    return fd(fc);
};

//可以看到，两部分逻辑已经分离，可以得到 javascript 中的 Y 组合子：

/*λ 表达式的等价形式*/
Y1 = y => (x => y(x(x)))(x => y(x(x)));

/*推导出的形式*/
Y2 = y => (x => y(n => x(x)(n)))(x => y(n => x(x)(n)));
