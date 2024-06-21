CREATE TRIGGER Update_Balance_On_Transaction_Update
BEFORE UPDATE ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id = old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' AND old.transaction_type = 'debit'
			THEN IIF((((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount) - new.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"),
				((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount) - new.amount)
			
			WHEN new.transaction_type = 'debit' AND old.transaction_type = 'credit'
			THEN IIF((((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount) - new.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"),
				((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount) - new.amount)
			
			WHEN new.transaction_type = 'credit' AND old.transaction_type = 'debit'
			THEN IIF((((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) + new.amount) < 0, 
				RAISE(ABORT, "ERROR: Not enough money"), 
				((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) + new.amount)
			
			WHEN new.transaction_type = 'credit' AND old.transaction_type = 'credit'
			THEN IIF((((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) + new.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"), 
				((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) + new.amount)
			
			ELSE RAISE(ABORT, "ERROR: ELSE UPDATE Accounts new=old")
			END
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN new.transaction_type = 'debit'
			THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"), 
				(SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount)

			WHEN new.transaction_type = 'credit'
			THEN IIF(((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"), 
				(SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount)

			ELSE RAISE(ABORT, "Error: Update Account new.account_id")
			END

	END WHERE id = new.account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN old.transaction_type = 'debit'
			THEN IIF(((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"), 
				(SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount)

			WHEN old.transaction_type = 'credit'
			THEN IIF(((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount) < 0,
				RAISE(ABORT, "ERROR: Not enough money"), 
				(SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount)

			ELSE RAISE(ABORT, "Error: Update Account old.account_id")
			END
	END WHERE id = old.account_id;
END;