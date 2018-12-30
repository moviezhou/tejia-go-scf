package main

import (
    "context"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

type PathParameters struct {
	Name string `string:"name"`
}

type DefineEvent struct {
    // test event define
    Key1 string `json:"key1"`
	Key2 string `json:"key2"`
	Path string `json:"path"`
	PathParameters PathParameters `json:"pathParameters"`

}

type Response struct{ 
    IsBase64 bool `json:"IsBase64"` 
    StatusCode uint32 `json:"statusCode"` 
    Headers map[string]string `json:"headers"` 
	Body string `json:"body"` 
} 


func getMarket(ctx context.Context, event DefineEvent) (Response, error) {
	fmt.Println("key1:", event)
	fmt.Println(ctx)
	fmt.Println("key2:", event.PathParameters.Name)

	var market string
	market = event.PathParameters.Name

	db, err := sql.Open("mysql", "root:JsJs6773@tcp(cd-cdb-hhl9rz48.sql.tencentcdb.com:63394)/tejia?charset=utf8")
	checkErr(err)

	// Query market
	query := "select * from market where market_name=?"
	rows, err := db.Query(query, market)
	checkErr(err)
	
	for rows.Next() {
		var id int
		var market_name string
		var created string
		err = rows.Scan(&id, &market_name, &created)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(market)
		fmt.Println(created)
	}

	checkErr(err)
	defer db.Close()

    ret := Response{} 
    ret.IsBase64 = false 
    ret.StatusCode = 200 
    ret.Headers = map[string]string{} 
    ret.Headers["Content-Type"] = "application/json" 
    
	ret.Body = "Ok"
    return ret, nil 
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
    // Make the handler available for Remote Procedure Call by Cloud Function
    cloudfunction.Start(getMarket)
}