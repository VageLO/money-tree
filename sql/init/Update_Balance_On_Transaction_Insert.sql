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
