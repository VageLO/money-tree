SELECT 
	SUM(CASE 
		WHEN transaction_type == 'debit'
		THEN amount
		
		WHEN transaction_type == 'transfer' AND Transactions.account_id == 1
		THEN amount
		
		ELSE 0 END) as debit,
	SUM (CASE 
		WHEN transaction_type == 'credit'
		THEN amount
		
		WHEN transaction_type == 'transfer' AND Transactions.account_id <> 1
		THEN CASE
			WHEN to_amount <> 0 AND to_amount IS NOT NULL
			THEN to_amount
			ELSE amount END
		ELSE 0 END) as credit, Categories.title FROM Transactions 
INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE (Transactions.account_id = 1 OR Transactions.to_account_id = 1) AND (date BETWEEN '2024-07-01' AND '2024-08-08')
GROUP BY Categories.title