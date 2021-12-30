 ="insert into table_test(code, init_date) values ("&A2&", '"&B2&"');"
 ="INSERT INTO student(id,name,age) VALUES("&A2&",'"&B2&"','"&C2&"')"
 ="INSERT INTO `d_aura_jike`.`enum_data_value_info`(`field_id`, `field_value_desc`, `field_value`) VALUES (41, '"&A2&"',"&ROW()-1&");"
 ="UPDATE product_info SET sort_number = "&E2&" WHERE id = "&B2&";"
 
# 常量字符串过长，wps无解
excel规范中说：
公式内容的长度 1,024 个字符
如果你的公式没有超过1024个字符却得到公式太长的提示，通常是因为公式中遗漏或多输入括号、逗号等。
如果公式确实超过1024字符。可以用定义名称的方法将公式字符数减少。
例如：插入》名称》定义 x =offset(sheet1!$a$1,,,counta(sheet1!$a:$a),counta(sheet1!$1:$1))
定义好名称后，就可以将长公式中的offset()部分用x代替，大大减少了公式的长度。例如原公式是
=index(offset(....),row(),column())
定义好名称后就可以改为
=index(x,row(),column())

# CONCAT
在单元格中输入公式=CONCAT(B2:B7&"、")，并同时按下Ctrl+Shift+Enter三键结束即可完成。

# 时间格式化
=TEXT(P2,"YYYY-MM-DD HH:MM:SS")