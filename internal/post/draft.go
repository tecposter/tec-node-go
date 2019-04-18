package post

import (
	//"log"
	"sync"
	"context"
	"errors"
	"fmt"
	"github.com/tecposter/tec-server-go/third/mapstructure"
)

type userDraftStT struct {
	sm sync.Map // [uid]*draftStT
}

type draftStT struct {
	sm sync.Map // [pid]*draftT
}

type draftT struct {
	Pid string
	Uid string
	Type string
	Content string
}

func ListDraft(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	uid, err := getUid(ctx)
	if err != nil {
		return nil, err
	}

	st := getUserDraftSt().getDraftSt(uid)
	list := make([]*draftT, 0)

	st.rg(func (_, val interface{}) bool {
		list = append(list, val.(*draftT))
		return true
	})

	return map[string]interface{}{"list": list}, nil

}

func FetchDraft(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	pid := input["pid"]

	draft := getDraftSt().fetch(pid.(string))
	if draft == nil {
		return nil, errors.New(fmt.Sprintf("Cannot find draft with pid: %s", pid))
	}

	return map[string]interface{}{
		"pid": draft.Pid,
		"uid": draft.Uid,
		"type": draft.Type,
		"content": draft.Content}, nil
}

func SaveDraft(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	uid, err := getUid(ctx)
	if err != nil {
		return nil, err
	}

	var draft draftT
	mapstructure.Decode(input, &draft)
	draft.Uid = uid;

	//log.Println("draft: ", draft)

	getUserDraftSt().saveDraft(&draft)

	return map[string]interface{}{"uid": draft.Uid, "pid": draft.Pid, "type": draft.Type, "content": draft.Content}, nil
}


// private funs
var once1 sync.Once
var once2 sync.Once
var userDraftSt *userDraftStT
var draftSt *draftStT

func getUserDraftSt() *userDraftStT {
	once1.Do(func() {
		userDraftSt = &userDraftStT{}
	})
	return userDraftSt
}

func getDraftSt() *draftStT {
	once2.Do(func() {
		draftSt = &draftStT{}
	})
	return draftSt
}

func (st *userDraftStT) saveDraft(draft *draftT) {
	st.getDraftSt(draft.Uid).save(draft)
	getDraftSt().save(draft)
}

func (st *userDraftStT) getDraftSt(uid string) *draftStT {
	if val, ok := st.sm.Load(uid); ok {
		return val.(*draftStT)
	}

	newSt := &draftStT{}
	st.sm.Store(uid, newSt)

	return newSt
}

func (st *draftStT) save(draft *draftT) {
	//log.Println("draft.Pid: ", draft.Pid)
	st.sm.Store(draft.Pid, draft)
}

func (st *draftStT) rg(f func(key, val interface{}) bool) {
	st.sm.Range(f)
}

func (st *draftStT) fetch(pid string) *draftT {
	if val, ok := st.sm.Load(pid); ok {
		return val.(*draftT)
	}
	return nil
}

/*var instance *singleton
var once sync.Once

func GetInstance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}*/
