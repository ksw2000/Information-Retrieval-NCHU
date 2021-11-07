CREATE TABLE "inverted_table" (
	"term"	TEXT,
	"aid"	INTEGER
);

CREATE TABLE "data" (
	"index_id"	INTEGER,
	"aid"	INTEGER,
	"articles"	TEXT,
	"title"	TEXT,
	PRIMARY KEY("index_id")
);