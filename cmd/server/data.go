package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
)

func save(ctx context.Context, dbPath, id, target, source string, sent, completed time.Time, duration int64) error {

	dbClient, err := spanner.NewClient(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("error while creating spanner db: %v", err)
	}

	_, err = dbClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		stmt := spanner.Statement{

			SQL: `INSERT pings (id, target, source, sent, completed, duration) VALUES (@id, @target, @source, @sent, @completed, @duration)`,
			Params: map[string]interface{}{
				"id":        id,
				"target":    target,
				"source":    source,
				"sent":      sent,
				"completed": completed,
				"duration":  duration,
			},
		}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return fmt.Errorf("error while inserting ping: %v", err)
		}

		if rowCount != 1 {
			return errors.New("expected ping insert to impact only 1 row")
		}

		return nil
	})
	return err

}
