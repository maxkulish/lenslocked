package models

import (
	"testing"
	"time"
)

const env = "test"

func testingUserService() (*UserService, error) {
	us, err := NewUserService(env)
	if err != nil {
		return nil, err
	}
	us.DB.LogMode(false)

	us.FullReset()
	us.DB.AutoMigrate(&User{})
	return us, nil
}

func TestCreateUserService(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}

	user := User{
		Name:  "Test User",
		Email: "test@gmail.com",
	}

	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d", user.ID)
	}

	if time.Since(user.CreatedAt) > 10*time.Second {
		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
	}

	if time.Since(user.UpdatedAt) > 10*time.Second {
		t.Errorf("Expected UpdateAt to be recent. Recieved %s", user.UpdatedAt)
	}
}
