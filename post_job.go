package main

import (
	"github.com/wayt/happyngine"
)

type postJobAction struct {
	happyngine.Action
}

func newPostJobAction(context *happyngine.Context) happyngine.ActionInterface {

	// Init
	this := &postJobAction{}
	this.Context = context

	this.Form = happyngine.NewForm(context,
		happyngine.NewFormElement("name", "invalid_job"))

	return this
}

func (this *postJobAction) Run() {

	// Get assotiated playbook
	job, err := GetJob(this.Form.Elem("name").FormValue())
	if err != nil {
		panic(err)
	}

	if job == nil {
		this.AddError(404, "Unknown job")
		return
	}

	if err := job.Run(); err != nil {
		panic(err)
	}

	data, err := job.Logs()
	if err != nil {
		panic(err)
	}

	this.Send(200, string(data))
}
