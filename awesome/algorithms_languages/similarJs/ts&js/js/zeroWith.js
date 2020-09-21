/**
 * 零宽空格（zero-width space, ZWSP）用于可能需要换行处。
 Unicode: U+200B  HTML: &#8203;
 零宽不连字 (zero-width non-joiner，ZWNJ)放在电子文本的两个字符之间，抑制本来会发生的连字，而是以这两个字符原本的字形来绘制。
 Unicode: U+200C  HTML: &#8204;
 零宽连字（zero-width joiner，ZWJ）是一个控制字符，放在某些需要复杂排版语言（如阿拉伯语、印地语）的两个字符之间，使得这两个本不会发生连字的字符产生了连字效果。
 Unicode: U+200D  HTML: &#8205;
 左至右符号（Left-to-right mark，LRM）是一种控制字符，用于计算机的双向文稿排版中。
 Unicode: U+200E  HTML: &lrm; &#x200E; 或&#8206;
 右至左符号（Right-to-left mark，RLM）是一种控制字符，用于计算机的双向文稿排版中。
 Unicode: U+200F  HTML: &rlm; &#x200F; 或&#8207;
 字节顺序标记（byte-order mark，BOM）常被用来当做标示文件是以UTF-8、UTF-16或UTF-32编码的标记。
 Unicode: U+FEFF
 * @param str
 * @returns {string}
 */
// str -> 零宽字符
function strToZeroWidth(str) {
    return str
        .split('')
        .map(char => char.charCodeAt(0).toString(2)) // 1 0 空格
        .join(' ')
        .split('')
        .map(binaryNum => {
            if (binaryNum === '1') {
                return '​'; // &#8203;
            } else if (binaryNum === '0') {
                return '‌'; // &#8204;
            } else {
                return '‍'; // &#8205;
            }
        })
        .join('‎') // &#8206;
}

// 零宽字符 -> str
function zeroWidthToStr(zeroWidthStr) {
    return zeroWidthStr
        .split('‎') // &#8206;
        .map(char => {
            if (char === '​') { // &#8203;
                return '1';
            } else if (char === '‌') { // &#8204;
                return '0';
            } else { // &#8205;
                return ' ';
            }
        })
        .join('')
        .split(' ')
        .map(binaryNum => String.fromCharCode(parseInt(binaryNum, 2)))
        .join('')
}

const rep = { // 替换用的数据，使用了4个零宽字符代理二进制
    '00': '\u200b',
    '01': '\u200c',
    '10': '\u200d',
    '11': '\uFEFF'
};

function hide(str) {
    str = str.replace(/[^\x00-\xff]/g, function (a) { // 转码 Latin-1 编码以外的字符。
        return escape(a).replace('%', '\\');
    });

    str = str.replace(/[\s\S]/g, function (a) { // 处理二进制数据并且进行数据替换
        a = a.charCodeAt().toString(2);
        a = a.length < 8 ? Array(9 - a.length).join('0') + a : a;
        return a.replace(/../g, function (a) {
            return rep[a];
        });
    });
    return str;
}

const tpl = '("@code".replace(/.{4}/g,function(a){const rep={"\u200b":"00","\u200c":"01","\u200d":"10","\uFEFF":"11"};return String.fromCharCode(parseInt(a.replace(/./g, a=>rep[a]),2))}))';

function hider(code, type) {
    let str = hide(code); // 生成零宽字符串

    str = tpl.replace('@code', str); // 生成模版
    if (type === 'eval') {
        str = 'eval' + str;
    } else {
        str = 'Function' + str + '()';
    }

    return str;
}

let a = strToZeroWidth("我爱你")
console.log(a)
console.log(rep['00'])
console.log(zeroWidthToStr(a))
console.log(hider('alert("我爱你")'))

const censored = '敏感词';
// 使用零宽度空格符对字符串进行分隔
console.log(Array.from(censored).join('\u200c'))// '敏​感​词'敏‌感‌词
