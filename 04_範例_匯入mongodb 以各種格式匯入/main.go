package main

import (
	"fmt"
	"time"

	"github.com/zxfonline/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type NodeA struct {
	Node1 int
	Node2 string
}
type NodeB struct {
	NodeAcode int
	NodeBcode string
	NodeASon  map[float64]*NodeA
	NodeBSon  []*NodeA
	NodeCcode time.Time
}
type NodeB1 struct {
	NodeAcode  int
	NodeBcode1 string
	NodeASon   map[float64]*NodeA
	NodeBSon   []*NodeA
}

func main() {
	node := &NodeB{NodeAcode: 21, NodeBcode: "21", NodeCcode: time.Now()}
	fmt.Println(node.NodeCcode.Format("2006-01-02 15:04:05"))

	sonst := make(map[float64]*NodeA)
	sonst[1.1] = &NodeA{Node1: 111, Node2: "111"}
	sonst[2.2] = &NodeA{Node1: 112, Node2: "112"}
	node.NodeASon = sonst
	sona := make([]*NodeA, 0)
	sona = append(sona, &NodeA{Node1: 211, Node2: "211"})
	sona = append(sona, &NodeA{Node1: 212, Node2: "212"})
	node.NodeBSon = sona

	b, err := json.Marshal(&node)
	if err != nil {
		panic(err)
	}
	fmt.Println("insert struct json=", string(b))

	var f interface{}
	err = json.Unmarshal([]byte(b), &f)
	if err != nil {
		panic(err)
	}
	var session *mgo.Session
	session, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	coll := session.DB("test").C("test")

	err = coll.Insert(f)
	if err != nil {
		panic(err)
	}

	var q []bson.M
	coll.Find(nil).All(&q)
	for _, info := range q {
		b, err = json.Marshal(info)
		if err != nil {
			panic(err)
		}
		fmt.Println("parse struct json=", string(b))
		nodeN := &NodeB{}
		err = json.Unmarshal([]byte(b), nodeN)
		if err != nil {
			panic(err)
		}
		//test content same?
		b, err = json.Marshal(&node)
		if err != nil {
			panic(err)
		}
		fmt.Println("new struct json=", string(b))
		fmt.Println(nodeN.NodeCcode.Format("2006-01-02 15:04:05"))
	}
}
