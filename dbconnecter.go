package dbconnecter

import (
	"database/sql"
	"fmt"
	"sync/atomic"
)

type Connection struct {
	*sql.DB
}

type MultipleDBConnecter struct {
	activeConnection atomic.Value
	connections      []Connection
}

func NewMultipleDBConnecter(connections ...*sql.DB) *MultipleDBConnecter {
	conns := make([]Connection, len(connections))
	for i := range conns {
		conns[i] = Connection{DB: connections[i]}
	}
	return &MultipleDBConnecter{
		connections: conns,
	}
}

func (c *MultipleDBConnecter) Connection() (Connection, error) {
	if err := c.setActiveConnection(); err != nil {
		return Connection{}, err
	}

	value, ok := c.activeConnection.Load().(Connection)
	if !ok {
		return Connection{}, fmt.Errorf("ошибка приведения типа")
	}

	return value, nil
}

func (c *MultipleDBConnecter) setActiveConnection() error {
	activeConnections := make(chan Connection)
	fails := make(chan struct{})

	for _, connection := range c.connections {
		go ping(connection, activeConnections, fails)
	}

	i := 0
	for {
		select {
		case <-fails:
			i++
			if i >= len(c.connections) {
				return fmt.Errorf("нет доступных соединений")
			}
		case firstConn := <-activeConnections:
			c.activeConnection.Store(firstConn)
			return nil
		}
	}
}

func ping(connection Connection, successConnections chan<- Connection, fails chan<- struct{}) {
	if err := connection.Ping(); err != nil {
		fails <- struct{}{}
	} else {
		successConnections <- connection
	}
}
