package models

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID           int64  `json:"id"`
	User         string `json:"user"`
	Pass         string `json:"pass"`
	Hash         string `json:"-"`
	Status       string `json:"status"`
	Otp          string `json:"otp"`
	OtpExp       string `json:"otp_exp"`
	Logged       int    `json:"logged"`
	Token        string `json:"token"`
	TokenExp     string `json:"token_exp"`
	Created      string `json:"created_dt"`
	Modified     string `json:"modified_dt"`
	ExpiredOtp   int    `json:"-"`
	ExpiredToken int    `json:"-"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	// just a post-process after a decode..
	u.Hash = fmt.Sprintf("%x", md5.Sum([]byte(u.Pass)))
	return nil
}

func (u *User) SanityCheck(data *User, which string) bool {

	switch which {
	case "ADD", "UPDATE", "LOG":
		if data.User == "" || data.Pass == "" {
			return false
		}
	case "ADD-LEN":
		//at least 4 chars?
		if len(data.User) < 4 || len(data.Pass) < 4 {
			return false
		}
	case "OTP":
		if data.User == "" || data.Otp == "" {
			return false
		}
	case "DELETE":
		if data.User == "" {
			return false
		}

	}
	return true
}

func (u *User) Get(ctx context.Context, db *sql.DB, who string) (*User, error) {
	r := `SELECT 
			ifnull(id,''), 
			ifnull(user,''), 
			ifnull(pass,''), 
			ifnull(status,''), 
			ifnull(otp,''), 
			ifnull(logged,0), 
			ifnull(token,''), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,''),
			ifnull((otp_exp   <now()),0),
			ifnull((token_exp <now()),0)
		FROM  users WHERE user = ?`
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("Get", err)
		return nil, err
	}
	defer stmt.Close()
	var data User
	err = stmt.QueryRowContext(ctx, who).Scan(
		&data.ID,
		&data.User,
		&data.Pass,
		&data.Status,
		&data.Otp,
		&data.Logged,
		&data.Token,
		&data.Created,
		&data.Modified,
		&data.ExpiredOtp,
		&data.ExpiredToken,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Get NOT_FOUND", who)
		return nil, errors.New("Info not found")
	case err != nil:
		log.Println("Get", err)
		return nil, err
	}
	//sounds good ;-)
	return &data, nil
}

func (u *User) Exists(ctx context.Context, db *sql.DB, who string) int {
	r := `SELECT count(id)
                FROM  users WHERE user = ?`

	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -1
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRowContext(ctx, who).Scan(&id)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -2
	}
	//sounds good ;-)
	return id
}

func (u *User) Create(ctx context.Context, db *sql.DB, data *User) int64 {
	//fmt
	r := `INSERT INTO users (
                user,
                pass,
                otp,
                otp_exp,
                status,
                created_dt)
              VALUES (?, ?, ?, ?, 'pending',Now())
              ON DUPLICATE KEY UPDATE
                 modified_dt = Now() `
	//exec
	result, err := db.ExecContext(ctx, r,
		data.User,
		data.Pass,
		data.Otp,
		data.OtpExp,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return -1
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("SQL_ERR", err)
		return -2
	}
	//sounds good ;-)
	data.ID = int64(id)
	return data.ID

}

func (u *User) Update(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                pass        = ?,
                modified_dt = Now()
              WHERE  user   = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Pass,
		data.User,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *User) Delete(ctx context.Context, db *sql.DB, who string) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                status      = 'deleted',
                token_exp   = date_add(now(), interval -1 day),
                modified_dt = Now()
              WHERE  user = ?`
	//exec
	result, err := db.ExecContext(ctx, r, who)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *User) SetLogged(ctx context.Context, db *sql.DB, who string, flag int) bool {

	r := `UPDATE users SET logged=?, modified_dt=Now() WHERE  user = ?`

	result, err := db.ExecContext(ctx, r, flag, who)
	if err != nil {
		log.Println("SetLogged", who, err)
		return false
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SetLogged", who, err)
		return false
	}
	//sounds good ;-)
	return true
}

func (u *User) UpdateUserOtp(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                otp         = ?,
                otp_exp     = ?,
                modified_dt = Now()
              WHERE  user   = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Otp,
		data.OtpExp,
		data.User,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *User) UpdateOtpStatus(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                status      = 'active',
                otp_exp     = date_add(now(), interval -1 day) ,
                modified_dt = Now()
              WHERE  user = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.User,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *User) SetUserLogStatus(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                logged      = 1,
                token       = ?,
                token_exp   = ?,
                modified_dt = Now()
              WHERE  user = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Token,
		data.TokenExp,
		data.User,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *User) GetByToken(ctx context.Context, db *sql.DB, token, who string) (*User, error) {
	r := `SELECT 
			ifnull(id,''), 
			ifnull(user,''), 
			ifnull(pass,''), 
			ifnull(status,''), 
			ifnull(otp,''), 
			ifnull(logged,0), 
			ifnull(token,''), 
			ifnull(created_dt,''), 
			ifnull(modified_dt,''),
			ifnull((otp_exp   <now()),0),
			ifnull((token_exp <now()),0)
		FROM  users WHERE user = ? AND  token = ? `
	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("Get", err)
		return nil, err
	}
	defer stmt.Close()
	var data User
	err = stmt.QueryRowContext(ctx, who, token).Scan(
		&data.ID,
		&data.User,
		&data.Pass,
		&data.Status,
		&data.Otp,
		&data.Logged,
		&data.Token,
		&data.Created,
		&data.Modified,
		&data.ExpiredOtp,
		&data.ExpiredToken,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Get NOT_FOUND", who)
		return nil, errors.New("Info not found")
	case err != nil:
		log.Println("Get", err)
		return nil, err
	}
	//sounds good ;-)
	return &data, nil
}

func (u *User) SetUserLogout(ctx context.Context, db *sql.DB, data *User) (bool, error) {
	//fmt
	r := `UPDATE users
                SET
                logged      = 0,
                token_exp   = date_add(Now(), INTERVAL -1 DAY),
                modified_dt = Now()
              WHERE  user = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.User,
	)
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}
