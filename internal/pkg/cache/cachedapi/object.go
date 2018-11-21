package cachedapi

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Toggly/core/internal/api"
	"github.com/Toggly/core/internal/domain"
)

type cachedObjectAPI struct {
	owner       string
	projectCode domain.ProjectCode
	envCode     domain.EnvironmentCode
	engine      api.ObjectAPI
	cache       DataCache
}

func (c *cachedObjectAPI) basePath() string {
	return fmt.Sprintf("/own/%s/project/%s/env/%s/object", c.owner, c.projectCode, c.envCode)
}

func (c *cachedObjectAPI) List() ([]*domain.Object, error) {
	key := c.basePath()
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.List()
	})
	if err != nil {
		return nil, err
	}
	var list []*domain.Object
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *cachedObjectAPI) Get(code domain.ObjectCode) (*domain.Object, error) {
	key := fmt.Sprintf("%s/%s", c.basePath(), code)
	bytes, err := withCache(c.cache, key, func() (interface{}, error) {
		return c.engine.Get(code)
	})
	if err != nil {
		return nil, err
	}
	env := &domain.Object{}
	err = json.Unmarshal(bytes, env)
	if err != nil {
		return nil, err
	}
	return env, err
}

func (c *cachedObjectAPI) Create(info *api.ObjectInfo) (*domain.Object, error) {
	obj, err := c.engine.Create(info)
	if err != nil {
		return nil, err
	}
	scopes := make([]string, 0)
	scopes = append(scopes, c.basePath())
	c.cache.Flush(scopes...)
	return obj, nil
}

func (c *cachedObjectAPI) Update(info *api.ObjectInfo) (*domain.Object, error) {
	obj, err := c.engine.Update(info)
	if err != nil {
		return nil, err
	}
	scopes := make([]string, 0)
	scopes = append(scopes, c.basePath())
	scopes = append(scopes, fmt.Sprintf("%s/%s", c.basePath(), obj.Code))
	inheritors, err := c.engine.InheritorsFlatList(obj.Code)
	if err != nil {
		return nil, err
	}
	for _, i := range inheritors {
		scopes = append(scopes, fmt.Sprintf("/own/%s/project/%s/env/%s/object", i.Owner, i.ProjectCode, i.EnvCode))
		scopes = append(scopes, fmt.Sprintf("/own/%s/project/%s/env/%s/object/%s", i.Owner, i.ProjectCode, i.EnvCode, i.Code))
	}
	c.cache.Flush(scopes...)
	return obj, nil
}

func (c *cachedObjectAPI) Delete(code domain.ObjectCode) error {
	if err := c.engine.Delete(code); err != nil {
		return err
	}
	scopes := make([]string, 0)
	scopes = append(scopes, c.basePath())
	scopes = append(scopes, fmt.Sprintf("%s/%s", c.basePath(), code))
	c.cache.Flush(scopes...)
	return nil
}

func (c *cachedObjectAPI) InheritorsFlatList(code domain.ObjectCode) ([]*domain.Object, error) {
	log.Print("[WARN] API method `InheritorsFlatList` not cached. It may cause performance problems. Consider to avoid it.")
	return c.engine.InheritorsFlatList(code)
}
