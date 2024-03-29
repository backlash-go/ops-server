package db

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisHMSet(t *testing.T) {
	type args struct {
		token     string
		keyFields map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestRedisHMSet",
			args: args{
				token: "k3",
				keyFields: map[string]interface{}{
					"name": "xxb3",
					"age":  1881,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisHMSet(tt.args.token, tt.args.keyFields); (err != nil) != tt.wantErr {
				t.Errorf("RedisHMSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisSetKeyTtl(t *testing.T) {
	type args struct {
		token  string
		expire time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestRedisSetKeyTtl",
			args: args{
				token:  "k2",
				expire: time.Minute * 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisSetKeyTtl(tt.args.token, tt.args.expire); (err != nil) != tt.wantErr {
				t.Errorf("RedisSetKeyTtl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis(t *testing.T) {

}

func TestRedisHMGet(t *testing.T) {
	type args struct {
		token     string
		keyFields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestRedisHMGet",
			args:    args{token: "k1", keyFields: []string{"age1", "name1"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisHMGet(tt.args.token, tt.args.keyFields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisHMGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Printf("got type is %T\n", got)
			fmt.Printf("got type is %+v\n", got)
			fmt.Printf("got-item type is %T\n", got[0])

		})
	}
}

func TestRedisKeyIsExist(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestRedisKeyIsExist", args: args{token: "k11"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisKeyIsExist(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisKeyIsExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("got is :%+v", got)

		})
	}
}

func TestRedisHGetAll(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestRedisHGetAll", args: args{token: "c8sodaauof2mc0aq92s01"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RedisHGetAll(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisHGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				fmt.Printf("got is nil")

			}

			fmt.Printf("got len is %+v", len(got))


			fmt.Printf("got is %+v", got)
		})
	}
}

func TestRedisDelKeys(t *testing.T) {
	type args struct {
		key []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "TestRedisDelKeys", args: args{key: []string{"k5"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RedisDelKeys(tt.args.key...); (err != nil) != tt.wantErr {
				t.Errorf("RedisDelKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}