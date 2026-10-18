package migrator

import "context"

type TransactionRunner func(ctx context.Context, fn func(ctx context.Context) error) error
