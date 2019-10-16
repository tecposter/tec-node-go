package draft

import (
	"database/sql"

	"github.com/tecposter/tec-node-go/lib/ws"
)

type Controller struct {
	db *sql.DB
}

func NewCtrl(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (ctrl *Controller) Handle(req ws.IRequest) (ws.IResponse, error) {
	return req, nil
}
