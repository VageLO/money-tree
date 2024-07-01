SELECT Transactions.*, Accounts.title as 'from_account', NULL as 'to_account', Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE Transactions.to_account_id IS NULL

UNION ALL

SELECT Transactions.*, from_account.title, to_account.title, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts as from_account ON from_account.id = Transactions.account_id
INNER JOIN Accounts as to_account ON to_account.id = Transactions.to_account_id
WHERE Transactions.to_account_id IS NOT NULL;
