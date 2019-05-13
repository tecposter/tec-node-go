package ws

import (
	"log"
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"fmt"
	"errors"
	"github.com/tecposter/tec-server-go/internal/post"
	"github.com/tecposter/tec-server-go/internal/user"
)

const (
	//postPublishCmd = "post.publish"
	userRegCmd = "user.reg"
	userLoginCmd = "user.login"
	postCreateCmd = "post.create"
	postEditCmd = "post.edit"
	postListCmd = "post.list"
	postFetchCmd = "post.fetch"
	postCommitCmd = "post.commit"
	postListDraftCmd = "post.listDraft"
	postFetchDraftCmd = "post.fetchDraft"
	postSaveDraftCmd = "post.saveDraft"
	postRemoveDraftCmd = "post.removeDraft"
)

type msgT struct {
	Cmd string `json:"cmd"`
	Token string `json:"token"`
	Data map[string]interface{} `json:"data"`
}

var isDebug = true

func filter(ctx context.Context, msg msgT) (context.Context, error) {
	if isDebug {
		return context.WithValue(ctx, "uid", "test-abc"), nil
	}

	if msg.Cmd == userRegCmd || msg.Cmd == userLoginCmd {
		return ctx, nil
	}

	if msg.Token == "" {
		return ctx, errors.New("token cannot be empty")
	}

	uid := user.GetUid(msg.Token)
	if uid == "" {
		return ctx, errors.New(fmt.Sprintf("Cannot find uid from token: %s", msg.Token))
	}

	//log.Println("uid:", uid)

	return context.WithValue(ctx, "uid", uid), nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, text, err := c.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			break
		}

		var msg msgT;
		json.Unmarshal(text, &msg)
		fmt.Printf("recv: %s, msg: %+v\n", text, msg)

		inner, err := handleMsg(msg)
		var m map[string]interface{}
		if err == nil {
			m = map[string]interface{} {
				"cmd": msg.Cmd,
				"status": "ok",
				"data": inner}
		} else {
			m = map[string]interface{} {
				"cmd": msg.Cmd,
				"status": "err",
				"err": err.Error()}
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Println("json.Marshal:", err)
			break
		}

		err = c.WriteMessage(mt, b)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type handleFunc func(context.Context, map[string]interface{}) (map[string]interface{}, error)
func handleMsg(msg msgT) (map[string]interface{}, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, err := filter(ctx, msg)
	if  err != nil {
		return nil, err
	}

	var h handleFunc
	switch msg.Cmd {
		case userLoginCmd:
			h = user.Login
		case postCreateCmd:
			h = post.Create
		case postEditCmd:
			h = post.Edit
		case postListCmd:
			h = post.List
		case postFetchCmd:
			h = post.Fetch
		case postCommitCmd:
			h = post.Commit
		case postListDraftCmd:
			h = post.ListDraft
		case postFetchDraftCmd:
			h = post.FetchDraft
		case postSaveDraftCmd:
			h = post.SaveDraft
		case postRemoveDraftCmd:
			h = post.RemoveDraft
		default:
			h = nil
	}

	if h == nil {
		return nil, errors.New("Unkown Message")
	}

	return h(ctx, msg.Data)
}

/*
type unknownMsgError struct {}
func (*unknownMsgError) Error() string {
	return fmt.Sprintf("Unkown Message")
}
*/
