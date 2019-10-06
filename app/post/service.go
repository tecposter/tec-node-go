package post

import (
	"database/sql"
	"github.com/tecposter/tec-node-go/lib/dto"
	"github.com/tecposter/tec-node-go/lib/uuid"
)

/*
const (
	cmdCreate = "post.create"
	cmdEdit   = "post.edit"
	cmdFetch  = "post.fetch"
	cmdCommit = "post.commit"
	cmdList   = "post.list"
	cmdSearch = "post.search"
)
*/

type service struct {
	db          *sql.DB
	commitRepo  *commitRepository
	contentRepo *contentRepository
	postRepo    *postRepository
}

func (s *service) create() (dto.ID, error) {
	id := uuid.NewID()
	created := time.Now().UnixNano()

	err = s.postRepo.create(id, created)
	if err != nil {
		return nil, err
	}
	return id, nil
}

/*

tx, err := db.Begin()
	if err != nil {
				log.Fatal(err)
					}
						stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
							if err != nil {
										log.Fatal(err)
											}
												defer stmt.Close()
												for i := 0; i < 100; i++ {
													_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
													if err != nil {
														log.Fatal(err)
													}
												}
																										tx.Commit()
package main

import (
		"context"
			"database/sql"
				"log"
			)

			var (
					ctx context.Context
						db  *sql.DB
					)

					func main() {
						tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
						if err != nil {
							log.Fatal(err)
						}
						id := 53
						_, err = tx.ExecContext(ctx, "UPDATE drivers SET status = ? WHERE id = ?;", "assigned", id)
						if err != nil {
							if rollbackErr := tx.Rollback(); rollbackErr != nil {
								log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
							}
							log.Fatal(err)
						}
						_, err = tx.ExecContext(ctx, "UPDATE pickups SET driver_id = $1;", id)
						if err != nil {
							if rollbackErr := tx.Rollback(); rollbackErr != nil {
								log.Fatalf("update failed: %v, unable to back: %v", err, rollbackErr)
							}
							log.Fatal(err)
						}
						if err := tx.Commit(); err != nil {
							log.Fatal(err)
						}
					}

*/
