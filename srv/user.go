package srv

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ability-sh/abi-db/client/service"
	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-micro/grpc"
	"github.com/ability-sh/abi-micro/micro"
)

func (s *Server) Login(ctx micro.Context, task *LoginTask) (*User, error) {

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

	id, err := collection.Get(cc, fmt.Sprintf("name/%s", task.Name))

	if err != nil {
		e, ok := err.(*errors.Error)
		if ok && e.Errno == 404 {
			return nil, errors.Errorf(ERRNO_LOGIN, "user or password error")
		}
		return nil, err
	}

	object, err := collection.GetObject(cc, fmt.Sprintf("id/%s", string(id)))

	if err != nil {
		e, ok := err.(*errors.Error)
		if ok && e.Errno == 404 {
			return nil, errors.Errorf(ERRNO_LOGIN, "user or password error")
		}
		return nil, err
	}

	u := &User{}

	dynamic.SetValue(&u, object)

	if u.Password != config.SecPassword(task.Password) {
		return nil, errors.Errorf(ERRNO_LOGIN, "user or password error")
	}

	u.Password = ""

	return u, nil
}

func (s *Server) Create(ctx micro.Context, task *UserCreateTask) (*User, error) {

	if task.Name == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter name is incorrect")
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

	id := config.NewID(ctx)

	u := &User{Id: id, Name: task.Name, Nick: task.Nick, Password: config.SecPassword(task.Password), Ctime: time.Now().Unix()}

	_, err = collection.Exec(cc, `
	var user = ${user};
	var k_name = collection + 'name/' + user.name;
	var k_id = collection + 'id/' + user.id;
	var k_nick = collection + 'nick/' + user.nick;
	if(get(k_name)) {
		throw 'name already exists'
	}
	if(user.nick && get(k_nick)) {
		throw 'nick already exists'
	}
	put(k_name,user.id);
	if(user.nick) {
		put(k_nick,user.id);
	}
	put(k_id,JSON.stringify(user));
	`, map[string]interface{}{"user": u})

	if err != nil {
		return nil, err
	}

	u.Password = ""

	return u, nil
}

func (s *Server) Get(ctx micro.Context, task *UserGetTask) (*User, error) {

	if task.Id == "" && task.Name == "" && task.Nick == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter id is incorrect")
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

	u := &User{}

	if task.Id != "" {

		b, err := collection.Get(cc, fmt.Sprintf("id/%s", task.Id))

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, u)

		if err != nil {
			return nil, err
		}

	} else if task.Name != "" {

		id, err := collection.Get(cc, fmt.Sprintf("name/%s", task.Name))

		if err != nil {

			e, ok := err.(*errors.Error)

			if ok && e.Errno == 404 {

				if task.AutoCreated {

					u, err := s.Create(ctx, &UserCreateTask{Name: task.Name})

					if err != nil {
						return nil, err
					}

					return u, nil

				} else {
					return nil, err
				}

			}

			return nil, err
		}

		b, err := collection.Get(cc, fmt.Sprintf("id/%s", string(id)))

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, u)

		if err != nil {
			return nil, err
		}

	} else if task.Nick != "" {

		id, err := collection.Get(cc, fmt.Sprintf("nick/%s", task.Nick))

		if err != nil {
			return nil, err
		}

		b, err := collection.Get(cc, fmt.Sprintf("id/%s", string(id)))

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(b, u)

		if err != nil {
			return nil, err
		}
	}

	u.Password = ""

	return u, nil
}

func (s *Server) Set(ctx micro.Context, task *UserSetTask) (*User, error) {

	if task.Id == "" {
		return nil, errors.Errorf(ERRNO_INPUT_DATA, "The parameter id is incorrect")
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

	data := map[string]interface{}{"id": task.Id}

	if task.Name != nil {
		data["name"] = *task.Name
	}

	if task.Nick != nil {
		data["nick"] = *task.Nick
	}

	if task.Password != nil {
		data["password"] = config.SecPassword(*task.Nick)
	}

	text, err := collection.Exec(cc, `(function(){
	var data = ${data};
	var k_id = collection + 'id/' + data.id;
	var text = get(k_id);
	if(!text) {
		throw 'The user does not exist'
	}
	var user = JSON.parse(text);
	var hasChanged = false;
	if(data.name && data.name != user.name) {
		var k_name = collection + 'name/' + data.name;
		if(get(k_name)) {
			throw 'name already exists'
		}
		put(k_name,data.id);
		del(collection + 'name/' + user.name)
		user.name = data.name;
		hasChanged = true;
	}
	if(data.nick !== undefined && data.nick != user.nick) {
		if(data.nick) {
			var k_nick = collection + 'nick/' + data.nick;
			if(get(k_nick)) {
				throw 'nick already exists'
			}
			put(k_nick,data.id);
		}
		if(user.nick) {
			del(collection + 'nick/' + user.nick)
		}
		user.nick = data.nick;
		hasChanged = true;
	}
	if(hasChanged) {
		var s = JSON.stringify(user);
		put(k_id,s);
		return s;
	}
	return text;
	})()
	`, map[string]interface{}{"data": data})

	if err != nil {
		return nil, err
	}

	u := &User{}

	err = json.Unmarshal([]byte(text), u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Server) BatchGet(ctx micro.Context, task *UserBatchGetTask) ([]*User, error) {

	vs := []*User{}

	if task.Ids == "" {
		return vs, nil
	}

	get := UserGetTask{}

	for _, id := range strings.Split(task.Ids, ",") {
		get.Id = id
		u, err := s.Get(ctx, &get)
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
