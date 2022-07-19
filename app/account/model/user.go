package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/eviltomorrow/rogue/lib/mysql"
	jsoniter "github.com/json-iterator/go"
)

func UserWithInsertOne(exec mysql.Exec, user *User, timeout time.Duration) (int64, error) {
	if user == nil {
		return 0, nil
	}

	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `insert into user(id, uuid, nick_name, email, phone, disabled, create_timestamp) values (null, ?, ?, ?, ?, 0, now())`
	result, err := exec.ExecContext(ctx, _sql, user.UUID, user.NickName.String, user.Email, user.Phone)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UserWithSelectRange(exec mysql.Exec, offset, limit int64, timeout time.Duration) ([]*User, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `select id, uuid, nick_name, email, phone, disabled, create_timestamp, modify_timestamp from user where disabled = 0 limit ?, ?`
	rows, err := exec.QueryContext(ctx, _sql, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = make([]*User, 0, limit)
	for rows.Next() {
		var user = &User{}
		if err := rows.Scan(&user.ID, &user.UUID, &user.NickName, &user.Email, &user.Phone, &user.Disabled, &user.CreateTimestamp, &user.ModifyTimestamp); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func UserWithSelectOneByEmail(exec mysql.Exec, email string, timeout time.Duration) (*User, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `select id, uuid, nick_name, email, phone, disabled, create_timestamp, modify_timestamp from user where email = ?`
	row := exec.QueryRowContext(ctx, _sql, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user = &User{}
	if err := row.Scan(&user.ID, &user.UUID, &user.NickName, &user.Email, &user.Phone, &user.Disabled, &user.CreateTimestamp, &user.ModifyTimestamp); err != nil {
		return nil, err
	}

	return user, nil
}

func UserWithSelectOneByPhone(exec mysql.Exec, phone string, timeout time.Duration) (*User, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `select id, uuid, nick_name, email, phone, disabled, create_timestamp, modify_timestamp from user where phone = ?`
	row := exec.QueryRowContext(ctx, _sql, phone)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user = &User{}
	if err := row.Scan(&user.ID, &user.UUID, &user.NickName, &user.Email, &user.Phone, &user.Disabled, &user.CreateTimestamp, &user.ModifyTimestamp); err != nil {
		return nil, err
	}

	return user, nil
}

func UserWithSelectOneByUUID(exec mysql.Exec, uuid string, timeout time.Duration) (*User, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `select id, uuid, nick_name, email, phone, disabled, create_timestamp, modify_timestamp from user where uuid = ?`
	row := exec.QueryRowContext(ctx, _sql, uuid)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user = &User{}
	if err := row.Scan(&user.ID, &user.UUID, &user.NickName, &user.Email, &user.Phone, &user.Disabled, &user.CreateTimestamp, &user.ModifyTimestamp); err != nil {
		return nil, err
	}

	return user, nil
}

func UserWithDeleteByUUID(exec mysql.Exec, uuid string, timeout time.Duration) (int64, error) {
	ctx, cannel := context.WithTimeout(context.Background(), timeout)
	defer cannel()

	var _sql = `delete from user where uuid = ?`
	result, err := exec.ExecContext(ctx, _sql, uuid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const (
	FieldUserID              = "id"
	FieldUserUUID            = "uuid"
	FieldUserNickName        = "nick_name"
	FieldUserEmail           = "email"
	FieldUserPhone           = "phone"
	FieldDisabled            = "disabled"
	FieldUserCreateTimestamp = "create_timestamp"
	FieldUserModifyTimestamp = "modify_timestamp"
)

var UserFeilds = []string{
	FieldUserUUID,
	FieldUserNickName,
	FieldUserEmail,
	FieldUserPhone,
	FieldDisabled,
	FieldUserCreateTimestamp,
}

type User struct {
	ID              int64          `json:"id"`
	UUID            string         `json:"uuid"`
	NickName        sql.NullString `json:"nick_name"`
	Email           string         `json:"email"`
	Phone           string         `json:"phone"`
	Disabled        int8           `json:"disabled"`
	CreateTimestamp time.Time      `json:"create_timestamp"`
	ModifyTimestamp sql.NullTime   `json:"modify_timestamp"`
}

func (q *User) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(q)
	return string(buf)
}
