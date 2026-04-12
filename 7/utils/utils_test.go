package utils_test

import (
	"sync"
	"testing"
	"time"

	"practice-7/utils"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := utils.HashPassword("mypassword")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if !utils.CheckPassword(hash, "mypassword") {
		t.Error("CheckPassword: expected true for correct password")
	}
	if utils.CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword: expected false for wrong password")
	}
}

func TestGenerateJWT(t *testing.T) {
	utils.SetJWTSecret("test-secret-key")

	id := uuid.New()
	token, err := utils.GenerateJWT(id, "user")
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}
	if token == "" {
		t.Error("GenerateJWT returned empty token")
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	rl := utils.NewRateLimiter(3, 10*time.Second)
	key := "test-user"

	for i := 1; i <= 3; i++ {
		if !rl.Allow(key) {
			t.Errorf("Allow returned false on request %d (limit is 3)", i)
		}
	}
	if rl.Allow(key) {
		t.Error("Allow should return false after limit is exceeded")
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	rl := utils.NewRateLimiter(2, 50*time.Millisecond)
	key := "reset-user"

	rl.Allow(key)
	rl.Allow(key)
	if rl.Allow(key) {
		t.Error("third request should be denied")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow(key) {
		t.Error("Allow should succeed after window reset")
	}
}

func TestRateLimiter_Concurrency(t *testing.T) {
	rl := utils.NewRateLimiter(200, time.Minute)
	key := "concurrent-key"

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Allow(key)
		}()
	}
	wg.Wait()
}

func TestRateLimiter_MultipleKeys(t *testing.T) {
	rl := utils.NewRateLimiter(1, time.Minute)

	if !rl.Allow("key-a") {
		t.Error("key-a first request should be allowed")
	}
	if rl.Allow("key-a") {
		t.Error("key-a second request should be denied")
	}
	if !rl.Allow("key-b") {
		t.Error("key-b first request should be allowed")
	}
}
