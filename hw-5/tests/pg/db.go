package pg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/pg"
	ordersch "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/order"
	wrapsch "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
)

const testEnvPath = "../../../test.env"

type TDB struct {
	ConnPool *pgxpool.Pool
}

func NewFromEnv() *TDB {
	if err := godotenv.Load(testEnvPath); err != nil {
		panic(err)
	}

	connPool, err := pg.NewPoolConn(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return &TDB{ConnPool: connPool}
}

func (db *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	db.truncateTable(context.Background(), tableName...)
}

func (db *TDB) TearDown(t *testing.T) {
	t.Helper()
	db.ConnPool.Reset()
}

func (db *TDB) truncateTable(ctx context.Context, tableName ...string) {
	q := fmt.Sprintf("TRUNCATE %s", strings.Join(tableName, ","))
	if _, err := db.ConnPool.Exec(ctx, q); err != nil {
		panic(err)
	}
}

func (db *TDB) FillOrders(records []order.Order) {
	for _, r := range records {
		_, err := db.ConnPool.Exec(context.Background(),
			`INSERT INTO orders(id,client_id,stored_until,weight,cost,status,status_changed_at,wrap_type,hash) 
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`,
			r.OrderID, r.ClientID, r.StoredUntil, r.Weight, r.Cost, r.Status, r.StatusChanged, r.WrapType, r.Hash,
		)
		if err != nil {
			panic(err)
		}
	}
}

func (db *TDB) GetAllOrders() []ordersch.OrderSchema {
	rows, err := db.ConnPool.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		panic(err)
	}

	var records []ordersch.OrderSchema

	if err = pgxscan.ScanAll(&records, rows); err != nil {
		panic(err)
	}

	return records
}

func (db *TDB) FillWraps(records []wrap.Wrap) {
	for _, r := range records {
		_, err := db.ConnPool.Exec(context.Background(),
			`INSERT INTO wrap(name, weight, cost) VALUES ($1,$2,$3);`,
			r.Name, r.Weight, r.Cost,
		)
		if err != nil {
			panic(err)
		}
	}
}

func (db *TDB) GetAllWraps() []wrapsch.WrapSchema {
	rows, err := db.ConnPool.Query(context.Background(), "SELECT * FROM wrap")
	if err != nil {
		panic(err)
	}

	var records []wrapsch.WrapSchema

	if err = pgxscan.ScanAll(&records, rows); err != nil {
		panic(err)
	}

	return records
}
