package srv

import (
	"fmt"
	"strings"

	"github.com/ability-sh/abi-db/client/service"
	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-micro/grpc"
	"github.com/ability-sh/abi-micro/micro"
)

func (s *Server) InfoGet(ctx micro.Context, task *InfoGetTask) (interface{}, error) {

	if task.Id == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter id is incorrect")
	}

	if task.Key == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter key is incorrect")
	}

	config, err := GetConfigService(ctx, SERVICE_CONFIG)

	if err != nil {
		return nil, err
	}

	client, err := service.GetClient(ctx, config.Db)

	if err != nil {
		return nil, err
	}

	cc := grpc.NewGRPCContext(ctx)

	collection := client.Collection(config.Collection)

	object, err := collection.GetObject(cc, fmt.Sprintf("info/%s/%s", task.Id, task.Key))

	if err != nil {
		e, ok := err.(*errors.Error)
		if ok && e.Errno == 404 {
			return nil, nil
		}
		return nil, err
	}

	return object, nil
}

func (s *Server) InfoSet(ctx micro.Context, task *InfoSetTask) (interface{}, error) {

	if task.Id == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter id is incorrect")
	}

	if task.Key == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter key is incorrect")
	}

	config, err := GetConfigService(ctx, SERVICE_CONFIG)

	if err != nil {
		return nil, err
	}

	client, err := service.GetClient(ctx, config.Db)

	if err != nil {
		return nil, err
	}

	cc := grpc.NewGRPCContext(ctx)

	collection := client.Collection(config.Collection)

	k := fmt.Sprintf("info/%s/%s", task.Id, task.Key)

	err = collection.MergeObject(cc, k, task.Info)

	if err != nil {
		return nil, err
	}

	return collection.GetObject(cc, k)
}

func (s *Server) InfoBatchGet(ctx micro.Context, task *InfoBatchGetTask) ([]interface{}, error) {

	if task.Key == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter key is incorrect")
	}

	vs := []interface{}{}

	if task.Ids == "" {
		return vs, nil
	}

	get := InfoGetTask{}

	for _, id := range strings.Split(task.Ids, ",") {
		get.Id = id
		get.Key = task.Key
		u, err := s.InfoGet(ctx, &get)
		if err != nil {
			e, ok := err.(*errors.Error)
			if ok && e.Errno == 404 {
				vs = append(vs, nil)
				continue
			}
			return nil, err
		}
		vs = append(vs, u)
	}

	return vs, nil
}
