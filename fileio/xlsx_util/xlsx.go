package xlsx_util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tealeg/xlsx"
)

type OnRowRender func(*xlsx.Row, interface{})

func ExportToTableExcel(sheetName string, array interface{}, onRowRender OnRowRender) (*bytes.Buffer, error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error
	file = xlsx.NewFile()
	sheet, err = file.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}
	if onRowRender != nil {
		for _, item := range array.([]interface{}) {
			row = sheet.AddRow()
			onRowRender(row, item)

		}
	} else {

		dataMapList, ok := array.([]map[string]interface{})
		if !ok {

			buff, err := json.Marshal(array)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(buff, &dataMapList)
			if err != nil {
				return nil, err
			}
		}

		columns := make([]string, 0)

		for _, dataRow := range dataMapList {
			row = sheet.AddRow()
			if len(columns) == 0 {
				for key := range dataRow {
					columns = append(columns, key)
				}
				sort.Strings(columns)
			}
			for _, column := range columns {
				cell := row.AddCell()
				cell.SetString(
					fmt.Sprintf("%v", dataRow[column]),
				)
			}
		}
	}

	bBuff := new(bytes.Buffer)

	err = file.Write(bBuff)
	if err != nil {
		return nil, err
	}

	return bBuff, nil
}
