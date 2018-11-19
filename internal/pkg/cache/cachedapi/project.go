package cachedapi

import (
	"encoding/json"
	"fmt"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
	"github.com/Toggly/core/internal/pkg/cache"
)

type cachedProjectAPI struct {
	owner  string
	engine api.ProjectAPI
	cache  cache.DataCache
}

func (c *cachedProjectAPI) basePath() string {
	return fmt.Sprintf("/own/%s/project", c.owner)
}

func (c *cachedProjectAPI) List() ([]*domain.Project, error) {
	key := c.basePath()
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.List()
	})
	if err != nil {
		return nil, err
	}
	var list []*domain.Project
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *cachedProjectAPI) Get(code domain.ProjectCode) (*domain.Project, error) {
	key := fmt.Sprintf("%s/%s", c.basePath(), code)
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.Get(code)
	})
	if err != nil {
		return nil, err
	}
	proj := &domain.Project{}
	err = json.Unmarshal(bytes, proj)
	if err != nil {
		return nil, err
	}
	return proj, err
}

func (c *cachedProjectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
	proj, err := c.engine.Create(info)
	if err != nil {
		return nil, err
	}
	c.cache.Flush([]string{
		c.basePath(),
	}...)
	return proj, nil
}

func (c *cachedProjectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
	proj, err := c.engine.Update(info)
	if err != nil {
		return nil, err
	}
	c.cache.Flush([]string{
		c.basePath(),
		fmt.Sprintf("%s/%s", c.basePath(), proj.Code),
	}...)
	return proj, nil
}

func (c *cachedProjectAPI) Delete(code domain.ProjectCode) error {
	if err := c.engine.Delete(code); err != nil {
		return err
	}
	c.cache.Flush([]string{
		c.basePath(),
		fmt.Sprintf("%s/%s", c.basePath(), code),
	}...)
	return nil
}

func (c *cachedProjectAPI) For(code domain.ProjectCode) api.ForProjectAPI {
	return &cachedForProjectAPI{
		owner:       c.owner,
		projectCode: code,
		engine:      c.engine.For(code).Environments(),
		cache:       c.cache,
	}
}

type cachedForProjectAPI struct {
	owner       string
	projectCode domain.ProjectCode
	engine      api.EnvironmentAPI
	cache       cache.DataCache
}

func (c *cachedForProjectAPI) Environments() api.EnvironmentAPI {
	return &cachedEnvAPI{
		owner:       c.owner,
		projectCode: c.projectCode,
		engine:      c.engine,
		cache:       c.cache,
	}
}
