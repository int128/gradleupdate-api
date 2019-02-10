package gateways

import (
	"context"
	"strings"
	"sync"

	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

func NewToggles() gateways.Toggles {
	return &togglesCache{
		Base: &togglesData{},
	}
}

type togglesCache struct {
	Base gateways.Toggles
	l    sync.Mutex
	v    *config.Toggles
}

func (r *togglesCache) Get(ctx context.Context) (*config.Toggles, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if r.v != nil {
		return r.v, nil
	}
	v, err := r.Base.Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting toggles")
	}
	r.v = v
	return r.v, nil
}

type togglesData struct{}

func (r *togglesData) Get(ctx context.Context) (*config.Toggles, error) {
	var e togglesEntity
	k := togglesKey(ctx, "DEFAULT")
	if err := datastore.Get(ctx, k, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return &config.Toggles{}, nil
		}
		return nil, errors.Wrapf(err, "error while getting the entity")
	}
	return &config.Toggles{
		BatchSendUpdatesOwners: strings.Split(e.BatchSendUpdatesOwners, ","),
	}, nil
}

func togglesKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, "Toggles", name, 0, nil)
}

type togglesEntity struct {
	BatchSendUpdatesOwners string // comma-separated strings
}
