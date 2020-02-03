package service

import (
	"gophr.v2/gophr.api/user"
	"gophr.v2/gophr.api/user/mocks"
	"testing"
)

func TestService_GetByID(t *testing.T) {
	repo := new(mocks.Repository)

	want := &user.User{
		ID: "Judith Kuliiiiiiiiiiittttt 123",
		Username: "judith.kulit",
		Note: "Love ko yaaan!",
	}

	svc := New(repo)
}

func TestService_GetByEmail(t *testing.T) {}

func TestService_GetByUsername(t *testing.T) {}

func TestService_Save(t *testing.T) {}

