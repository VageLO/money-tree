CREATE TRIGGER Update_Balance_On_Transaction_Update
AFTER UPDATE ON Transactions
FOR EACH ROW

BEGIN

	UPDATE Accounts
	SET balance = CASE
		-- Update new.account if new.account_id EQUAL old.account_id
		WHEN new.account_id = old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'debit' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'credit' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) + new.amount, 2)
			
			WHEN new.transaction_type = 'credit' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) + new.amount, 2)
			
			WHEN new.transaction_type = 'transfer' AND (old.transaction_type = 'debit' OR old.transaction_type = 'transfer')
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) + old.amount) - new.amount, 2)
			
			WHEN new.transaction_type = 'transfer' AND old.transaction_type = 'credit'
			THEN ROUND(((SELECT balance FROM Accounts WHERE id = new.account_id) - old.amount) - new.amount, 2)

			ELSE RAISE(ABORT, "-- Update new.account if new.account_id EQUAL old.account_id")
			END
			
		WHEN new.account_id <> old.account_id
		THEN CASE
			WHEN new.transaction_type = 'debit' OR new.transaction_type = 'transfer'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount, 2)

			WHEN new.transaction_type = 'credit'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount, 2)

			ELSE RAISE(ABORT, "-- Update new.account if new.account_id NOT EQUAL old.account_id")
			END

	END WHERE id = new.account_id;
	
	-- Update old.account if new.account_id NOT EQUAL old.account_id
	UPDATE Accounts
	SET balance = CASE
		WHEN new.account_id <> old.account_id
		THEN CASE 
			WHEN old.transaction_type = 'debit' OR old.transaction_type = 'transfer'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) + old.amount, 2)

			WHEN old.transaction_type = 'credit'
			THEN ROUND((SELECT balance FROM Accounts WHERE id = old.account_id) - old.amount, 2)

			ELSE RAISE(ABORT, "Update old.account if new.account_id NOT EQUAL old.account_id")
			END
		ELSE RAISE(IGNORE)
	END WHERE id = old.account_id;
	
END;
