package main

import (
	"time"
)

func TxCancelation(c *CustomContext) {
	for {
		// Cancellation takes place every N minutes,
		// where N is set in configuration file, with cancel_time parameter
		time.Sleep(time.Duration(c.Config.CancellationInterval) * time.Minute)

		// Lock mutex, to secure DB from parallel access.
		c.DbMut.Lock()

		var user User
		c.Db.First(&user)

		// Double check that user exists.
		if !(user.ID > 0) {
			c.Lg.Fatal("User account not exist.")
		}

		// Take last 10 odd transaction:
		// take 1..3.. ... ..17..19
		// skip 2..4.. .. ..16..18
		// where transaction 1 is last transaction, 2 is second last and so on
		var transactions []Transaction
		c.Db.Order("id DESC").Limit(19).Find(&transactions)

		var transactionCanceled int
		for i := 0; i < len(transactions); i += 2 {
			transaction := transactions[i]
			if transaction.Canceled {
				continue
			}

			if err := user.CancelTransaction(&transaction); err != nil {
				c.Lg.Error(err)
				break
			}

			c.Db.Save(&transaction)
			c.Db.Save(&user)

			transactionCanceled++
		}

		c.Lg.Printf("Transactions canceled: %d", transactionCanceled)

		c.DbMut.Unlock()
	}
}
