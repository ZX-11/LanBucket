package main

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type TableModel struct {
	walk.ReflectTableModelBase
}

func (m *TableModel) Items() interface{} { return files }

var tableModel = new(TableModel)
var _ walk.ReflectTableModel = new(TableModel)
var fileUpload = make(chan struct{}, 8)
var leAddr *walk.LineEdit

func init() {
	go func() {
		for range fileUpload {
			tableModel.PublishRowsReset()
		}
	}()
}

var mainWindow = MainWindow{
	Title:   "LanBucket",
	Size:    Size{800, 600},
	MinSize: Size{800, 600},
	Layout:  HBox{},
	OnDropFiles: func(files []string) {
		for _, f := range files {
			if err := add(f); err != nil {
				fmt.Println(err)
			}
		}
		tableModel.PublishRowsReset()
	},
	Children: []Widget{
		Composite{
			Layout:  VBox{},
			MaxSize: Size{600, 600},
			Children: []Widget{
				Composite{
					Layout: HBox{},
					Children: []Widget{
						LineEdit{
							AssignTo: &leAddr,
							Text:     webAddr(),
							ReadOnly: true,
						},
						PushButton{
							Text:    "浏览",
							MaxSize: Size{60, 20},
							OnClicked: func() {
								cmd := exec.Command(`cmd`, `/c`, `start`, webAddr())
								cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
								cmd.Run()
							},
						},
						PushButton{
							Text:      "刷新",
							MaxSize:   Size{60, 20},
							OnClicked: tableModel.PublishRowsReset,
						},
						PushButton{
							Text:    "清空",
							MaxSize: Size{60, 20},
							OnClicked: func() {
								deleteAll()
								tableModel.PublishRowsReset()
							},
						},
					},
				},
				TableView{
					Alignment: AlignHNearVNear,
					Model:     tableModel,
					Columns: []TableViewColumn{
						TableViewColumn{
							DataMember: "Name",
							Title:      "文件名",
							Width:      450,
						},
						TableViewColumn{
							DataMember: "Size",
							Title:      "大小",
							Format:     "%v",
							Width:      100,
						},
					},
				},
			},
		},
		Composite{
			Layout:  VBox{},
			MinSize: Size{100, 600},
			DataBinder: DataBinder{
				DataSource: settings,
				AutoSubmit: true,
			},
			Children: []Widget{
				Label{
					Text:      "设置",
					Alignment: AlignHCenterVNear,
				},
				RadioButtonGroupBox{
					DataMember: "EnableUpload",
					Title:      "上传功能",
					Layout:     HBox{},
					Buttons: []RadioButton{
						RadioButton{
							Name:  "DisableUpload",
							Text:  "关闭",
							Value: false,
						},
						RadioButton{
							Name:  "EnableUpload",
							Text:  "开启",
							Value: true,
						},
					},
				},
				RadioButtonGroupBox{
					DataMember: "EnableIPv6",
					Title:      "使用IPv6地址",
					Layout:     HBox{},
					Buttons: []RadioButton{
						RadioButton{
							Name:  "DisableIPv6",
							Text:  "关闭",
							Value: false,
							OnClicked: func() {
								leAddr.SetText(webAddr())
							},
						},
						RadioButton{
							Name:  "EnableIPv6",
							Text:  "开启",
							Value: true,
							OnClicked: func() {
								leAddr.SetText(webAddr())
							},
						},
					},
				},
			},
		},
	},
}
