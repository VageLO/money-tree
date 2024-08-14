package structs

const TransactionsWhereAccountId = `
SELECT Transactions.*, Accounts.title as 'from_account', NULL as 'to_account', Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE Transactions.to_account_id IS NULL AND Transactions.account_id = ?

UNION ALL

SELECT Transactions.*, from_account.title, to_account.title, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts as from_account ON from_account.id = Transactions.account_id
INNER JOIN Accounts as to_account ON to_account.id = Transactions.to_account_id
WHERE Transactions.to_account_id IS NOT NULL AND (Transactions.account_id = ? OR Transactions.to_account_id = ?);
`
const TransactionsWhereCategoryId = `
SELECT Transactions.*, Accounts.title as 'from_account', NULL as 'to_account', Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE Transactions.to_account_id IS NULL AND Transactions.category_id = ?

UNION ALL

SELECT Transactions.*, from_account.title, to_account.title, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts as from_account ON from_account.id = Transactions.account_id
INNER JOIN Accounts as to_account ON to_account.id = Transactions.to_account_id
WHERE Transactions.to_account_id IS NOT NULL AND Transactions.category_id = ?;
`

const Transactions = `
SELECT Transactions.*, Accounts.title as 'from_account', NULL as 'to_account', Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE Transactions.to_account_id IS NULL

UNION ALL

SELECT Transactions.*, from_account.title, to_account.title, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts as from_account ON from_account.id = Transactions.account_id
INNER JOIN Accounts as to_account ON to_account.id = Transactions.to_account_id
WHERE Transactions.to_account_id IS NOT NULL;
`
const StatisticsQuery = `
SELECT 
	SUM(CASE 
		WHEN transaction_type == 'debit'
		THEN amount
		
		WHEN transaction_type == 'transfer' AND Transactions.account_id == ?
		THEN amount
		
		ELSE 0 END) as debit,
	SUM (CASE 
		WHEN transaction_type == 'credit'
		THEN amount
		
		WHEN transaction_type == 'transfer' AND Transactions.account_id <> ?
		THEN CASE
			WHEN to_amount <> 0 AND to_amount IS NOT NULL
			THEN to_amount
			ELSE amount END
		ELSE 0 END) as credit, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE (Transactions.account_id = ? OR Transactions.to_account_id = ?) AND (date BETWEEN 'FIRST' AND 'LAST')
GROUP BY Categories.title
`
const InitTransactions = `
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
`

const InitAccounts = `
CREATE TABLE "Accounts" (
  "id"	INTEGER,
  "title"	TEXT NOT NULL UNIQUE,
  "currency"	TEXT NOT NULL,
  "balance"	NUMERIC NOT NULL DEFAULT 0,
  PRIMARY KEY("id" AUTOINCREMENT)
)
`

const InitCategories = `
CREATE TABLE "Categories" (
	"id"	INTEGER,
	"parent_id"	INTEGER,
	"title"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("id" AUTOINCREMENT)
)
`
const UpdateOnDelete = `
CREATE TRIGGER Update_Balance_On_Transaction_Delete
AFTER DELETE ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Accounts
	SET balance = CASE
		WHEN old.transaction_type = 'debit' OR old.transaction_type = 'transfer'
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount, 2)
		
		WHEN old.transaction_type = 'credit'
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount, 2)
		
		ELSE RAISE(ABORT, "ELSE UPDATE Accounts ON DELETE")
	END WHERE id = old.account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN old.to_account_id IS NOT NULL AND (old.to_amount IS NOT NULL AND old.to_amount <> 0)
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.to_account_id) - old.to_amount, 2)
		
		WHEN old.to_account_id IS NOT NULL AND old.amount IS NOT NULL
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.to_account_id) - old.amount, 2)
		
		ELSE RAISE(IGNORE)
	END WHERE id = old.to_account_id;
END
`

const UpdateOnInsert = `
CREATE TRIGGER Update_Balance_On_Transaction_Insert
AFTER INSERT ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Accounts
	SET balance = CASE
		WHEN new.transaction_type = 'debit' OR new.transaction_type = 'transfer'
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount, 2)
		
		WHEN new.transaction_type = 'credit'
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount, 2)
		
		ELSE RAISE(IGNORE)
	END WHERE id = new.account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN new.to_account_id IS NOT NULL AND (new.to_amount IS NOT NULL AND new.to_amount <> 0)
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.to_account_id) + new.to_amount, 2)
		
		WHEN new.to_account_id IS NOT NULL AND new.amount IS NOT NULL
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.to_account_id) + new.amount, 2)
		
		ELSE RAISE(IGNORE)
	END WHERE id = new.to_account_id;
END;
`

const UpdateOnUpdate = `
CREATE TRIGGER Update_Balance_On_Transaction_Update
AFTER UPDATE ON Transactions
FOR EACH ROW

BEGIN

	UPDATE Accounts
	SET balance = CASE
		-- Update new.account if new.account_id EQUAL old.account_id
		WHEN new.account_id = old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'debit' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'credit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) + new.amount, 2)
			
			WHEN new.transaction_type = 'credit' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) + new.amount, 2)
			
			WHEN new.transaction_type = 'transfer' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'transfer' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) - new.amount, 2)

			ELSE RAISE(ABORT, "-- Update new.account if new.account_id EQUAL old.account_id")
			END
			
		WHEN new.account_id <> old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' OR new.transaction_type = 'transfer'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount, 2)

			WHEN new.transaction_type = 'credit'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount, 2)

			ELSE RAISE(ABORT, "-- Update new.account if new.account_id NOT EQUAL old.account_id")
			END

	END WHERE id = new.account_id;
	
	-- Update old.account if new.account_id NOT EQUAL old.account_id
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN old.transaction_type = 'debit' OR old.transaction_type = 'transfer'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount, 2)

			WHEN old.transaction_type = 'credit'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount, 2)

			ELSE RAISE(ABORT, "Update old.account if new.account_id NOT EQUAL old.account_id")
			END
		ELSE RAISE(IGNORE)
	END WHERE id = old.account_id;
	
END;
`
const UpdateToAccount = `
CREATE TRIGGER Update_ToAccount_Balance
AFTER UPDATE ON Transactions
FOR EACH ROW

BEGIN
	
	-- Update new.to_account with NEW data
	UPDATE Accounts
	SET balance = CASE
		WHEN new.to_account_id IS NOT NULL AND (new.to_amount IS NOT NULL AND new.to_amount <> 0)
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.to_account_id) + new.to_amount, 2)
		
		WHEN new.to_account_id IS NOT NULL AND new.amount IS NOT NULL
		THEN ROUND((SELECT balance FROM Accounts WHERE id = new.to_account_id) + new.amount, 2)
		
		ELSE RAISE(ABORT, "Update new.to_account with NEW data")
	END WHERE id = new.to_account_id;
	
	-- Update old.to_account with OLD data
	UPDATE Accounts
	SET balance = CASE
		WHEN old.to_account_id IS NOT NULL AND (old.to_amount IS NOT NULL AND old.to_amount <> 0)
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.to_account_id) - old.to_amount, 2)
		
		WHEN old.to_account_id IS NOT NULL AND old.amount IS NOT NULL
		THEN ROUND((SELECT balance FROM Accounts WHERE id = old.to_account_id) - old.amount, 2)
		
		ELSE RAISE(ABORT, "Update old.to_account with OLD data")
	END WHERE id = old.to_account_id;
	
END
`
