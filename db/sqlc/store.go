package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore 提供了执行SQL数据库查询和事务的所有功能
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx 数据库事务中执行Tx
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("TX 错误:%v, 回滚错误: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams 包含传输事务的输入参数
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult 转账交易结果
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx from toAccount to fromAccount
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// 转出
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 转进
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			//result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	ID:     arg.FromAccountID,
			//	Amount: -arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
			//
			//result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	ID:     arg.ToAccountID,
			//	Amount: arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
			//result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	ID:     arg.ToAccountID,
			//	Amount: arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
			//
			//result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			//	ID:     arg.FromAccountID,
			//	Amount: -arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}

//var txKey = struct{}{}

// TransFerTx 执行从一个account到另一个account的资金转移
// 它在单个数据库事务中创建转账记录、添加account条目和更新account余额

// 只能From Account1 to Account2 传帐
//func (store *Store) TransFerTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
//	var result TransferTxResult
//
//	err := store.execTx(ctx, func(q *Queries) error {
//		var err error
//
//		//txName := ctx.Value(txKey)
//		//fmt.Println(txName, "create transfer")
//		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
//			FromAccountID: arg.FromAccountID,
//			ToAccountID:   arg.ToAccountID,
//			Amount:        arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		// 转出
//		//fmt.Println(txName, "create entry1")
//		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
//			AccountID: arg.FromAccountID,
//			Amount:    -arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		// 转进
//		//fmt.Println(txName, "create entry2")
//		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
//			AccountID: arg.ToAccountID,
//			Amount:    arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		// update account
//		// get account -> update its balance
//		//fmt.Println(txName, "get account 1")
//		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
//		//if err != nil {
//		//	return err
//		//}
//		//fmt.Println(txName, "update account 1")
//		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//			ID:     arg.FromAccountID,
//			Amount: -arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//		//fmt.Println(txName, "get account 2")
//		//account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
//		//if err != nil {
//		//	return err
//		//}
//		//fmt.Println(txName, "update account 2")
//		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
//			ID:     arg.ToAccountID,
//			Amount: arg.Amount,
//		})
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	return result, err
//}
