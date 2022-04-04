package main

import (
	"math/rand"
	"time"
)

var Conf Config
var slotChange chan string
var ifSlotChange bool // this var used for break condition.!!!
var nodes []node      // 所有的节点
var nodeNumber = 4

const channelBuffer = 100

var ifTranMessage int // judge if there is a node send a message,count the number
var LNumber int       // Leader Number
var ifPrint bool      // Control test output
var tot float64
var cnt float64
var countLeader int // 计算要选出多少次Leader需要超过f
var countRound int

func main() {
	ifPrint = false
	ifSlotChange = false
	ifTranMessage = 0
	countRound = 0
	LNumber = 0
	tot = 0
	cnt = 0
	rand.Seed(int64(time.Now().Nanosecond())) // 随机数种子
	slotChange = make(chan string, channelBuffer)
	_ = ConfigInitial()
	nodeNumber = Conf.nodeNumber
	slotDuration = time.Duration(Conf.timeSlot) * time.Millisecond
	NodeInit()
	//time.Sleep(200 * time.Millisecond)
	//go timeGenerate() // 生成时间
	//go simulation() // 原版的
	//for {
	//	select {
	//	case _ = <-slotChange:
	//		//fmt.Println(v)
	//	}
	//}
	go simulationImproved()
	select {}
}
