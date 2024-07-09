package admin

import (
	"errors"
	"net"

	"github.com/Clever/gearadmin"
)

var ErrNotConnected = errors.New("not connected")

type GearAdminClient struct {
	addr string
	conn net.Conn
}

func NewGearAdminClient(addr string) *GearAdminClient {
	return &GearAdminClient{addr: addr}
}

func (g *GearAdminClient) Connect() error {
	var err error
	g.conn, err = net.Dial("tcp", g.addr)
	return err
}

func mapTo[T1 any, T2 any](m []T1, trans func(T1) T2) []T2 {
	r := make([]T2, 0, len(m))
	for _, v := range m {
		r = append(r, trans(v))
	}
	return r
}

func (g *GearAdminClient) Status() ([]Status, error) {
	if g.conn == nil {
		return nil, ErrNotConnected
	}
	client := gearadmin.NewGearmanAdmin(g.conn)
	status, err := client.Status()
	if err != nil {
		return nil, err
	}
	r := mapTo(status, func(m gearadmin.Status) Status {
		return Status{
			Function:         m.Function,
			AvailableWorkers: m.AvailableWorkers,
			Running:          m.Running,
			Total:            m.Total,
		}
	})
	return r, nil
}

func (g *GearAdminClient) Workers() ([]Worker, error) {
	if g.conn == nil {
		return nil, ErrNotConnected
	}
	client := gearadmin.NewGearmanAdmin(g.conn)
	workers, err := client.Workers()
	if err != nil {
		return nil, err
	}
	r := mapTo(workers, func(m gearadmin.Worker) Worker {
		return Worker{
			Fd:        m.Fd,
			IPAddress: m.IPAddress,
			ClientID:  m.ClientID,
			Functions: m.Functions,
		}
	})
	return r, nil

}

func (g *GearAdminClient) Close() error {
	if g.conn != nil {
		err := g.conn.Close()
		g.conn = nil
		return err
	}
	return nil
}

type Status struct {
	Function         string `json:"name"`
	Total            int    `json:"total"`
	Running          int    `json:"running"`
	AvailableWorkers int    `json:"available_workers"`
}

// Worker represents a worker connected to gearman as returned by the "workers" command.
type Worker struct {
	Fd        string   `json:"fd"`
	IPAddress string   `json:"ip_address"`
	ClientID  string   `json:"client_id"`
	Functions []string `json:"functions"`
}

type GearmanStats struct {
	Status  []Status
	Workers []Worker
}

func Load(addr string) (*GearmanStats, error) {
	c := NewGearAdminClient(addr)
	err := c.Connect()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	r := GearmanStats{}
	r.Status, err = c.Status()
	if err != nil {
		return nil, err
	}
	r.Workers, err = c.Workers()
	if err != nil {
		return nil, err
	}
	return &r, nil
}
