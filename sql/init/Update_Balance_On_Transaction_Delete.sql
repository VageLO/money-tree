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