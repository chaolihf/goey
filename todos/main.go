package main

import (
	"goey"
	"strconv"
)

var (
	mainWindow *goey.MainWindow
)

func main() {
	mw, err := goey.NewMainWindow("Example", render())
	if err != nil {
		println(err.Error())
		return
	}
	defer mw.Close()
	mainWindow = mw

	goey.Run()
}

func update() {
	err := mainWindow.SetChildren(render())
	if err != nil {
		println("Error: ", err.Error())
	}
}

func render() []goey.Widget {
	labelText := "What needs to be done?"
	if len(Model) > 0 {
		labelText = labelText + "  (" + strconv.Itoa(len(Model)) + ")"
	}
	widgets := []goey.Widget{
		&goey.Label{Text: "What needs to be done:"},
		&goey.TextInput{Placeholder: "Enter todo description.", OnEnterKey: onNewTodoItem},
	}
	count1, count2 := getItemCounts()
	if count2 > 0 {
		widgets = append(widgets, &goey.HR{})
		widgets = append(widgets, &goey.Label{Text: "There are " + strconv.Itoa(count2) + " waiting item(s)."})
		for i, v := range Model {
			if !v.Completed {
				index := i
				widgets = append(widgets, &goey.Checkbox{Text: v.Text, Value: v.Completed,
					OnChange: func(newValue bool) {
						Model[index].Completed = newValue
						update()
					}})
			}
		}
	}
	if count1 > 0 {
		widgets = append(widgets, &goey.HR{})
		widgets = append(widgets, &goey.Label{Text: "There are " + strconv.Itoa(count1) + " completed item(s)."})
		for i, v := range Model {
			if v.Completed {
				index := i
				widgets = append(widgets, &goey.Checkbox{Text: v.Text, Value: v.Completed,
					OnChange: func(newValue bool) {
						Model[index].Completed = newValue
						update()
					}})
			}
		}
	}

	return widgets
}

func onNewTodoItem(value string) {
	Model = append(Model, TodoItem{Text: value})
	mainWindow.SetChildren(render())
}
