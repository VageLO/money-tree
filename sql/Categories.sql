CREATE TABLE "Categories" (
	"id"	INTEGER,
	"parent_id"	INTEGER,
	"title"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("id" AUTOINCREMENT)
)