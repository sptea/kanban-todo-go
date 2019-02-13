package entity

type Tile struct {
	TileID string `json:"tileId"`
	Title  string `json:"title"`
	RowID  string `json:"rowId"`
	Text   string `json:"text"`
}
