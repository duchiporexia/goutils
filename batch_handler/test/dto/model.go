package dto

type DogPo struct {
	Name string `msg:"name"`
	Age  int    `msg:"age"`
}

type Dog struct {
	Id   string `msg:"id"`
	Name string `msg:"name"`
	Age  int    `msg:"age"`
}
