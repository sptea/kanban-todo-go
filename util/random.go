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

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func GenerateOriginToken(tokenLength int) (string, error) {

	var tokenId string

	for {
		tokenId = randString(tokenLength)

		rowList, err := Db.Query(`select count(*) as count from origin_token where token_id = ?`,
			tokenId,
		)
		if err != nil {
			return "", err
		}

		if checkCount(rowList) < 1 {
			_, err := Db.Exec(
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
