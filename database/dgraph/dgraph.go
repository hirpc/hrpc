package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/hirpc/hrpc/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Proxy represents a collection of methods to operate dGraph
type Proxy interface {
	// Alter can be used to set schema.
	// For example:
	// 	op := &api.Operation{
	//		Schema: 		 `name: string @index(exact) .`,
	//		RunInBackground: true
	//	}
	Alter(ctx context.Context, op *api.Operation) error
	// Transaction should be called normally to run a Mutation or Upsert(Query + Mutation).
	// For example:
	// 	(Mutation)
	//		type Person struct {
	//			Name string `json:"name,omitempty"`
	//		}
	//		p := Persion{Name: "xxx"}
	//		b, err := json.Marshal(p)
	//		mu := &api.Mutation{
	//			SetJson: b,
	//		}
	//		req := &api.Request{CommitNow: true, Mutations: []*api.Mutation{mu}}
	//		rsp, err := Transaction(ctx, req)
	//	(Query + Mutation) WITHOUT conditions
	//		q := `query {user as var(func: eq(email, "wrong_email@example.com"))}`
	//		mu := &api.Mutation{SetNquads: []byte(`uid(user) <email> "correct_email@example.com" .`)}
	//		req := &api.Request{Query: q, Mutations: []*api.Mutation{mu}, CommitNow: true}
	//		rsp, err := Transaction(ctx, req)
	//	(Query + Mutation) WITH conditions
	//		q := `query {user as var(func: eq(email, "wrong_email@example.com"))}`
	//		mu := &api.Mutation{
	//			Cond: `@if(eq(len(user), 1))`,
	//			SetNquads: []byte(`uid(user) <email> "correct_email@example.com" .`),
	//		}
	//		req := &api.Request{Query: q, Mutations: []*api.Mutation{mu}, CommitNow: true}
	//		rsp, err := Transaction(ctx, req)
	Transaction(ctx context.Context, req *api.Request) ([]byte, error)
	// Query can be used in two ways, including general query and query with variables.
	// It is useful to increase read speed.
	// the param of `vars` can be nil if you do not want to pass variables to the Query statements.
	Query(ctx context.Context, q string, vars map[string]string) ([]byte, error)
}

var dg *dGraph

type dGraph struct {
	client *dgo.Dgraph
	conns  []*grpc.ClientConn

	option Option
}

func New() *dGraph {
	if dg != nil {
		dg.Destory()
	}
	dg = &dGraph{}
	return dg
}

// ------------ for implemation database.Database START-------- //

func (d *dGraph) Load(src []byte) error {
	if err := json.Unmarshal(src, &d.option); err != nil {
		return err
	}
	return nil
}
func (d *dGraph) Connect() error {
	d.Destory()

	var clients []api.DgraphClient
	for _, target := range d.option.Targets {
		c, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		d.conns = append(d.conns, c)
		clients = append(clients, api.NewDgraphClient(c))
	}
	d.client = dgo.NewDgraphClient(clients...)
	if d.option.Credential.User != "" && d.option.Credential.Password != "" {
		if d.option.Credential.Namespace == 0 {
			d.client.Login(
				context.Background(),
				d.option.Credential.User,
				d.option.Credential.Password,
			)
		} else {
			// SHOULD ONLY valid for Dgraph v21.03 or above
			d.client.LoginIntoNamespace(
				context.Background(),
				d.option.Credential.User,
				d.option.Credential.Password,
				d.option.Credential.Namespace,
			)
		}
	}
	return nil
}
func (d *dGraph) Name() string {
	return "dgraph"
}
func (d *dGraph) Destory() {
	for _, c := range d.conns {
		c.Close()
	}
}

// ------------ for implemation database.Database END-------- //

func (d *dGraph) Alter(ctx context.Context, op *api.Operation) error {
	return dg.client.Alter(ctx, op)
}

func (d *dGraph) Transaction(ctx context.Context, req *api.Request) ([]byte, error) {
	txn := dg.client.NewTxn()
	defer txn.Discard(ctx)

	rsp, err := txn.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := txn.Commit(ctx); err != nil {
		return nil, err
	}
	return rsp.GetJson(), nil
}

func (d *dGraph) Query(ctx context.Context, q string, vars map[string]string) ([]byte, error) {
	txn := dg.client.NewReadOnlyTxn()
	defer txn.Discard(ctx)

	if vars != nil {
		rsp, err := txn.QueryWithVars(ctx, q, vars)
		if err != nil {
			return nil, err
		}
		return rsp.GetJson(), nil
	}

	rsp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	return rsp.GetJson(), nil
}

// Client returns the handler to operate mysql if success
func Client() Proxy {
	return dg
}

var _ database.Database = (*dGraph)(nil)
