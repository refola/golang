package core

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
)

type User struct {
	name       string
	data       map[string]interface{}
	passSHA512 []byte
}

func (u *User) GetInfo(what string) interface{} {
	return u.data[what]
}
func (u *User) SetInfo(what string, info interface{}) {
	u.data[what] = info
}

func passHash(pass string) []byte {
	hash := sha512.New()
	hash.Write([]byte(pass))
	return hash.Sum(nil)
}
func (u *User) ValidatePassword(guess string) bool {
	for i, v := range passHash(guess) {
		if u.passSHA512[i] != v {
			return false
		}
	}
	return true
}
func (u *User) ChangePassword(pass string) {
	u.passSHA512 = passHash(pass)
}

func (u *User) Name() string {
	return u.name
}

type UserStore struct {
	data DataStore
}

func NewUserStore(data DataStore, location string) *UserStore {
	return &UserStore{NewSubDataStore(location, data)}
}

func (us *UserStore) GetUser(name string) (*User, error) {
	raw, err := us.data.Load(name)
	if err != nil {
		return nil, err
	}
	u := new(User)
	buf := bytes.NewBuffer(raw)
	err = gob.NewDecoder(buf).Decode(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *UserStore) SaveUser(user *User) error {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(user)
	if err != nil {
		return err
	}
	return us.data.Save(user.name, buf.Bytes())
}

func (us *UserStore) DeleteUser(user *User) error {
	return us.data.Remove(user.name)
}
