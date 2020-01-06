package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

// init(dbLogin, truncate tables) -> run tests -> clear data(truncate tables)

func clearTables(){
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}


func TestMain(m *testing.M){
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T){
	t.Run("Add",testAddUser)
	t.Run("Get",testGetUser)
	t.Run("Delete",testDeleteUser)
	t.Run("reGet",testReGetUser)
}
func testAddUser(t *testing.T){
	err:=AddUserCredential("hikari","123")
	if err!=nil{
		t.Errorf("Error of addUser: %v",err)
	}
}
func testGetUser(t *testing.T){
	pwd,err:=GetUserCredential("hikari")
	if pwd!="123" || err!=nil{
		t.Errorf("Error of getUser")
	}
}
func testDeleteUser(t *testing.T){
	err:=DeleteUser("hikari","123")
	if err!=nil{
		t.Errorf("Error of delUser: %v",err)
	}
}
func testReGetUser(t *testing.T){
	pwd,err:=GetUserCredential("hikari")
	if err!=nil{
		t.Errorf("Error of ReGetUser: %v",err)
	}
	if pwd!=""{
		t.Errorf("Delete user test failed")
	}
}

func testAddComments(t *testing.T){
	vid:="12345"
	aid:=1
	content:="I like this video first"
	err:=AddNewComments(vid,aid,content)
	if err!=nil{
		t.Errorf("Error of AddComments: %v",err)
	}
}
func testListComments(t *testing.T){
	vid:="12345"
	from :=1514764800
	to,_:=strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000,10))
	res,err:=ListComments(vid,from,to)
	if err!=nil{
		t.Errorf("Errof of list Comments: %v",err)
	}
	for i,ele:=range res{
		fmt.Printf("comment %d, %v \n",i,ele)
	}
}



func TestComments(t *testing.T){
	clearTables()
	t.Run("AddUser",testAddUser)
	t.Run("AddComments",testAddComments)
	t.Run("ListComments",testListComments)
}





