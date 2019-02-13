package entity

import (
	"database/sql"
	"encoding/json"
	"log"

	. "../util"
)

type Board struct {
	Row  []Row  `json:"row"`
	Tile []Tile `json:"tile"`
}

func getBoardFromDb(board *Board) error {

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Printf("Couldnt open database file. FIlePath: " + DbPath)
		return err
	}
	defer db.Close()

	rowList, err := db.Query(`select * from row`)
	if err != nil {
		panic(err)
	}

	for rowList.Next() {
		var rowId string
		var title string
		if err := rowList.Scan(&rowId, &title); err != nil {
			return err
		}
		board.Row = append(board.Row, Row{RowID: rowId, Title: title})
	}

	tileList, err := db.Query(`select * from tile`)
	if err != nil {
		return err
	}

	for tileList.Next() {
		var tileId string
		var title string
		var rowId string
		var text string
		if err := tileList.Scan(&tileId, &title, &rowId, &text); err != nil {
			return err
		}
		board.Tile = append(board.Tile, Tile{TileID: tileId, Title: title, RowID: rowId, Text: text})
	}

	return nil
}

func (board *Board) GetBoardFromDb() {
	err := getBoardFromDb(board)
	if err != nil {
		log.Printf("Error occuerred during database operation.")
		panic(err)
	}
}

func (board *Board) ToJsonString() string {
	jsonBytes, err := json.Marshal(board)
	if err != nil {
		log.Println("JSON Marshal error:", err)
		panic(err)
	}

	return string(jsonBytes)
}

func (board *Board) WriteBoardToDb() error {
	prevBoard := Board{}
	prevBoard.GetBoardFromDb()

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		log.Printf("Couldnt open database file. FIlePath: " + DbPath)
		return err
	}
	defer db.Close()

	newRowIdList := make(map[string]struct{})
	for _, row := range board.Row {
		_, err := db.Exec("replace into row(row_id,title) values(?,?);", row.RowID, row.Title)
		if err != nil {
			return err
		}

		newRowIdList[row.RowID] = struct{}{}
	}

	for _, row := range prevBoard.Row {
		if _, exist := newRowIdList[row.RowID]; !exist {
			_, err := db.Exec("delete from row where row_id = ?;", row.RowID)
			if err != nil {
				return err
			}
		}
	}

	newTileIdList := make(map[string]struct{})
	for _, tile := range board.Tile {
		_, err := db.Exec("replace into tile(tile_id, title, row_id, text) values(?,?,?,?);", tile.TileID, tile.Title, tile.RowID, tile.Text)
		if err != nil {
			return err
		}

		newTileIdList[tile.TileID] = struct{}{}
	}

	for _, tile := range prevBoard.Tile {
		if _, exist := newTileIdList[tile.TileID]; !exist {
			_, err := db.Exec("delete from tile where tile_id = ?;", tile.TileID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
