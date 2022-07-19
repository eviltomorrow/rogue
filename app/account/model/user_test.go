package model

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/rogue/lib/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var timeout = 10 * time.Second

func init() {
	mysql.DSN = "root:root@tcp(127.0.0.1:3306)/rogue_account?charset=utf8mb4&parseTime=true&loc=Local"
	mysql.Build()

}

func truncateUser() {
	var _sql = `truncate table user`
	_, err := mysql.DB.Exec(_sql)
	if err != nil {
		log.Fatal(err)
	}
}

var u1 = &User{
	UUID:     uuid.NewString(),
	NickName: sql.NullString{String: "shepard"},
	Email:    "eviltomorrow@163.com",
	Phone:    "12345678902",
}

var u2 = &User{
	UUID:     uuid.NewString(),
	NickName: sql.NullString{},
	Email:    "eviltomorrow@gamil.com",
	Phone:    "9658123547",
}

func TestUserWithInsertOne(t *testing.T) {
	_assert := assert.New(t)
	truncateUser()

	tx, err := mysql.DB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := UserWithInsertOne(tx, u1, timeout)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()

	user, err := UserWithSelectOneByUUID(mysql.DB, u1.UUID, timeout)
	_assert.Nil(err)

	_assert.Equal(u1.UUID, user.UUID)
	_assert.Equal(u1.NickName.String, user.NickName.String)
	_assert.Equal(u1.Email, user.Email)
	_assert.Equal(u1.Phone, user.Phone)
}

func TestUserWithSelectRange(t *testing.T) {
	_assert := assert.New(t)
	truncateUser()

	tx, err := mysql.DB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := UserWithInsertOne(tx, u1, timeout)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = UserWithInsertOne(tx, u2, timeout)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()

	users, err := UserWithSelectRange(mysql.DB, 0, 1, timeout)
	_assert.Nil(err)
	_assert.Equal(1, len(users))

	users, err = UserWithSelectRange(mysql.DB, 0, 2, timeout)
	_assert.Nil(err)
	_assert.Equal(2, len(users))

	users, err = UserWithSelectRange(mysql.DB, 0, 10, timeout)
	_assert.Nil(err)
	_assert.Equal(2, len(users))
}

func TestUserWithDeleteByUUID(t *testing.T) {
	_assert := assert.New(t)
	truncateUser()

	tx, err := mysql.DB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	affected, err := UserWithInsertOne(tx, u1, timeout)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	affected, err = UserWithDeleteByUUID(tx, u1.UUID, timeout)
	_assert.Nil(err)
	_assert.Equal(int64(1), affected)

	tx.Commit()
}
