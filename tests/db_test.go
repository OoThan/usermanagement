package tests

import (
	"testing"

	"github.com/OoThan/usermanagement/internal/model"
	"github.com/OoThan/usermanagement/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func newDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	db.AutoMigrate(&model.User{})
	return db, err
}

func TestUserCreate(t *testing.T) {
	db, err := newDB()
	if err != nil {
		logger.Sugar.Error(err)
	}
	defer db.Close()

	t.Run("Create", func(t *testing.T) {
		user := &model.User{
			Name:  "Mya Mya",
			Email: "myamya@email.com",
		}

		err := db.Create(user).Error
		assert.NoError(t, err)

		assert.NotEqual(t, uint64(0), user.Id, "Expected user ID to be set, but got 0")
	})
}

func TestUserUpdate(t *testing.T) {
	db, err := newDB()
	if err != nil {
		logger.Sugar.Error(err)
	}
	defer db.Close()

	t.Run("Update", func(t *testing.T) {
		user := &model.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		err := db.Create(user).Error
		assert.NoError(t, err)

		user.Name = "Updated Name"
		err = db.Save(user).Error
		assert.NoError(t, err)

		updatedUser := &model.User{}
		err = db.First(updatedUser, user.Id).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Name, updatedUser.Name, "Expected user name to be updated")
	})
}

func TestGetUserByID(t *testing.T) {
	db, err := newDB()
	if err != nil {
		logger.Sugar.Error(err)
	}
	defer db.Close()

	t.Run("GetUserByID", func(t *testing.T) {
		expectedUser := &model.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		err := db.Create(expectedUser).Error
		assert.NoError(t, err)
		assert.NotNil(t, expectedUser, "Expected user to be found, but got nil")

		user := &model.User{}
		err = db.First(user, expectedUser.Id).Error
		assert.NoError(t, err)

		assert.Equal(t, expectedUser.Name, user.Name, "Expected user names to match")
		assert.Equal(t, expectedUser.Email, user.Email, "Expected user emails to match")
	})
}

func TestUserDelete(t *testing.T) {
	db, err := newDB()
	if err != nil {
		logger.Sugar.Error(err)
	}
	defer db.Close()

	t.Run("Delete", func(t *testing.T) {
		user := &model.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		err := db.Create(user).Error
		assert.NoError(t, err)

		err = db.Delete(&model.User{}, user.Id).Error
		assert.NoError(t, err)

		deletedUser := &model.User{}
		err = db.First(deletedUser, user.Id).Error
		assert.Error(t, err, "Expected user to be deleted, but it was found")
		// assert.Nil(t, deletedUser, "Expected user to be nil, but got a value")
	})
}
