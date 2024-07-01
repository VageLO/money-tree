CREATE TABLE "Transactions" (
	"id"	INTEGER,
	"account_id"	INTEGER NOT NULL,
	"to_account_id"	INTEGER,
	"category_id"	INTEGER NOT NULL,
	"transaction_type"	TEXT NOT NULL,
	"date"	TEXT NOT NULL,
	"amount"	NUMERIC NOT NULL,
	"to_amount"	NUMERIC,
	"description"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("category_id") REFERENCES "Categories"("id") ON DELETE CASCADE,
	FOREIGN KEY("account_id") REFERENCES "Accounts"("id") ON DELETE CASCADE,
	FOREIGN KEY("to_account_id") REFERENCES "Accounts"("id") ON DELETE CASCADE
)