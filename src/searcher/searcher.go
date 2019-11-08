package searcher

import (
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	searcher = riot.Engine{}
)

func Init() {
	searcher.Init(types.EngineOpts{
		NotUseGse: true,
	})

	text := "Google Is Experimenting With Virtual Reality Advertising"
	text1 := `Google accidentally pushed Bluetooth update for Home
			speaker early`
	text2 := `Google is testing another Search results layout with 
					rounded cards, new colors, and the 4 mysterious colored dots again`

	// Add the document to the index, docId starts at 1
	searcher.Index("1", types.DocData{Content: text})
	searcher.Index("2", types.DocData{Content: text1}, false)
	searcher.IndexDoc("3", types.DocData{Content: text2}, true)

	// Wait for the index to refresh
	searcher.Flush()
	// engine.FlushIndex()

	// The search output format is found in the types.SearchResp structure
	log.Print(searcher.Search(types.SearchReq{Text: "google testing"}))
}

func Close() {
	searcher.Close()
}

func Index(id, content string) {
	log.Println("Index: ", id, content)
	searcher.Index(id, types.DocData{Content: content})
	searcher.Flush()
}

func Search(query string) types.SearchResp {
	log.Println("Search: ", query)
	return searcher.Search(types.SearchReq{Text: query})
}

/*
import (
	"fmt"
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}

	text  = "Google Is Experimenting With Virtual Reality Advertising"
	text1 = `Google accidentally pushed Bluetooth update for Home
	speaker early`
	text2 = `Google is testing another Search results layout with
	rounded cards, new colors, and the 4 mysterious colored dots again`

	opts = types.EngineOpts{
		Using: 1,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		UseStore: true,
		// StoreFolder: path,
		StoreEngine: "bg", // bg: badger, lbd: leveldb, bolt: bolt
		// GseDict: "../../data/dict/dictionary.txt",
		GseDict:       "../../testdata/test_dict.txt",
		StopTokenFile: "../../data/dict/stop_tokens.txt",
	}
)

func initEngine() {
	// gob.Register(MyAttriStruct{})

	// var path = "./riot-index"
	searcher.Init(opts)
	defer searcher.Close()
	// os.MkdirAll(path, 0777)

	// Add the document to the index, docId starts at 1
	searcher.Index("1", types.DocData{Content: text})
	searcher.Index("2", types.DocData{Content: text1})
	searcher.Index("3", types.DocData{Content: text2})
	searcher.Index("5", types.DocData{Content: text2})

	searcher.RemoveDoc("5")

	// Wait for the index to refresh
	searcher.Flush()

	log.Println("Created index number: ", searcher.NumDocsIndexed())
}

func restoreIndex() {
	// var path = "./riot-index"
	searcher.Init(opts)
	defer searcher.Close()
	// os.MkdirAll(path, 0777)

	// Wait for the index to refresh
	searcher.Flush()

	log.Println("recover index number: ", searcher.NumDocsIndexed())
}

func main() {
	initEngine()
	// restoreIndex()

	sea := searcher.Search(types.SearchReq{
		Text: "google testing",
		RankOpts: &types.RankOpts{
			OutputOffset: 0,
			MaxOutputs:   100,
		}})

	fmt.Println("search response: ", sea, "; docs = ", sea.Docs)

	// os.RemoveAll("riot-index")
}
*/
