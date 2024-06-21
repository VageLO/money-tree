CREATE TRIGGER Update_Balance_On_Transaction_Insert
AFTER INSERT ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Transactions
	SET balance = CASE
		WHEN new.transaction_type = 'debit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount)
		WHEN new.transaction_type = 'credit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount)
		ELSE RAISE(ABORT, "ERROR: ELSE UPDATE Transactions ON Insert")		
	END WHERE id = new.id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN new.transaction_type = 'debit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount)
		WHEN new.transaction_type = 'credit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount)
		ELSE RAISE(ABORT, "ERROR: ELSE UPDATE Accounts ON Insert")
	END WHERE id = new.account_id;
END;