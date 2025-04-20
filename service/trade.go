package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "github/ukilolll/trade/database"
	"github/ukilolll/trade/pkg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var(
	dbCon =db.Connect()
	Dashboard = make(map[assetName]*assetData)//map use struct must be pointer
	upgrader = &ws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {return true},
	}
)

func ShowDashboad(name assetName){
	data:= Dashboard[name] //map should cache data 
	fmt.Printf("%v price:%v time:%v\n",data.Stock,data.Price,time.Unix(data.Time,0).UTC())
}

func RunDashboard(){
	println("running board")
	var res response
	url := fmt.Sprintf("wss://ws.finnhub.io?token=%v",os.Getenv("FINHUB_TOKEN"))
	conn,_,err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Panicln(err)
	}
	rows,err := dbCon.Query("SELECT * FROM assets")
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	var i int
	for rows.Next(){
		var name assetName
		var data assetData
		err := rows.Scan(&data.Id,&name,&data.Type)
		if err != nil {
			log.Panic(err)
		}
		Dashboard[name]= &data
		i++;
	}
	//send data to query stock data
	for k,_ := range Dashboard{
		msg,_ := json.Marshal(map[string]any{"type": "subscribe", "symbol":k})
		conn.WriteMessage(ws.TextMessage,msg)
	}
	//LOOP get stock data
    //ticker := time.NewTicker(time.Second * 5)
	ticker := time.NewTicker(time.Millisecond * 100)
	for{
		select{
		case _ = <- ticker.C:
			conn.ReadJSON(&res)
			for k,_ := range Dashboard{
				for _,resData := range res.Data{
					if string(k) == resData.Stock{
						Dashboard[k].responseData = resData
						//ShowDashboad(k)
					}
				}
			}
		}

	}

}

func CheckAsset(ctx *gin.Context){
	var body bodyAsset
	var user_id,_ = strconv.Atoi(ctx.MustGet("id").(string))
	err := ctx.Bind(&body)
	if err != nil {
		ctx.String(http.StatusBadRequest,err.Error())
// Key: 'bodyAsset.Name' Error:Field validation for 'Name' failed on the 'required' tag
// Key: 'bodyAsset.Quntity' Error:Field validation for 'Quntity' failed on the 'required' tag
		ctx.Abort();return;
	}

	if body.Quntity < 1{
		pkg.BadRequest.SendErr("quntity should more than 0",ctx)
		ctx.Abort();return;
	}

	asset,ok := Dashboard[assetName(body.Name)]
	if !ok{
		ctx.String(http.StatusBadRequest,"invalid asset")
		ctx.Abort();return;
	}

	transition := &transition{
		UserId: user_id,
		AssetId: asset.Id,
		Quantity:body.Quntity ,
		Price: asset.Price,
	}

	ctx.Set("transition",transition)
	ctx.Next()
}

func makeTransition(tradeType string ,data *transition) error{
	command:=fmt.Sprintf(`INSERT INTO transition (trade_type, price, quantity, user_id, asset_id) VALUES ('%v', ?, ?, ?, ?);`,
	tradeType)
	_, err := dbCon.Exec(command,data.Price, data.Quantity, data.UserId, data.AssetId)
	return err
}

func BuyAsset(ctx *gin.Context) {
	data := ctx.MustGet("transition").(*transition)
	trans ,_ := dbCon.Begin()
	log.Println(data.AssetId,data.UserId)
	err := makeTransition("Buy",data)
	if err != nil {
		trans.Rollback()
		ctx.String(500,"")
		log.Panic(err)
	} 

	var oldQuantity,oldPrice float64 
	command := "SELECT quantity,price FROM user_assets WHERE user_id=? AND asset_id=?;"
	err = trans.QueryRow(command,data.UserId,data.AssetId).Scan(&oldQuantity,&oldPrice)

	if (err != sql.ErrNoRows && err != nil){
		trans.Rollback()
		ctx.String(500,"")
		log.Panic(err)
	}
	//first buy
	if(err == sql.ErrNoRows){
		trans.Exec("INSERT INTO user_assets(asset_id, user_id, quantity,price) VALUES(?,?,?,?);",
		data.AssetId,data.UserId,data.Quantity,data.Price)
	}else{//not first buy
		newQuantity:= data.Quantity+oldQuantity
		newPrice := ((oldQuantity*oldPrice) + (data.Quantity*data.Price)) / newQuantity

		_,err = trans.Exec("UPDATE user_assets SET quantity=?,price=? WHERE user_id=? AND asset_id=?;",
		newQuantity,newPrice,data.UserId,data.AssetId)
		if err != nil {
			trans.Rollback()
			ctx.String(500,"")
			log.Panic(err)
		}
	}

	ctx.String(200,"buy success")
	trans.Commit()
}


func SellAsset(ctx *gin.Context)  {
	data := ctx.MustGet("transition").(*transition)
	trans,_ := dbCon.Begin()

	err := makeTransition("Sell",data)
	if err != nil {
		trans.Rollback()
		ctx.String(500,"")
		log.Panic(err)
	}

	var oldQuantity float64
	command := "SELECT quantity FROM user_assets WHERE user_id=? AND asset_id=?;"
	err =dbCon.QueryRow(command,data.UserId,data.AssetId).Scan(&oldQuantity)

	if data.Quantity > oldQuantity{
		ctx.String(http.StatusBadRequest,"can't sell")
		return;
	}
	newQuantity := oldQuantity - data.Quantity

	_,err =trans.Exec("UPDATE user_assets SET quantity=? WHERE user_id=? AND asset_id=?;",
	newQuantity,data.UserId,data.AssetId)
	if err != nil {
		ctx.String(500,"")
		trans.Rollback()
		log.Panic(err)
	}

	ctx.String(200,"sell success")
	trans.Commit()
}


func LookAsset(ctx *gin.Context){
	var user_id,_ = strconv.Atoi(ctx.MustGet("id").(string))
	var res []showAsset
	var userQuantity,userPrice float64
	
	command := `SELECT name,quantity,price FROM user_assets 
	INNER JOIN assets ON user_assets.asset_id=assets.asset_id 
	WHERE user_id=?;`
	rows ,err := dbCon.Query(command,user_id)
	defer rows.Close()
	if err != nil {
		log.Panic(err)
	}

	for rows.Next(){
		var r showAsset
		rows.Scan(&r.Name,&userQuantity,&userPrice)
		res = append(res,r)
	}

	if(len(res) == 0){
		ctx.String(http.StatusNotFound,"user have no asset now!")
		return
	}

	conn,err := upgrader.Upgrade(ctx.Writer,ctx.Request,nil)	
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second)

	for{
	select{
	case _ = <-ticker.C:
		for i,v := range res{
			data := Dashboard[assetName(v.Name)]
			difference := data.Price/userPrice
			NowQuantity := userQuantity*difference
	
			var showPercent string
			percent := (difference-1)*100//difference 0.0 - infintie
			if percent > 0{
				showPercent = fmt.Sprintf("+%.6f",percent)
			}else{
				showPercent = fmt.Sprintf("%.6f",percent)
			}
	
			res[i].Income = fmt.Sprintf("%.6f(%v%%)",NowQuantity,showPercent)
			res[i].StockPriceNow = data.Price
		}
		conn.WriteJSON(res)
	}
	}	

}
