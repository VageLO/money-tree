CREATE TRIGGER Update_Balance_On_Transaction_Insert
AFTER INSERT ON Transactions
FOR EACH ROW

BEGIN
	UPDATE Transactions
	SET balance = (SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount WHERE id = new.id AND new.transaction_type = 'debit';
	UPDATE Transactions
	SET balance = (SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount WHERE id = new.id AND new.transaction_type = 'credit';
	UPDATE Accounts
	SET balance = (SELECT balance FROM Accounts WHERE id = new.account_id) - new.amount WHERE id = new.account_id AND new.transaction_type = 'debit';
	UPDATE Accounts
	SET balance = (SELECT balance FROM Accounts WHERE id = new.account_id) + new.amount WHERE id = new.account_id AND new.transaction_type = 'credit';
END;