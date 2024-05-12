UPDATE `draft` set info = JSON_SET(info,"$.expressRule",JSON_ARRAY("1"),"$.serviceRule[0]","包邮") where id = 1;

UPDATE `draft` set info = REPLACE(info,'"expressRule":["1"]','"expressRule":["1"]')