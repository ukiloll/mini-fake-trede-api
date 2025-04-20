package service

type bodyAsset struct{
	Name string `json:"name" form:"name" binding:"required"`
	Quntity float64  `json:"quntity" form:"quntity" binding:"required"`
}

type transition struct{
	UserId int
	AssetId int
	Price float64
	Quantity float64
}

type responseData struct{
	Price float64 `json:"p"`
	Stock string `json:"s"`
	Time int64 `json:"t"`
	Volume float64 `json:"v"`
}
type response struct{
	Data []responseData `json:"data"`
	Type string  `json:"type"`
}

type assetData struct{
	Id int 
	Type string
	responseData
}

type assetName string

type showAsset struct{
	Name string `json:"s"`
	StockPriceNow float64 `json:"stock_price_now"`
	Income string  `json:"income"`
}