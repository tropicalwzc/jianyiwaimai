// websockets.go
package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"log"
	"net/http"
	"strings"

	_ "github.com/alexbrainman/odbc"
	"github.com/axgle/mahonia"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024000,
	WriteBufferSize: 1024000,
}

type account struct {
	aid     string
	level   string
	balance string
}
type sender struct {
	sid     string
	sname   string
	sphone  string
	aid     string
	level   string
	balance string
}
type orders struct {
	oid, raddr, rname, rphone, caddr, cname, cphone, fname, price, mood string
}
type food struct {
	fid   string
	fname string
	price string
	rid   string
}
type restaurant struct {
	rid     string
	rname   string
	raddr   string
	rphone  string
	aid     string
	level   string
	balance string
}
type customer struct {
	cid     string
	cname   string
	cphone  string
	caddr   string
	aid     string
	level   string
	balance string
}
type foodqu struct {
	fname                  string
	price                  string
	rname                  string
	raddr                  string
	rphone                 string
	fid                    string
	oid, sid, sphone, mood string
}

func floatTostring(aim float32) string {
	return fmt.Sprintf("%f", aim)
}
func intTostring(aim int) string {
	return fmt.Sprintf("%d", aim)
}
func utf8togbk(aim string) string {
	return mahonia.NewEncoder("gbk").ConvertString(aim)
}
func gbktoutf8(aim string) string {
	return mahonia.NewDecoder("gbk").ConvertString(aim)
}

