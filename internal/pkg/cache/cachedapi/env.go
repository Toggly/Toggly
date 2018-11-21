package cachedapi

import (
	"encoding/json"
	"fmt"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/pkg/cache"
)

type cachedEnvAPI struct {
	owner       string
	projectCode domain.ProjectCode
	engine      api.EnvironmentAPI
	cache       cache.DataCache
}

func (c *cachedEnvAPI) basePath() string {
	return fmt.Sprintf("/own/%s/project/%s/env", c.owner, c.projectCode)
}

func (c *cachedEnvAPI) List() ([]*domain.Environment, error) {
	key := c.basePath()
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.List()
	})
	if err != nil {
		return nil, err
	}
	var list []*domain.Environment
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *cachedEnvAPI) Get(code domain.EnvironmentCode) (*domain.Environment, error) {
	key := fmt.Sprintf("%s/%s", c.basePath(), code)
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.Get(code)
	})
	if err != nil {
		return nil, err
	}
	env := &domain.Environment{}
	err = json.Unmarshal(bytes, env)
	if err != nil {
		return nil, err
	}
	return env, err
}

func (c *cachedEnvAPI) Create(info *api.EnvironmentInfo) (*domain.Environment, error) {
	env, err := c.engine.Create(info)
	if err != nil {
		return nil, err
	}
	c.cache.Flush([]string{
		c.basePath(),
	}...)
	return env, nil
}

func (c *cachedEnvAPI) Update(info *api.EnvironmentInfo) (*domain.Environment, error) {
	env, err := c.engine.Update(info)
	if err != nil {
		return nil, err
	}
	c.cache.Flush([]string{
		c.basePath(),
		fmt.Sprintf("%s/%s", c.basePath(), env.Code),
	}...)
	return env, nil
}

func (c *cachedEnvAPI) Delete(code domain.EnvironmentCode) error {
	if err := c.engine.Delete(code); err != nil {
		return err
	}
	c.cache.Flush([]string{
		c.basePath(),
		fmt.Sprintf("%s/%s", c.basePath(), code),
	}...)
	return nil
}

func (c *cachedEnvAPI) For(code domain.EnvironmentCode) api.ForObjectAPI {
	return &cachedForObjectAPI{
		owner:       c.owner,
		projectCode: c.projectCode,
		envCode:     code,
		engine:      c.engine.For(code).Objects(),
		cache:       c.cache,
	}
}

type cachedForObjectAPI struct {
	owner       string
	projectCode domain.ProjectCode
	envCode     domain.EnvironmentCode
	engine      api.ObjectAPI
	cache       cache.DataCache
}

func (c *cachedForObjectAPI) Objects() api.ObjectAPI {
	return &cachedObjectAPI{
		owner:       c.owner,
		projectCode: c.projectCode,
		envCode:     c.envCode,
		engine:      c.engine,
		cache:       c.cache,
	}
}
