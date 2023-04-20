package postgresdriver

import (
	"database/sql"
	"database/sql/driver"
	"net"
	"time"

	// PQ import is required
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Open(s string) (_ driver.Conn, err error) {
	return pq.DialOpen(self, s)
}

func (self *ViaSSHDialer) Dial(network, address string) (net.Conn, error) {
	return self.client.Dial(network, address)
}

func (self *ViaSSHDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return self.client.Dial(network, address)
}

// The PostgresDriver struct satisfies the Driver interface which defines all database driver methods
type PostgresDriver struct {
	*Queries
	db *sql.DB
}

/* NewPostgresDriver returns PostgresDriver instance from Postgres connection string */
func NewPostgresDriver(connectionString string) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	driver := &PostgresDriver{
		Queries: New(db),
		db:      db,
	}

	return driver, nil
}

/* NewPostgresDriverFromDBInstance returns PostgresDriver instance from sdl.DB instance */
// mostly used for mocking tests
func NewPostgresDriverFromDBInstance(db *sql.DB) *PostgresDriver {
	driver := &PostgresDriver{
		Queries: New(db),
	}

	return driver
}

/* NewPostgresDriver returns PostgresDriver instance from Postgres connection string */
func NewPostgresDriverWithSSH(connectionString string, sshcon *ssh.Client) (*PostgresDriver, error) {
	// Now we register the ViaSSHDialer with the ssh connection as a parameter
	sql.Register("postgres+ssh", &ViaSSHDialer{sshcon})

	db, err := sql.Open("postgres+ssh", connectionString)
	if err != nil {
		return nil, err
	}

	driver := &PostgresDriver{
		Queries: New(db),
		db:      db,
	}

	return driver, nil
}
