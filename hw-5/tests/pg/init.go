package pg

var (
	DB *TDB
)

func init() {
	DB = NewFromEnv()
}
