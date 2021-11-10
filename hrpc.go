package hrpc

import (
	"github.com/hirpc/arsenal/uniqueid"
	"github.com/hirpc/hrpc/database"
	"github.com/hirpc/hrpc/database/category"
	"github.com/hirpc/hrpc/option"
	"github.com/hirpc/hrpc/server"
)

func register() error {
	if err := uniqueid.Register(); err != nil {
		return err
	}
	return nil
}

// NewServer is the entrance of the framwork
func NewServer(opts ...option.Option) (server.Server, error) {
	// register some dependent components
	if err := register(); err != nil {
		return nil, err
	}

	var opt = &option.Options{
		ID:         uniqueid.String(),
		ListenPort: 8888,
		ENV:        option.Development,
		DBs:        make(map[category.Category]database.Database),
	}
	for _, o := range opts {
		o(opt)
	}
	if err := opt.Valid(); err != nil {
		return nil, err
	}
	// returns a new grpc server with desired options
	return server.NewGRPC(opt)
}
