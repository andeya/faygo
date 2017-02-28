package main

import (
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/xorm"
	"mime/multipart"
)

// MultiUsage multi usage handler
type MultiUsage struct {
	Id         int64                 `param:"-" xorm:"not null pk autoincr INT(11)"`
	Name       string                `param:"<in:formData><required>" xorm:"not null VARCHAR(25)"`
	Age        uint8                 `param:"<in:formData><range: 1:100>" xorm:"INT(3)"`
	AvatarFile *multipart.FileHeader `param:"<in:formData><maxmb:30><name:avatar><desc:(not more than 30 MB)>" xorm:"-"`
	Avatar     string                `param:"-" xorm:"not null VARCHAR(250)"`
}

func init() {
	// Create the multi_usage form
	xorm.MustDB().Sync2(new(MultiUsage))
}

// TableName gives xorm the table name
func (m *MultiUsage) TableName() string {
	return "multi_usage"
}

// AddUser add user to database
func (m *MultiUsage) AddUser() error {
	_, err := xorm.MustDB().Insert(m)
	return err
}

// Serve impletes Handler
func (m *MultiUsage) Serve(ctx *faygo.Context) error {
	info, err := ctx.SaveFile("avatar", false, m.Name)
	if err != nil {
		return ctx.String(412, err.Error())
	}
	m.Avatar = info.Url
	err = m.AddUser()
	if err != nil {
		return ctx.String(503, "error:%v", err)
	}
	return ctx.String(200, "Success added user:\nname: %s\nage: %d\navatar: %s",
		m.Name,
		m.Age,
		ctx.Site()+m.Avatar,
	)
}

// Doc returns the API's note.
func (m *MultiUsage) Doc() faygo.Doc {
	return faygo.Doc{
		Note: "struct handler's multi usage",
	}
}

func main() {
	// new app
	app := faygo.New("myapp", "1.0")
	// Register the route, and set the default value for apidoc
	app.POST("/multi", &MultiUsage{
		Name: "henrylee",
		Age:  30,
	})
	// Start the service
	faygo.Run()

	// PS: By visiting /apidoc to test.
}
