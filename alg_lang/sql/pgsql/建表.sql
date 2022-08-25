CREATE TABLE "bilibili"."video"
(
    "aid"    int8 NOT NULL,
    "cid"    int8 NOT NULL,
    "data" jsonb,
    "record" bool,
    "created_at" timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "video_pkey" PRIMARY KEY ("aid", "cid")
)
;

ALTER TABLE "bilibili"."video"
    OWNER TO "postgres";

CREATE INDEX "idx_video_created_at" ON "bilibili"."video" USING btree (
    "created_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
    );