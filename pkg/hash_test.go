package pkg

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	// Тестируем создание хеша
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword() error = %v", err)
	}

	// Проверяем, что хеш не пустой
	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	// Проверяем, что хеш отличается от оригинального пароля
	if hash == password {
		t.Error("HashPassword() returned the same value as input password")
	}

	// Тестируем проверку пароля
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash() failed for correct password")
	}

	// Тестируем проверку неправильного пароля
	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash() succeeded for wrong password")
	}
}
