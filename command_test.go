package noredis

// import (
// 	"testing"
// )

// func Test_SET(t *testing.T) {

// }

// func Test_SET_options(t *testing.T) {

// 	args := []string{"ex", "1000", "xx"}
// 	opt := new(Setoptions)
// 	err := opt.fromArgs(args)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if opt.expireInSecond != true {
// 		t.Error(opt.expireInSecond)
// 	}
// 	if opt.expireValue != "1000" {
// 		t.Error(opt.expireValue)
// 	}
// 	if opt.keyShouldExist != true {
// 		t.Error(opt.keyShouldExist)
// 	}

// 	if opt.keyShouldNotExist == true {
// 		t.Error("should false")
// 	}
// }

// type MockQuerys struct {
// 	cmd  string
// 	args []string
// }

// func (mq *MockQuerys) Cmd() string {
// 	return mq.cmd
// }

// func (mq *MockQuerys) Args() []string {
// 	return mq.args
// }

// func Test_Set_Command(t *testing.T) {
// 	q := MockQuerys{
// 		args: []string{"jhon", "doe", "ex", "1000"},
// 	}
// 	db := NewDB()

// 	err := setCommand(q.args, db)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
