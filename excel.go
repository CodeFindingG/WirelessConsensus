package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func StreamWriterFunc(contents []string) {
	//打开工作簿
	file, err := excelize.OpenFile("data.xlsx")
	if err != nil {
		file = excelize.NewFile()
		file.SetCellValue("Sheet1", "A1", "nodeNumber")
		file.SetCellValue("Sheet1", "B1", "pv")
		file.SetCellValue("Sheet1", "C1", "k")
		file.SetCellValue("Sheet1", "D1", "R")
		file.SetCellValue("Sheet1", "E1", "仿真次数")
		file.SetCellValue("Sheet1", "F1", "平均共识轮次")
	}
	sheet_name := "Sheet1"

	//获取流式写入器
	streamWriter, _ := file.NewStreamWriter(sheet_name)

	rows, _ := file.GetRows(sheet_name) //获取行内容
	cols, _ := file.GetCols(sheet_name) //获取列内容

	//将源文件内容先写入excel
	for rowid, row_pre := range rows {
		row_p := make([]interface{}, len(cols))
		for colID_p := 0; colID_p < len(cols); colID_p++ {
			if row_pre == nil {
				row_p[colID_p] = nil
			} else {
				row_p[colID_p] = row_pre[colID_p]
			}
		}
		cell_pre, _ := excelize.CoordinatesToCellName(1, rowid+1)
		if err := streamWriter.SetRow(cell_pre, row_p); err != nil {
			fmt.Println(err)
		}
	}

	//将新加contents写进流式写入器
	row := make([]interface{}, len(contents))
	for colID := 0; colID < len(contents); colID++ {
		row[colID] = contents[colID]
	}
	cell, _ := excelize.CoordinatesToCellName(1, len(rows)+1) //决定写入的位置
	if err := streamWriter.SetRow(cell, row); err != nil {
		fmt.Println(err)
	}

	//结束流式写入过程
	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
	}
	//保存工作簿
	if err := file.SaveAs("data.xlsx"); err != nil {
		fmt.Println(err)
	}
}
