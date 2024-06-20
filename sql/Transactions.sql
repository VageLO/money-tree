CREATE TABLE "Transactions" (
	"id"	INTEGER,
	"account_id"	INTEGER NOT NULL,
	"category_id"	INTEGER NOT NULL,
	"transaction_type"	TEXT NOT NULL,
	"date"	TEXT NOT NULL,
	"amount"	NUMERIC NOT NULL,
	"balance"	NUMERIC NOT NULL,
	FOREIGN KEY("category_id") REFERENCES "Categories"("id") ON DELETE CASCADE,
	FOREIGN KEY("account_id") REFERENCES "Accounts"("id") ON DELETE CASCADE,
	PRIMARY KEY("id" AUTOINCREMENT)
)