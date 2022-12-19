INSERT INTO bilibili.USER ( "id", "name", "face" ) (
    SELECT DISTINCT ON( A.OWNER -> 'mid' )
        CAST ( A.OWNER -> 'mid' AS BIGINT ) ID,
        A.OWNER ->> 'name' NAME,
        A.OWNER ->> 'face' face
    FROM
        ( SELECT DATA -> 'owner' OWNER FROM bilibili.VIEW WHERE DATA ? 'owner' ) A
    WHERE
        A.OWNER ? 'mid'
);

INSERT INTO bilibili.view_v2 (
    "aid",
    "bvid",
    "uid",
    "title",
    "desc",
    "dynamic",
    "tid",
    "pic",
    "ctime",
    "tname",
    "videos",
    "pubdate",
    "record",
    "created_at",
    "updated_at",
    "deleted_at"
) (
    SELECT A
               .aid,
           A.bvid,
           CAST ( A.OWNER -> 'mid' AS BIGINT ) uid,
           A.title,
           A.desc,
           A.dynamic,
           CAST ( A.tid AS int ) tid,
           A.pic,
           to_timestamp(CAST( A.ctime AS float8)) ctime,
           A.tname,
           CAST (A.videos AS int) videos,
           to_timestamp(CAST( A.pubdate AS float8)) pubdate,
           CAST(A.cover_record AS int) record,
           A.created_at,
           A.created_at updated_at,
           A.deleted_at
    FROM
        (
            SELECT
                aid,
                bvid,
                DATA ->> 'tid' TID,
                DATA ->> 'title' title,
                DATA ->> 'dynamic' DYNAMIC,
                DATA ->> 'desc' DESC,
			DATA ->> 'pic' pic,
			DATA -> 'ctime' ctime,
			DATA -> 'tname' tname,
			DATA -> 'videos' videos,
			DATA -> 'pubdate' pubdate,
			DATA -> 'owner' OWNER,
			cover_record,
			created_at,
			deleted_at
            FROM
                bilibili.VIEW
        ) A
);

INSERT INTO bilibili.view_bak_v2 (
    "aid",
    "uid",
    "title",
    "desc",
    "dynamic",
    "tid",
    "pic",
    "ctime",
    "tname",
    "videos",
    "pubdate",
    "created_at",
    "updated_at"
) (
    SELECT DISTINCT ON(A.aid) A
        .aid,
        CAST ( A.OWNER -> 'mid' AS BIGINT ) uid,
        A.title,
        A.desc,
        A.dynamic,
        CAST ( A.tid AS int ) tid,
        A.pic,
        to_timestamp(CAST( A.ctime AS float8)) ctime,
        A.tname,
        CAST (A.videos AS int) videos,
        to_timestamp(CAST( A.pubdate AS float8)) pubdate,
        A.created_at,
        A.created_at updated_at
    FROM
        (
        SELECT
        aid,
        DATA ->> 'tid' TID,
        DATA ->> 'title' title,
        DATA ->> 'dynamic' DYNAMIC,
        DATA ->> 'desc' DESC,
        DATA ->> 'pic' pic,
        DATA -> 'ctime' ctime,
        DATA -> 'tname' tname,
        DATA -> 'videos' videos,
        DATA -> 'pubdate' pubdate,
        DATA -> 'owner' OWNER,
        created_at
        FROM
        bilibili.view_bak
        ) A
);
SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part, b.created_at,b.record
FROM bilibili.video  b
         LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p, deleted_at FROM bilibili.view)  a ON (a.p->'cid')::int8 = b.cid LIMIT 10;

INSERT INTO bilibili.video_v2 ( aid, CID, part, page, accept_format, video_codecid, duration, accept_quality, record, created_at, updated_at, deleted_at ) SELECT
                                                                                                                                                               b.aid,
                                                                                                                                                               b.CID,
                                                                                                                                                               A.P ->> 'part' part,
                                                                                                                                                               CAST ( A.P -> 'page' AS INT ) page,
                                                                                                                                                               b.DATA ->> 'accept_format' accept_format,
                                                                                                                                                               CAST ( b.DATA -> 'video_codecid' AS INT ) video_codecid,
                                                                                                                                                               CAST ( b.DATA -> 'timelength' AS INT ) duration,
                                                                                                                                                               b.DATA -> 'accept_quality' accept_quality,
                                                                                                                                                               b.record,
                                                                                                                                                               b.created_at,
                                                                                                                                                               b.created_at,
                                                                                                                                                               b.deleted_at
FROM
    bilibili.video b
        LEFT JOIN ( SELECT DATA -> 'title' title, DATA -> 'owner' OWNER, jsonb_path_query ( DATA, '$.pages[*]' ) P, deleted_at FROM bilibili.VIEW ) A ON ( A.P -> 'cid' ) :: INT8 = b.CID;

SELECT cid, ARRAY(SELECT jsonb_array_elements_text(accept_quality)) FROM  bilibili.video;
UPDATE bilibili.video SET accept_quality_v2 = d.array FROM bilibili.video v
LEFT JOIN ( SELECT CID, CAST(ARRAY ( SELECT jsonb_array_elements_text ( accept_quality ) ) AS int[])  FROM bilibili.video ) d ON v.CID = d.CID ;