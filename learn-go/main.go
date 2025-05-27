package main

import (
	"fmt"
	"log"

	"github.com/jinwook-song/learn-go/banking"
)

func main() {
	account := banking.NewBankAccount("jinwook")
	fmt.Printf("최초 계좌 소유자: %s\n", account.Owner())

	account.Deposit(1000)
	fmt.Printf("현재 잔액: %d원\n", account.Balance())

	// 계좌 소유자 변경
	account.ChangeOwner("song")
	fmt.Printf("변경된 계좌 소유자: %s\n", account.Owner())

	// 첫 번째 출금 시도
	if err := account.Withdraw(500); err != nil {
		log.Printf("출금 실패: %v\n", err)
	} else {
		fmt.Printf("출금 성공! 현재 잔액: %d원\n", account.Balance())
	}

	// 두 번째 출금 시도
	if err := account.Withdraw(500); err != nil {
		log.Printf("출금 실패: %v\n", err)
	} else {
		fmt.Printf("출금 성공! 현재 잔액: %d원\n", account.Balance())
	}

	// 잔액 초과 출금 시도
	if err := account.Withdraw(500); err != nil {
		log.Printf("출금 실패: %v\n", err)
	} else {
		fmt.Printf("출금 성공! 현재 잔액: %d원\n", account.Balance())
	}

	fmt.Printf("최종 잔액 (소유자: %s): %d원\n", account.Owner(), account.Balance())

	fmt.Println(account)
}
