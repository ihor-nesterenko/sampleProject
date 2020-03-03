package db

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

type csvQI struct {
	CsvReader *csv.Reader
	CsvWriter *csv.Writer
}

func Init(connPath string) (QI, error) {
	file, err := os.OpenFile(connPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", connPath)
	}

	return csvQI{
		CsvReader: csv.NewReader(file),
		CsvWriter: csv.NewWriter(file),
	}, nil
}

func (c csvQI) UserQI() UserQI {
	return &CsvUserQI{c}
}

type CsvUserQI struct {
	csvQI
}

func (c CsvUserQI) SaveUser(newUser User) error {
	oldUser, err := c.findUser(newUser.Login)
	if err != nil {
		return errors.Wrap(err, "failed to find user by login")
	}
	if oldUser != nil {
		return errors.New("user with such login already exists")
	}

	err = c.CsvWriter.Write([]string{newUser.Login, newUser.Password})
	if err != nil {
		return errors.Wrap(err, "failed to save new user")
	}

	c.CsvWriter.Flush()
	return c.CsvWriter.Error()
}

func (c CsvUserQI) GetUser(login string) (*User, error) {
	return c.findUser(login)
}

func (c CsvUserQI) findUser(nickname string) (*User, error) {
	rawRecords, err := c.CsvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read users")
	}

	for _, rawRecord := range rawRecords {
		if rawRecord[0] == nickname {
			return &User{
				Login:    rawRecord[0],
				Password: rawRecord[1],
			}, nil
		}
	}

	return nil, nil
}
