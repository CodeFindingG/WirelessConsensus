package main

import (
	"colorout"
	"fmt"
	"math/rand"
	"strconv"
)

func simulationImprovedInit() {
	for i := 0; i < nodeNumber; i++ {
		nodes[i].State = "C"                          // The initial status
		nodes[i].CountV = 0                           // The initial status
		nodes[i].CV = "null"                          // The initial status
		nodes[i].PosX = rand.Float64() * Conf.maxPosX // Generate X position
		nodes[i].PosY = rand.Float64() * Conf.maxPosY // Generate Y position
		nodes[i].ifSelected = false
		nodes[i].ifBad = false
	}
	for i := 0; i < Conf.badNodes; i++ {
		tmp := rand.Intn(Conf.nodeNumber)
		for nodes[tmp].ifBad {
			tmp = rand.Intn(Conf.nodeNumber)
		}
		nodes[tmp].ifBad = true
	}
	countRound = 0  // Init
	countLeader = 0 // Init
}
func checkLImprove(nodeID int) int {
	ans := 0
	for i := 0; i < nodeNumber; i++ {
		if i == nodeID || nodes[i].State != "L" {
			continue
		}
		if nodes[i].ifBad || nodes[i].ifSelected {
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
func checkMessageImprove(nodeID int) bool {
	for i := 0; i < nodeNumber; i++ {
		if i == nodeID || nodes[i].ifMessage == false { // 没发的不管
			continue
		}
		if nodes[i].ifBad {
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
func roundPrint() {
	if ifPrint {
		fmt.Print(colorout.Green("Now round" + strconv.Itoa(countRound)))
	}
	for i := 0; i < nodeNumber; i++ { // 遍历所有节点
		nodes[i].ifMessage = false
		if ifPrint {
			fmt.Print(colorout.Green(" " + nodes[i].State + nodes[i].CV))
		}
	}
	if ifPrint {
		fmt.Println(" ")
	}
}
func simulationImproved() {
	ifContinue := true
	simulationImprovedInit()
	for ifContinue {
		for i := 0; i < nodeNumber; i++ {
			if nodes[i].ifSelected {
				continue
			}
			nodes[i].State = "C"
			nodes[i].CountV = 0
		}
		for i := 0; i < nodeNumber; i++ { // 循环所有节点
			if nodes[i].State == "C" { // 如果位于C状态。
				ifTran := rand.Float64() // Generate a [0-1) number用于看要不要传消息
				if ifTran <= Conf.pv {   //发消息
					nodes[i].ifMessage = true // 节点发送一条消息
				}
			}
		}
		for i := 0; i < nodeNumber; i++ { // 循环所有节点
			if nodes[i].State == "C" && nodes[i].ifMessage == false { // 如果位于C状态。且没发消息所以在监听。
				if checkMessageImprove(i) { // 有消息了
					nodes[i].State = "S" // 开始静默
				} else {
					nodes[i].CountV++
					if float64(nodes[i].CountV) > Conf.k {
						nodes[i].State = "L"
					}
				}
			}
		}
		countRound += 1 // 执行完时隙1了，轮次+1
		roundPrint()
		// 检查是否符合要求
		mark := true
		for i := 0; i < nodeNumber; i++ {
			if nodes[i].State == "S" {
				if checkLImprove(i) != 1 { // 检查范围内本轮有几个L
					mark = false
					break
				}
			}
		}
		if !mark { // 这轮如果报废了就继续下一轮
			continue
		}
		countLeader += 1
		for j := 0; j < nodeNumber; j++ {
			if nodes[j].State == "L" {
				nodes[j].ifSelected = true //被选中了
			}
		}

		if countLeader > Conf.badNodes {
			tot += float64(countRound)
			cnt += 1
			fmt.Print(colorout.Cyan("已结束，共识耗时轮次" + strconv.Itoa(countRound) + " 节点总数：" + strconv.Itoa(Conf.nodeNumber) + " 坏节点总数：" + strconv.Itoa(Conf.badNodes)))
			fmt.Println(colorout.Cyan(" 平均耗时轮次" + strconv.FormatFloat(tot/cnt, 'f', 10, 32)))
			simulationImprovedInit()
			//time.Sleep(2 * time.Second)
		}
	}
}
