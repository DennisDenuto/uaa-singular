// +build js

package main

import (
	"github.com/DennisDenuto/go-uaa-singular/oidc"
	"net/http"
	"fmt"
	"github.com/DennisDenuto/go-uaa-singular/oidc/session_management/iframes/oidc_party"
	"bytes"
	"github.com/gopherjs/jquery"
	"github.com/gopherjs/todomvc/utils"
	"html/template"
)

func main1() {
	authURL, err := oidc.NewSilentAuthURL("http://localhost:8080/uaa/", "http://localhost:9090", "app", "some-state")
	if err != nil {
		panic(err)
	}

	println(authURL)

	silentAuthResponse, err := http.Get(authURL)
	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("%v", silentAuthResponse))
}

func main2() {
	authURL, err := oidc.NewSilentAuthURL("http://localhost:8080/uaa/", "http://localhost:9090", "app", "some-state")
	if err != nil {
		panic(err)
	}

	println(authURL)
	client := oidc_party.OpenIDProviderClient{
		SilentAuthURL: authURL,
		RedirectURL:   "http://localhost:9090",
		AuthResponse:  nil,
	}
	_, err = client.CheckSession("app", "some-state")
	if err != nil {
		panic(err)
	}
}

var jQuery = jquery.NewJQuery //for convenience

const (
	KEY        = "TodoMVC4GopherJS"
	ENTER_KEY  = 13
	ESCAPE_KEY = 27
)

func main() {
	app := NewApp()
	app.bindEvents()
	app.initRouter()
	app.render()
}

type ToDo struct {
	Id        string
	Text      string
	Completed bool
}
type App struct {
	todos           []ToDo
	todoTmpl        *template.Template
	footerTmpl      *template.Template
	todoAppJQuery   jquery.JQuery
	headerJQuery    jquery.JQuery
	mainJQuery      jquery.JQuery
	footerJQuery    jquery.JQuery
	newTodoJQuery   jquery.JQuery
	toggleAllJQuery jquery.JQuery
	todoListJQuery  jquery.JQuery
	countJQuery     jquery.JQuery
	clearBtnJQuery  jquery.JQuery
	filter          string
}