//顾客付钱
func customerpay(money string, cid string) {
	basesearch := "with aimid as(" +
		"select aid " +
		"from Customer " +
		"where cid='" + cid + "') " +
		"update account " +
		"set balance = balance - " + money +
		"  where aid in(select aid from aimid );"

		//	println(basesearch)
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func searchrestaurant(rid string) string {
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := "select rname,balance,level,rid,raddr,rphone from restaurant,account where rid='" + rid + "' and restaurant.aid = account.aid"
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*restaurant
	//遍历每一行
	for rows.Next() {
		var row = new(restaurant)
		rows.Scan(&row.rname, &row.balance, &row.level, &row.rid, &row.raddr, &row.rphone)
		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		restr += ar.rname + "+" + ar.balance + "+" + ar.level + "+" + ar.rid + "+" + ar.raddr + "+" + ar.rphone + "\n"
	}
	return gbktoutf8(restr)
}
func restaurantfoodlist(rid string) string {
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := "select fid,fname,price from Food where rid = '" + rid + "'"
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*food
	//遍历每一行
	for rows.Next() {
		var row = new(food)
		rows.Scan(&row.fid, &row.fname, &row.price)
		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		restr += ar.fid + "+" + ar.fname + "+" + ar.price + "\n"
	}
	return gbktoutf8(restr)
}

func foodchangeprice(newprice string, fid string) {
	basesearch := "update food set price = '" + newprice + "' where fid = '" + fid + "';"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}
func foodchangename(newname string, fid string) {
	basesearch := "update food set fname = '" + newname + "' where fid = '" + fid + "';"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func newfoodinsert(newname string, newprice string, rid string) {
	fid := idgenerator(10)
	basesearch := "insert into food values('" + fid + "','" + newprice + "','" + newname + "','" + rid + "');"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func newrestinsert(newrid string, newname string, newaddr string, newphone string, newaid string) {

	basesearch := "insert into restaurant values('" + newrid + "','" + newaddr + "','" + newphone + "','" + newname + "','" + newaid + "');"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func newcustomerinsert(newcid string, newname string, newaddr string, newphone string, newaid string) {

	basesearch := "insert into customer values('" + newcid + "','" + newname + "','" + newphone + "','" + newaddr + "','" + newaid + "');"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}
func newsenderinsert(newsid string, newname string, newphone string, newaid string) {

	basesearch := "insert into sender values('" + newsid + "','" + newname + "','" + newphone + "','" + newaid + "');"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//饭店收钱依据fid
func restaurantrecev(money string, fid string) {
	basesearch := "with restid as(select rid from Food where fid='" + fid + "')," +
		"restaid as(select aid from Restaurant where rid in(select rid from restid ) )" +
		"update account " +
		"set balance = balance + " + money +
		"  where aid in(select aid from restaid);"

		//	println(basesearch)
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//骑手收钱sid
func senderrevev(money string, sid string) {
	basesearch := "with aimid as(" +
		"select aid " +
		"from Sender " +
		"where sid='" + sid + "') " +
		"update account " +
		"set balance = balance + " + money +
		"where aid in(select aid from aimid );"

	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func idgenerator(leng int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	oidbcode := [40]string{"A", "B", "C", "D", "E", "F", "G", "H", "i", "J", "K", "f", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "e", "q", "t", "r", "y"}
	oid := ""
	for i := 0; i < leng; i++ {
		oid += oidbcode[r.Intn(40)]
	}
	return oid
}

//增加新订单
func insertneworder(cid string, fid string) {

	oid := idgenerator(10)
	basesearch := "insert into orders values('" + oid + "','" + cid + "', NULL ,'" + fid + "',0)"
	//	fmt.Println(oid)
	//	fmt.Println(basesearch)
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//运送订单
func deliverorder(oid string, sid string) {
	basesearch := "update orders set sid = '" + sid + "' where oid = '" + oid + "';"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//删除套餐
func deletefood(fid string) {
	basesearch := "delete from food where fid = '" + fid + "';"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//删除已完成订单
func deletefinishedorder(sid string) {
	basesearch := "delete from orders where sid = '" + sid + "' and mood != 0;"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

//收到并且评价订单
func recevandcomment(oid string, mood string) {
	basesearch := "update orders set mood = " + mood + " where oid = '" + oid + "' ;"
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	db.Query(basesearch)
}

func bookedQuary(customerid string, nulloption int) string {
	basesearch := ""
	if nulloption == 0 {
		basesearch = "select fname,price,rname,food.fid,oid,sphone " +
			"from food,Restaurant,orders,Sender " +
			"where food.rid=Restaurant.rid and Orders.fid=Food.fid and Orders.sid=Sender.sid and cid = '" + customerid + "' and mood=0;"
	} else if nulloption == 1 {
		basesearch = "select fname,price,rname,food.fid,oid from food,Restaurant,orders " +
			"where food.rid=Restaurant.rid and Orders.fid=Food.fid and Orders.sid is NULL and cid = '" + customerid + "' and mood=0;"
	} else {
		basesearch = "select fname,price,rname,food.fid,oid,Sender.sid,sphone,mood from food,Restaurant,orders,Sender  " +
			"where food.rid=Restaurant.rid and Orders.fid=Food.fid and Orders.sid is not NULL and cid = '" + customerid + "' and mood=0 and Sender.sid=Orders.sid;"
	}

	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	var rowsData []*foodqu
	//遍历每一行
	for rows.Next() {
		var row = new(foodqu)
		if nulloption == 0 {
			rows.Scan(&row.fname, &row.price, &row.rname, &row.fid, &row.oid, &row.sphone)
		} else if nulloption == 1 {
			rows.Scan(&row.fname, &row.price, &row.rname, &row.fid, &row.oid)
		} else {
			rows.Scan(&row.fname, &row.price, &row.rname, &row.fid, &row.oid, &row.sid, &row.sphone, &row.mood)
		}

		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		if nulloption == 0 {
			restr += ar.fname + "+" + ar.price + "+" + ar.rname + "+" + ar.fid + "+" + ar.oid + "+" + ar.sphone + "\n"
		} else if nulloption == 1 {
			restr += ar.fname + "+" + ar.price + "+" + ar.rname + "+" + ar.fid + "+" + ar.oid + "\n"
		} else {
			restr += ar.fname + "+" + ar.price + "+" + ar.rname + "+" + ar.fid + "+" + ar.oid + "+" + ar.sid + "+" + ar.sphone + "+" + ar.mood + "\n"
		}

	}
	return gbktoutf8(restr)
}

//查询所有套餐
func foodQuary(opt string) string {

	basesearch := "select fname,price,rname,raddr,rphone,fid from food,Restaurant where food.rid=Restaurant.rid"
	if opt != "all" {
		basesearch += " and fname like '%" + opt + "%'"
	}
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	var rowsData []*foodqu
	//遍历每一行
	for rows.Next() {
		var row = new(foodqu)
		rows.Scan(&row.fname, &row.price, &row.rname, &row.raddr, &row.rphone, &row.fid)
		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		restr += ar.fname + "+" + ar.price + "+" + ar.rname + "+" + ar.raddr + "+" + ar.rphone + "+" + ar.fid + "\n"
	}
	return gbktoutf8(restr)
}

//根据sid搜索骑手
func searchsender(usrid string) string {
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := "select sname,balance,level,sid from sender,account where sid='" + usrid + "' and sender.aid = account.aid "
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*sender
	//遍历每一行
	for rows.Next() {
		var row = new(sender)
		rows.Scan(&row.sname, &row.balance, &row.level, &row.sid)
		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		restr += ar.sname + "+" + ar.balance + "+" + ar.level + "+" + ar.sid + "\n"
	}
	return gbktoutf8(restr)
}

//根据cid搜索顾客
func searchcustomer(usrid string) string {
	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := "select cname,balance,level,cid from Customer,account where cid='" + usrid + "' and Customer.aid = account.aid"
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*customer
	//遍历每一行
	for rows.Next() {
		var row = new(customer)
		rows.Scan(&row.cname, &row.balance, &row.level, &row.cid)
		rowsData = append(rowsData, row)
	}
	restr := ""
	//打印数组
	for _, ar := range rowsData {
		restr += ar.cname + "+" + ar.balance + "+" + ar.level + "+" + ar.cid + "\n"
	}
	return gbktoutf8(restr)
}

//根据cid查询账户
func accountQuary(classname string) string {

	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := "select * from [account] where aid='"
	basesearch += classname
	basesearch += "'"
	if string(classname) == "all_class" {
		basesearch = "select * from [account] "
	}
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	var rowsData []*account
	//遍历每一行
	for rows.Next() {
		var row = new(account)
		rows.Scan(&row.aid, &row.balance, &row.level)
		rowsData = append(rowsData, row)
	}
	restr := ""
	resnum := 0
	//打印数组
	for _, ar := range rowsData {

		restr += ar.aid + "+" + ar.balance + "+" + ar.level + "\n"
		resnum++
	}
	if resnum == 0 {
		restr = "未找到名字为" + classname + "的账户"
	}
	return gbktoutf8(restr)
}

//骑手查询订单
func orderQuary(ordername string, cat string) string {

	db, err := sql.Open("odbc", "driver={sql server};server=LAPTOP-2OEQ2KS5\\PROJ2008;port=1433;uid=control;pwd=qwerzxcv123;database=waimai")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//通过连接对象执行查询
	basesearch := " select oid,raddr,rname, rphone, caddr,cname, cphone,fname,price,mood" +
		" from orders,Customer,food,Restaurant " +
		" where Orders.cid = Customer.cid and food.fid=Orders.fid and food.rid=Restaurant.rid "
	if cat != "all" {
		basesearch += " and orders.sid = " + cat
	}
	if string(ordername) != "all" {
		basesearch += " and oid = " + ordername
	}
	if cat == "all" && ordername == "all" { // 寻找无人配送的订单
		basesearch += " and sid is NULL"
	}
	rows, err := db.Query(basesearch)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	var rowsData []*orders
	//遍历每一行
	for rows.Next() {
		var row = new(orders)
		rows.Scan(&row.oid, &row.raddr, &row.rname, &row.rphone, &row.caddr, &row.cname, &row.cphone, &row.fname, &row.price, &row.mood)
		rowsData = append(rowsData, row)
	}
	restr := ""
	resnum := 0
	//打印数组
	for _, ar := range rowsData {

		restr += ar.oid + "+" + ar.raddr + "+" + ar.rname + "+" + ar.rphone + "+" + ar.caddr + "+" + ar.cname + "+" + ar.cphone + "+" + ar.fname + "+" + ar.price
		if cat == "all" {
			restr += "\n"
		} else {
			restr += "+" + ar.mood + "\n"
		}
		resnum++
	}

	return gbktoutf8(restr)
}
func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			strmsg := string(msg)
			strlines := strings.Split(strmsg, " ")
			//	fmt.Printf("%s\n", strmsg)
			displaystr := "type | msg"
			rply := true

			if len(strlines) > 2 {
				category := strlines[0]
				order := strlines[1]
				if category == "account" {
					if order == "select" {
						displaystr = "account_select|"
						displaystr += accountQuary(strlines[2])
					}
				}
				if category == "update" {
					if len(strlines) < 4 {
						//		fmt.Println("Too few")
						displaystr = "err|too few"
					} else {
						whom := strlines[2]
						acer := len(strlines) - 1
						valer := strlines[acer]
						//fmt.Println("restaurant money " + valer)
						if order == "customer" {
							customerpay(valer, whom)
						}
						if order == "restaurant" {
							restaurantrecev(valer, whom)
						}
						if order == "sender" {
							senderrevev(valer, whom)
						}
					}
					rply = false
				}
				if category == "sender" {
					//fmt.Println("sender recev ")
					if order == "select" {
						//fmt.Println("sender all ")
						displaystr = "sender_select|"
						displaystr += orderQuary(strlines[2], "all")
					}
					if order == "ongoing" {
						displaystr = "sender_ongoing|"
						displaystr += orderQuary("all", strlines[2])
					}
					if order == "pick" {
						//fmt.Println("pick " + strlines[2] + " " + strlines[3])
						deliverorder(strlines[2], strlines[3])
					}
					if order == "search" {
						displaystr = "sender_search|"
						displaystr += searchsender(strlines[2])
					}
					if order == "clfinish" {
						//fmt.Println("del finished order " + strlines[2])
						deletefinishedorder(strlines[2])
					}
					if order == "new" {
						newsenderinsert(strlines[2], strlines[3], strlines[4], strlines[5])
					}
				}
				if category == "book" {
					if order == "notpick" {
						displaystr = "book_notpick|"
						displaystr += bookedQuary(strlines[2], 1)
					}
					if order == "pick" {
						displaystr = "book_pick|"
						displaystr += bookedQuary(strlines[2], 2)
					}
					if order == "dilivering" {
						displaystr = "book_dilivering|"
						displaystr += bookedQuary(strlines[2], 2)
					}
				}
				if category == "customer" {
					if order == "search" {
						displaystr = "customer_search|"
						displaystr += searchcustomer(strlines[2])
					}
					if order == "request" {
						if len(strlines) < 4 {
							displaystr = "err|too few argument "
						} else {
							cidstr := strlines[2]
							fidstr := strlines[3]
							insertneworder(cidstr, fidstr)
						}
					}
					if order == "comment" {
						recevandcomment(strlines[2], strlines[3])
					}
					if order == "new" {
						newcustomerinsert(strlines[2], strlines[3], strlines[4], strlines[5], strlines[6])
					}
				}
				if category == "food" {
					if order == "select" {
						displaystr = "food_select|"
						//				fmt.Println("find food")
						displaystr += foodQuary("all")
					}
					if order == "search" {
						displaystr = "food_search|"
						displaystr += foodQuary(strlines[2])
					}
					if order == "changeprice" {
						//fmt.Println("price " + strlines[2] + " fid " + strlines[3])
						foodchangeprice(strlines[2], strlines[3]) // newprice fid
					}
					if order == "changename" {
						//fmt.Println("name " + strlines[2] + " fid " + strlines[3])
						foodchangename(strlines[2], strlines[3]) // newname fid
					}
					if order == "new" {
						newfoodinsert(strlines[2], strlines[3], strlines[4]) //fname fprice rid
					}
					if order == "delete" {
						deletefood(strlines[2]) //fid
					}
				}
				if category == "restaurant" {
					if order == "search" {
						//fmt.Println("find rest " + strlines[2])
						displaystr = "rest_search|"
						displaystr += searchrestaurant(strlines[2])
					}
					if order == "foodlist" {
						displaystr = "rest_foodlist|"
						displaystr += restaurantfoodlist(strlines[2])
					}
					if order == "new" {
						newrestinsert(strlines[2], strlines[3], strlines[4], strlines[5], strlines[6])
					}
				}
			}

			// Print the message to the console
			// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(displaystr))
			msg = []byte(displaystr)
			// Write message back to browser
			if rply == true {
				if err = conn.WriteMessage(msgType, msg); err != nil {
					return
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "C:\\go_workplace\\gowebtest\\websocket.html")
	})

	http.ListenAndServe(":7111", nil)
}
