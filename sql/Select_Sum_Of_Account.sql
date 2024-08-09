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