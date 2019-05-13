package post

import (
	"context"
	"errors"
	"sync"
	"time"
	"regexp"
	"github.com/google/uuid"
	"github.com/btcsuite/btcutil/base58"
	"github.com/tecposter/tec-server-go/third/ipfs"
	"github.com/tecposter/tec-server-go/internal/wicc"
	//"github.com/tecposter/tec-server-go/third/mapstructure"
)

type postT struct {
	Pid string `json:"pid"`
	Uid string `json:"uid"`
	Type string `json:"type"`
	Cid string `json:"cid"`
	Prev string `json:"prev"`
}

func Create(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	pid, err := uuid.NewUUID()
    if err !=nil {
		return nil, err
        // handle error
    }

	b, err := pid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	s := base58.Encode(b)
	SaveDraft(ctx, map[string]interface{}{"pid": s, "type": "markdown", "content": "# title \n"})
	return map[string]interface{}{"pid": s}, nil
}

func Edit(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	cid := input["cid"].(string)

	post, e := fetchPost(cid)
	if e != nil {
		return nil, e
	}

	draft := getDraftSt().fetch(post["pid"].(string))

	if draft == nil || draft.Prev != cid {
		post["prev"] = post["cid"]
		delete(post, "cid")
		SaveDraft(ctx, post)
	}

	return map[string]interface{}{"pid": post["pid"]}, nil
}


func List(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	uid, err := getUid(ctx)
	if err != nil {
		return nil, err
	}

	var (
		currTxnCid string
		txn *txnT
		e error
	)
	currTxnCid, _ = getUserTxnSt().get(uid)
	//list := make([]map[string]interface{}, 0)
	m := make(map[string]bool)
	list := make([]interface{}, 0)

	for currTxnCid != "" {
		txn, e = getTxn(currTxnCid)
		if e != nil {
			return nil, e
		}
		post, e := fetchPost(txn.Cid)
		if e != nil {
			return nil, e
		}
		pid := post["pid"].(string)
		if _, ok := m[pid]; !ok {
			m[pid] = true
			list = append(list, post)
		}
		currTxnCid = txn.Prev
	}

	/*
	list := make([]interface{}, len(m))
	idx := 0
	for _, val := range m {
		list[idx] = val
		idx++
	}
	*/

	return map[string]interface{}{"list": list}, nil
	//return map[string]interface{}{"list": m}, nil
}

func fetchPost(postCid string) (map[string]interface{}, error) {
	var post postT
	err := ipfs.RecvByCid(postCid, &post)
	if err != nil {
		return nil, err
	}

	content, err := ipfs.FetchStr(post.Cid)
	if err != nil {
		return nil, err
	}

	return map[string]interface{} {
		"pid": post.Pid,
		"cid": postCid,
		"uid": post.Uid,
		"type": post.Type,
		"title": extractMdTitle(content),
		"content": content}, nil
}

func Fetch(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	postCid := input["cid"].(string)
	return fetchPost(postCid)
}

func Commit(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	uid, err := getUid(ctx)
	if err != nil {
		return nil, err
	}

	pid := input["pid"].(string)
	typ := input["type"].(string)
	content := input["content"].(string)

	cid, err := ipfs.AddStr(content)
	if err != nil {
		return nil, err
	}

	postCid, err := ipfs.Add(map[string]string{
		"uid": uid,
		"pid": pid,
		"type": typ,
		"cid": cid})
	if err != nil {
		return nil, err
	}

	txn, err := addTxn(uid, postCid)
	if err != nil {
		return nil, err
	}

	RemoveDraft(ctx, map[string]interface{}{"pid": pid})

	wicc.Send("cid", txn["cid"].(string))
	return txn, nil
}

func getUid(ctx context.Context) (string, error) {
	if uid, ok := ctx.Value("uid").(string); ok {
		return uid, nil
	}

	return "", errors.New("cannot find uid")
}

type userTxnStT struct {
	sm sync.Map
}

var once sync.Once
var userTxnSt *userTxnStT

func getUserTxnSt() *userTxnStT {
	once.Do(func () {
		userTxnSt = &userTxnStT{}
	})
	return userTxnSt
}

func (st *userTxnStT) set(uid, cid string) {
	st.sm.Store(uid, cid)
}

func (st *userTxnStT) get(uid string) (string, bool) {
	if val, ok := st.sm.Load(uid); ok {
		return val.(string), true
	}
	return "", false
}

func addTxn(uid, postCid string) (map[string]interface{}, error) {
	lastTxnCid, _ := getUserTxnSt().get(uid)

	m := map[string]interface{} {
		"cid": postCid,
		"utc": time.Now().UTC(),
		"prev": lastTxnCid}

	currentTxnCid, err := ipfs.Add(m)
	if err != nil {
		return nil, err
	}

	getUserTxnSt().set(uid, currentTxnCid)
	return m, nil
}

type txnT struct {
	Cid string
	Utc time.Time
	Prev string
}

func getTxn(cid string) (*txnT, error) {
	var txn txnT
	err := ipfs.RecvByCid(cid, &txn)
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func extractMdTitle(content string) string {
    //$matched = preg_match('/# ([^#\n]+)/', $content, $matches);
    re := regexp.MustCompile(`# ([^#\n]+)\n`)
    founds := re.FindStringSubmatch(content)

    //fmt.Println(founds)
    if len(founds) >= 2 {
        return founds[1]
    } else {
		l := len(content)
		if l > 100 {
			return content[0:100]
		} else {
			return content
		}
    }
}
