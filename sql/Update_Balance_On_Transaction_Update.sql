CREATE TRIGGER Update_Balance_On_Transaction_Update
AFTER UPDATE ON Transactions
FOR EACH ROW

BEGIN
	-- TODO: Fix where transaction amount grater then balance
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id = old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount) - new.amount
			
			WHEN new.transaction_type = 'debit' AND old.transaction_type = 'credit'
			THEN ((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount) - new.amount
			
			WHEN new.transaction_type = 'credit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) + new.amount
			
			WHEN new.transaction_type = 'credit' AND old.transaction_type = 'credit'
			THEN ((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) + new.amount

			ELSE RAISE(ABORT, "ELSE UPDATE Accounts new=old")
			END
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN new.transaction_type = 'debit' OR new.transaction_type = 'transfer'
			THEN (SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount

			WHEN new.transaction_type = 'credit'
			THEN (SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount

			ELSE RAISE(ABORT, "ELSE UPDATE Accounts new.account_id")
			END

	END WHERE id = new.account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN old.transaction_type = 'debit' OR old.transaction_type = 'transfer'
			THEN (SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount

			WHEN old.transaction_type = 'credit'
			THEN (SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount

			ELSE RAISE(ABORT, "ELSE UPDATE Accounts old.account_id")
			END
		ELSE RAISE(IGNORE)
	END WHERE id = old.account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN new.to_account_id IS NOT NULL AND new.to_amount IS NOT NULL
		THEN (SELECT balance FROM Accounts WHERE id = new.to_account_id) + new.to_amount
		ELSE RAISE(IGNORE)
	END WHERE id = new.to_account_id;
	
	UPDATE Accounts
	SET balance = CASE
		WHEN old.to_account_id IS NOT NULL AND old.to_amount IS NOT NULL
		THEN (SELECT balance FROM Accounts WHERE id = old.to_account_id) - old.to_amount
		ELSE RAISE(ABORT, "OLD")
	END WHERE id = old.to_account_id;
	
END;
