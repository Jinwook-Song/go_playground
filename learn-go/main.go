package main

import (
	"fmt"
	"log"

	"github.com/jinwook-song/learn-go/banking"
	"github.com/jinwook-song/learn-go/my_dict"
)

func main() {
	// 은행 계좌 예제
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

	// 사전 예제
	fmt.Println("\n[사전 예제]")
	dictionary := my_dict.NewDictionary()

	// 단어 추가
	err := dictionary.Add("hello", "안녕하세요")
	if err != nil {
		fmt.Println("단어 추가 실패:", err)
	}

	// 단어 검색
	def, err := dictionary.Search("hello")
	if err != nil {
		fmt.Println("단어 검색 실패:", err)
	} else {
		fmt.Printf("hello의 뜻: %s\n", def)
	}

	// 존재하지 않는 단어 검색
	_, err = dictionary.Search("없는단어")
	if err != nil {
		fmt.Println("없는 단어 검색 결과:", err)
	}

	// 단어 수정
	err = dictionary.Update("hello", "안녕!")
	if err != nil {
		fmt.Println("단어 수정 실패:", err)
	}

	// 수정된 단어 검색
	def, _ = dictionary.Search("hello")
	fmt.Printf("수정된 hello의 뜻: %s\n", def)

	// 단어 삭제
	err = dictionary.Delete("hello")
	if err != nil {
		fmt.Println("단어 삭제 실패:", err)
	}

	// 삭제된 단어 검색
	_, err = dictionary.Search("hello")
	if err != nil {
		fmt.Println("삭제된 단어 검색 결과:", err)
	}
}
