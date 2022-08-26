UPDATE video SET data = data - '{seek_type,seek_param}'::text[];

SELECT
    aid,p->'cid' cid ,p->'page' page,p->'part' part
FROM
    "bilibili"."view",jsonb_path_query(data,'$.pages[*]') AS p
WHERE aid = 33640070

SELECT
    a.aid,b.cid,,a.data->'title' title,
    p->'page' page,p->'part' part
FROM
    "bilibili"."view" a,jsonb_path_query(a.data,'$.pages[*]') AS p
LEFT JOIN "bilibili"."video" b ON (p->'cid')::int8 = b.cid
WHERE b.record = false
LIMIT 10;

SELECT t.cid FROM "bilibili"."view", jsonb_to_recordset(jsonb_path_query(data,'$.pages[*]') - '{vid,from,weblink,duration,dimension,first_frame}'::text[]) AS t(cid int8,page int2,part text) LIMIT 10;

SELECT
    aid,t.cid,t.page,t.part
FROM
    "bilibili"."view",jsonb_path_query(data,'$.pages[*]') AS p,jsonb_to_record(p) AS t(cid int8, page int2, part text)
WHERE aid = 33640070;

SELECT * FROM "view" WHERE (data-> 'cid')::int8 = 760279439;

SELECT data #> '{accept_quality,0}' quality FROM "video" LIMIT 10;


SELECT b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM "bilibili"."video" b
         LEFT JOIN (SELECT data->'title' title ,jsonb_path_query(data,'$.pages[*]') p FROM "bilibili"."view")  a ON (a.p->'cid')::int8 = b.cid
WHERE b.record = false AND b.aid < 10000000000  ORDER BY b.aid DESC
    LIMIT 20;