package docker

import (
	"context"
	"github.com/docker/docker/client"
	"reflect"
	"testing"
)

func Client() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	//defer cli.Close()
	return cli
}

func TestDocker_IsRunning(t *testing.T) {
	type fields struct {
		ctx    context.Context
		client *client.Client
		n      int64
	}
	type args struct {
		ctx           context.Context
		containerName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "success",
			fields: fields{ctx: context.Background(), client: Client(), n: int64(10)},
			args:   args{ctx: context.Background(), containerName: "/kind-control-plane"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Docker{
				ctx:    tt.fields.ctx,
				client: tt.fields.client,
				n:      tt.fields.n,
			}
			if got := d.IsRunning(tt.args.ctx, tt.args.containerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}
