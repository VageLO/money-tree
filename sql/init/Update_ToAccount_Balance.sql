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