func NewApp() *App {
	somethingToDo := make([]ToDo, 0)
	utils.Retrieve(KEY, &somethingToDo)

	todoHtml := jQuery("#todo-template").Html()
	todoTmpl := template.Must(template.New("todo").Parse(todoHtml))

	footerHtml := jQuery("#footer-template").Html()
	footerTmpl := template.Must(template.New("footer").Parse(footerHtml))

	todoAppJQuery := jQuery("#todoapp")
	headerJQuery := todoAppJQuery.Find("#header")
	mainJQuery := todoAppJQuery.Find("#main")
	footerJQuery := todoAppJQuery.Find("#footer")
	newTodoJQuery := headerJQuery.Find("#new-todo")
	toggleAllJQuery := mainJQuery.Find("#toggle-all")
	todoListJQuery := mainJQuery.Find("#todo-list")
	countJQuery := footerJQuery.Find("#todo-count")
	clearBtnJQuery := footerJQuery.Find("#clear-completed")
	filter := "all"

	return &App{somethingToDo, todoTmpl, footerTmpl, todoAppJQuery, headerJQuery, mainJQuery, footerJQuery, newTodoJQuery, toggleAllJQuery, todoListJQuery, countJQuery, clearBtnJQuery, filter}
}
func (a *App) bindEvents() {

	a.newTodoJQuery.On(jquery.KEYUP, a.create)
	a.toggleAllJQuery.On(jquery.CHANGE, a.toggleAll)
	a.footerJQuery.On(jquery.CLICK, "#clear-completed", a.destroyCompleted)
	a.todoListJQuery.On(jquery.CHANGE, ".toggle", a.toggle)
	a.todoListJQuery.On(jquery.DBLCLICK, "label", a.edit)
	a.todoListJQuery.On(jquery.KEYUP, ".edit", a.blurOnEnter)
	a.todoListJQuery.On(jquery.FOCUSOUT, ".edit", a.update)
	a.todoListJQuery.On(jquery.CLICK, ".destroy", a.destroy)
}
func (a *App) initRouter() {
	router := utils.NewRouter()
	router.On("/:filter", func(filter string) {
		a.filter = filter
		a.render()
	})
	router.Init("/all")
}
func (a *App) render() {
	todos := a.getFilteredTodos()

	var b bytes.Buffer
	a.todoTmpl.Execute(&b, todos)
	strtodoTmpl := b.String()

	a.todoListJQuery.SetHtml(strtodoTmpl)
	a.mainJQuery.Toggle(len(a.todos) > 0)
	a.toggleAllJQuery.SetProp("checked", len(a.getActiveTodos()) != 0)
	a.renderfooter()
	a.newTodoJQuery.Focus()
	utils.Store(KEY, a.todos)
}
func (a *App) renderfooter() {
	activeTodoCount := len(a.getActiveTodos())
	activeTodoWord := utils.Pluralize(activeTodoCount, "item")
	completedTodos := len(a.todos) - activeTodoCount
	filter := a.filter
	footerData := struct {
		ActiveTodoCount int
		ActiveTodoWord  string
		CompletedTodos  int
		Filter          string
	}{
		activeTodoCount, activeTodoWord, completedTodos, filter,
	}
	var b bytes.Buffer
	a.footerTmpl.Execute(&b, footerData)
	footerJQueryStr := b.String()

	a.footerJQuery.Toggle(len(a.todos) > 0).SetHtml(footerJQueryStr)
}
func (a *App) toggleAll(e jquery.Event) {
	checked := !a.toggleAllJQuery.Prop("checked").(bool)
	for idx := range a.todos {
		a.todos[idx].Completed = checked
	}
	a.render()
}
func (a *App) getActiveTodos() []ToDo {
	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if !val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	return todosTmp
}
func (a *App) getCompletedTodos() []ToDo {
	todosTmp := make([]ToDo, 0)
	for _, val := range a.todos {
		if val.Completed {
			todosTmp = append(todosTmp, val)
		}
	}
	return todosTmp
}
func (a *App) getFilteredTodos() []ToDo {
	switch a.filter {
	case "active":
		return a.getActiveTodos()
	case "completed":
		return a.getCompletedTodos()
	default:
		return a.todos
	}
}
func (a *App) destroyCompleted(e jquery.Event) {
	a.todos = a.getActiveTodos()
	a.filter = "all"
	a.render()
}
func (a *App) indexFromEl(e jquery.Event) int {
	id := jQuery(e.Target).Closest("li").Data("id")
	for idx, val := range a.todos {
		if val.Id == id {
			return idx
		}
	}
	return -1
}
func (a *App) create(e jquery.Event) {
	val := jquery.Trim(a.newTodoJQuery.Val())
	if len(val) == 0 || e.KeyCode != ENTER_KEY {
		return
	}
	newToDo := ToDo{Id: utils.Uuid(), Text: val, Completed: false}
	a.todos = append(a.todos, newToDo)
	a.newTodoJQuery.SetVal("")
	a.render()
}
func (a *App) toggle(e jquery.Event) {
	idx := a.indexFromEl(e)
	a.todos[idx].Completed = !a.todos[idx].Completed
	a.render()
}
func (a *App) edit(e jquery.Event) {
	input := jQuery(e.Target).Closest("li").AddClass("editing").Find(".edit")
	input.SetVal(input.Val()).Focus()
}
func (a *App) blurOnEnter(e jquery.Event) {
	switch e.KeyCode {
	case ENTER_KEY:
		jQuery(e.Target).Blur()
	case ESCAPE_KEY:
		jQuery(e.Target).SetData("abort", "true").Blur()
	}
}
func (a *App) update(e jquery.Event) {

	thisJQuery := jQuery(e.Target)
	val := jquery.Trim(thisJQuery.Val())
	if thisJQuery.Data("abort") == "true" {
		thisJQuery.SetData("abort", "false")
		a.render()
		return
	}
	idx := a.indexFromEl(e)
	if len(val) > 0 {
		a.todos[idx].Text = val
	} else {
		a.todos = append(a.todos[:idx], a.todos[idx+1:]...)
	}
	a.render()
}
func (a *App) destroy(e jquery.Event) {
	idx := a.indexFromEl(e)
	a.todos = append(a.todos[:idx], a.todos[idx+1:]...)
	a.render()
}