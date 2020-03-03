/*
	CSV was used as storage for ability to git clone and run without any db connection setup
	Hardcoded filepath was used for simplicity as well
*/
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/ihor-nesterenko/sampleProject/db"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	salt = "was put into global variable for simplicity"

	loginFlag    = "login"
	passwordFlag = "password"

	filepath = "./db.csv"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     loginFlag,
		Usage:    "User nickname",
		Value:    "",
		Required: true,
	},
	&cli.StringFlag{
		Name:     passwordFlag,
		Usage:    "User password",
		Value:    "",
		Required: true,
	},
}

func main() {
	db, err := db.Init(filepath)
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "login/registration flow",
		Usage: "resister and login user",
		Commands: []*cli.Command{
			{
				Name:   "login",
				Usage:  "login user",
				Flags:  flags,
				Action: login(db),
			},
			{
				Name:   "registration",
				Usage:  "register user",
				Flags:  flags,
				Action: register(db),
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func login(qi db.QI) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		login := c.String(loginFlag)
		password := c.String(passwordFlag)

		user, err := qi.UserQI().GetUser(login)
		if err != nil {
			log.Println("failed to get user")
			return nil
		}

		if user == nil {
			log.Println("invalid login or password")
			return nil
		}
		passwordEncrypted := encryptPassword(password)
		if user.Password != passwordEncrypted {
			log.Println("invalid login or password")
			return nil
		}

		log.Println("Login successful")
		return nil
	}
}

func register(qi db.QI) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		login := c.String(loginFlag)
		password := c.String(passwordFlag)
		user := db.User{
			Login:    login,
			Password: password,
		}

		err := user.Validate()
		if err != nil {
			log.Println("both login and password must be set")
			return nil
		}

		user.Password = encryptPassword(password)

		err = qi.UserQI().SaveUser(user)
		if err != nil {
			log.Println(err)
			return nil
		}

		log.Println("Registration successful")
		return nil
	}
}

func encryptPassword(password string) string {
	saltedPassword := password + salt
	return hex.EncodeToString(sha256.New().Sum([]byte(saltedPassword)))
}
