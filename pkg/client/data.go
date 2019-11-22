package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
)

// SavePing saves the ping results into DB
func SavePing(ctx context.Context, dbName, id, region string, sent time.Time) error {

	dbClient, err := spanner.NewClient(ctx, dbName)
	if err != nil {
		return fmt.Errorf("error while creating spanner db: %v", err)
	}

	_, err = dbClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		stmt := spanner.Statement{
			SQL: `INSERT pings (id, region, sent) VALUES (@id, @region, @sent)`,
			Params: map[string]interface{}{
				"id":     id,
				"region": region,
				"sent":   sent,
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

// CompletePing updates previously saved ping in the BD
func CompletePing(ctx context.Context, dbName, id string, completed time.Time) error {

	dbClient, err := spanner.NewClient(ctx, dbName)
	if err != nil {
		return fmt.Errorf("error while creating spanner db: %v", err)
	}

	_, err = dbClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		stmt := spanner.Statement{
			SQL: `UPDATE pings set completed = @completed where id = @id`,
			Params: map[string]interface{}{
				"id":        id,
				"completed": completed,
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
