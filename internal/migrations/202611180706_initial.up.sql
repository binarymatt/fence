CREATE TABLE IF NOT EXISTS "tenants" (
  "id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  "updated_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
  PRIMARY KEY ("id")
)

--bun:split

PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS "entities" (
  "id" VARCHAR NOT NULL, 
  "type" VARCHAR NOT NULL, 
  "tenant_id" VARCHAR NOT NULL, 
  "parents" json, 
  "attributes" json, 
  "tags" json, 
  "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  "updated_at" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  PRIMARY KEY ("id", "type"),
  FOREIGN KEY ("tenant_id") REFERENCES tenants ("id")
)
--bun:split
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS "policies" (
  "id" VARCHAR NOT NULL,
  "tenant_id" VARCHAR ,
  "content" VARCHAR,
  "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp, 
  "updated_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("tenant_id") REFERENCES tenants ("id")
)
