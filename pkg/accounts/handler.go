package accounts

import db "github.com/Ritwiksrivastava0809/go-bank/pkg/db/sqlc"

func Reverse(accounts []db.Account) {
	for i, j := 0, len(accounts)-1; i < j; i, j = i+1, j-1 {
		accounts[i], accounts[j] = accounts[j], accounts[i]
	}
}
