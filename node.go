package main

import (
	"GZsimulation/colorout"
	"fmt"
	"math/rand"
	"strconv"
)

type node struct {
	State     string  // A:mining C:candidate S:silent L:leader
	CountV    int     // Whatever I don't Know
	CV        string  // null / 0 / 1
	ifMessage bool    // 是否发了消息。
	Message   string  // Not use
	Block     string  // Not use
	PosX      float64 // X postion
	PosY      float64 // Y postion
}

func slotBreak() {
	ifSlotChange = false
	ifTranMessage = 0 // 重置消息池
	if ifPrint {
		fmt.Print(colorout.Cyan("Now time Slot" + strconv.Itoa(timeSlot)))
	}
	for i := 0; i < nodeNumber; i++ { // 遍历所有节点
		nodes[i].ifMessage = false
		if ifPrint {
			fmt.Print(colorout.Cyan(" " + nodes[i].State + nodes[i].CV))
		}
	}
	if ifPrint {
		fmt.Println(" ")
	}
}
func checkL(nodeID int) int {
	ans := 0
	for i := 0; i < nodeNumber; i++ {
		if i == nodeID || nodes[i].State != "L" {
			continue
		}
		dis := (nodes[i].PosX-nodes[nodeID].PosX)*(nodes[i].PosX-nodes[nodeID].PosX) + (nodes[i].PosY-nodes[nodeID].PosY)*(nodes[i].PosY-nodes[nodeID].PosY)
		if dis > Conf.receiveR*Conf.receiveR {
			continue
		} else {
			ans += 1
		}
	}
	return ans
}
func checkMessage(nodeID int) bool {
	for i := 0; i < nodeNumber; i++ {
		if i == nodeID || nodes[i].ifMessage == false { // 没发的不管
			continue
		}
		dis := (nodes[i].PosX-nodes[nodeID].PosX)*(nodes[i].PosX-nodes[nodeID].PosX) + (nodes[i].PosY-nodes[nodeID].PosY)*(nodes[i].PosY-nodes[nodeID].PosY)
		if dis > Conf.receiveR*Conf.receiveR {
			continue
		} else {
			return true
		}
	}
	return false
}

var tot float64
var cnt float64
var ifContinue bool = true
var data []string

func endProcess() {
	cnt += 1
	tmp := (timeSlot + 1) / 5
	tot += float64(tmp)
	fmt.Println("已结束,当前仿真耗时轮次:", tmp, " 平均达成共识轮次:", tot/cnt)
	if cnt == Conf.simulationTimes {
		ifContinue = false
		data = []string{}
		data = append(data, strconv.Itoa(nodeNumber), strconv.FormatFloat(Conf.pv, 'f', 2, 64), strconv.FormatFloat(Conf.k, 'f', 2, 64), strconv.FormatFloat(Conf.receiveR, 'f', 2, 64), strconv.FormatFloat(Conf.simulationTimes, 'f', 0, 64), strconv.FormatFloat(tot/cnt, 'f', 2, 64))
		StreamWriterFunc(data)
		fmt.Println(colorout.Blue("结束程序,本次数据将追加到data.xlsx!"))
	}
}
func simulation() {
	for ifContinue {
		// Slot 1
		for ifSlotChange == false {
			for i := 0; i < nodeNumber; i++ { // 遍历所有节点
				if nodes[i].State == "C" || nodes[i].State == "L" {
					// 发送消息
					ifTranMessage += 1
				} else if nodes[i].State == "A" {
					if ifTranMessage >= 1 {
						nodes[i].State = "S"
					} else {
						nodes[i].State = "C"
					}
				}
			}
		}
		slotBreak()
		// Slot 2
		for ifSlotChange == false {
			for i := 0; i < nodeNumber; i++ { // 循环所有节点
				if nodes[i].State == "C" { // 如果位于C状态。
					ifTran := rand.Float64() // Generate a [0-1) number用于看要不要传消息
					//fmt.Println(ifTran)
					if ifTran <= Conf.pv { //发消息
						nodes[i].ifMessage = true // 节点发送一条消息 消息池+1
					} else { //监听
						if checkMessage(i) { // 有消息了
							nodes[i].State = "S" // 开始静默
						} else {
							//tmp := Conf.k * math.Log10(Conf.nn[0])
							tmp := Conf.k //就是搞了多少次countv++呗。
							if float64(nodes[i].CountV) > tmp {
								nodes[i].State = "L"
								LNumber += 1
							}
						}
					}
					nodes[i].CountV++
				}
			}
		}
		slotBreak()
		// Slot 3
		for ifSlotChange == false {
			for i := 0; i < nodeNumber; i++ {
				if nodes[i].State == "L" {
					// 跳过，其他节点直接检查有几个L
				}
				if nodes[i].State == "S" {
					if checkL(i) == 1 { // 检查范围内有几个L
						nodes[i].CV = "0"
					} else {
						nodes[i].CV = "1"
					}
				}
			}
		}
		slotBreak()
		// Slot 4
		for ifSlotChange == false {
			for i := 0; i < nodeNumber; i++ {
				if nodes[i].State == "L" {
					if checkMessage(i) {
						nodes[i].CV = "1"
					} else {
						nodes[i].CV = "0"
					}
				}
				if nodes[i].State == "S" && nodes[i].CV == "1" {
					nodes[i].ifMessage = true
				}
			}
		}
		slotBreak()
		// Slot 5
		for ifSlotChange == false {
			countL := 0
			for i := 0; i < nodeNumber; i++ {
				if nodes[i].State == "L" {
					countL += 1 // 有几个主节点？
				}
			}
			if countL == 1 { // 只有1个
				for i := 0; i < nodeNumber; i++ { // 遍历所有节点
					fmt.Print(colorout.Green(nodes[i].State + nodes[i].CV + " "))
				}
				fmt.Println("")
				for i := 0; i < nodeNumber; i++ {
					if nodes[i].State == "L" && nodes[i].CV == "0" { // 满足调节。结束程序吧。
						endProcess()
						//ifContinue = false
						timeSlot = 1 // 持续
					}
				}
			}
			for i := 0; i < nodeNumber; i++ {
				nodes[i].State = "A"
				nodes[i].CountV = 0
				nodes[i].CV = "null"
			}
			for ifSlotChange == false {
			}
		}
		slotBreak()
		LNumber = 0
	}
}

func NodeInit() {
	nodes = make([]node, nodeNumber) // 创建节点数组
	for i := 0; i < nodeNumber; i++ {
		//nodes[i].SlotChanges = make(chan string, channelBuffer)
		//nodes[i].Noise = make(chan float64)
		//nodes[i].Sv = make(chan float64)
		nodes[i].State = "A"                          // The initial status
		nodes[i].CountV = 0                           // The initial status
		nodes[i].CV = "null"                          // The initial status
		nodes[i].PosX = rand.Float64() * Conf.maxPosX // Generate X position
		nodes[i].PosY = rand.Float64() * Conf.maxPosY // Generate Y position

		//go nodeStart(i) // 启动很多很多的节点 Not Use! We just simulate it by for circle.
	}
}
