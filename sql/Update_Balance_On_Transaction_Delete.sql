CREATE TRIGGER Update_Balance_On_Transaction_Delete
AFTER DELETE ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Accounts
	SET balance = CASE
		WHEN old.transaction_type = 'debit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount)
		WHEN old.transaction_type = 'credit'
		THEN IIF(((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount) < 0,
			RAISE(ABORT, "ERROR: Not enough money"),
			(SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount)
		ELSE RAISE(ABORT, "ERROR: ELSE UPDATE Accounts ON Delete")
	END WHERE id = old.account_id;
END