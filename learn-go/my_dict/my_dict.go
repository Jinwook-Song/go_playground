package my_dict

import "fmt"

type dictionary map[string]string

// 에러 정의
var (
	ErrNotFound   = fmt.Errorf("단어를 찾을 수 없습니다")
	ErrWordExists = fmt.Errorf("해당 단어가 이미 존재합니다")
	ErrCantUpdate = fmt.Errorf("업데이트할 단어가 사전에 없습니다")
	ErrCantDelete = fmt.Errorf("삭제할 단어가 사전에 없습니다")
)

// new dictionary
func NewDictionary() dictionary {
	return make(dictionary)
}

// Search 단어 검색
func (d dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if !exists {
		return "", ErrNotFound
	}
	return value, nil
}

// Add 단어 추가
func (d dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case ErrNotFound:
		d[word] = def
	case nil:
		return ErrWordExists
	}
	return nil
}

// Update 단어 수정
func (d dictionary) Update(word, newDef string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = newDef
	case ErrNotFound:
		return ErrCantUpdate
	}
	return nil
}

// Delete 단어 삭제
func (d dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		delete(d, word)
	case ErrNotFound:
		return ErrCantDelete
	}
	return nil
}
