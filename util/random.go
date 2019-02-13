package util

import (
	"database/sql"
	"math/rand"
)

const rs2Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkCount(rowList *sql.Rows) (count int) {
	for rowList.Next() {
		err := rowList.Scan(&count)
		checkErr(err)
	}
	return count
}

func GenerateOriginToken(tokenLength int) (string, error) {

	var tokenId string

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		return "", err
	}

	for {
		tokenId = randString(tokenLength)

		rowList, err := db.Query(`select count(*) as count from origin_token where token_id = ?`,
			tokenId,
		)
		if err != nil {
			return "", err
		}

		if checkCount(rowList) < 1 {
			_, err := db.Exec(
				`insert into origin_token (token_id, created_at) VALUES (?, datetime('now', 'localtime'))`,
				tokenId,
			)
			if err != nil {
				return "", err
			}

			break
		}
	}

	return tokenId, nil
}

func randString(length int) string {
	returnString := make([]byte, length)
	for i := range returnString {
		returnString[i] = rs2Letters[rand.Intn(len(rs2Letters))]
	}
	return string(returnString)
}
