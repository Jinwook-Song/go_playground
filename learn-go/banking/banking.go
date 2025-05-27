package banking

import "fmt"

type bankAccount struct {
	owner   string
	balance int
}

var ErrorInsufficientBalance = fmt.Errorf("잔액이 부족합니다")

func NewBankAccount(owner string) *bankAccount {
	return &bankAccount{owner: owner, balance: 0}
}

func (b *bankAccount) Deposit(amount int) {
	b.balance += amount
}

func (b *bankAccount) Withdraw(amount int) error {
	if b.balance < amount {
		return ErrorInsufficientBalance
	}
	b.balance -= amount
	return nil
}

func (b *bankAccount) Balance() int {
	return b.balance
}

func (b *bankAccount) Owner() string {
	return b.owner
}

func (b *bankAccount) ChangeOwner(newOwner string) {
	b.owner = newOwner
}

func (b *bankAccount) String() string {
	return fmt.Sprintf("계좌 소유자: %s, 잔액: %d원", b.owner, b.balance)
}
