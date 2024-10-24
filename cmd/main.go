package main

import (
	"fmt"
	"log"

	"proj1/internal/pkg/server"
	"proj1/internal/pkg/storage"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}
	srv := server.New(":8090", &s)
	srv.Start()

	stor2, err := storage.NewSliceStorage()
	if err != nil {
		log.Fatal(err)
	}
	// stor2.LoadFromFile("slice_storage.json")
	// stor2.LPush("key1", []string{"hhh", "gggg"})
	// stor2.LPush("key1", []string{"llll", "aaaa"})
	// stor2.RPush("key2", []string{"kkkk", "pppp"})
	// stor2.RAddToSet("key1", []string{"jjjj", "hhh", "gggg", "mmmm"})
	// fmt.Println(stor2)
	// fmt.Println(stor2.LGet("key1", 3))
	// fmt.Println(stor2.RPop("key1", 2, 78))
	// fmt.Println(stor2)
	// fmt.Println(stor2.LPop("key1", -6, 4))
	// fmt.Println(stor2.LPop("key2", 1))
	// fmt.Println(stor2)
	s.Set("key1", "val1")
	s.Set("key2", "754")
	res1, ok := s.Get("key1")
	if ok {
		fmt.Println(res1, s.GetKind("key1"))
	}
	res2, ok := s.Get("key2")
	if ok {
		fmt.Println(res2, s.GetKind("key2"))
	}
	err = stor2.SaveToFile("slice_storage.json")
	if err != nil {
		fmt.Println("Error saving storage:", err)
	}
	//s.SaveToFile("storage.json")

}
