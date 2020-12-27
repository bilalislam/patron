package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize = 1024 * 1024
	target  = "bufnet"
)

var lis *bufconn.Listener

type server struct{}

func (s *server) SayHello(_ context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + in.Name}, nil
}

func TestMain(m *testing.M) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	code := m.Run()

	s.GracefulStop()

	os.Exit(code)
}

func bufDialer(_ context.Context, _address string) (net.Conn, error) {
	return lis.Dial()
}

func TestDial(t *testing.T) {
	conn, err := Dial(target, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	assert.NoError(t, conn.Close())
}

func TestDialContext(t *testing.T) {
	type args struct {
		opts []grpc.DialOption
	}
	tests := map[string]struct {
		args        args
		expectedErr string
	}{
		"success": {
			args: args{
				opts: []grpc.DialOption{grpc.WithContextDialer(bufDialer), grpc.WithInsecure()},
			},
		},
		"failure missing grpc.WithInsecure()": {
			args:        args{},
			expectedErr: "grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotConn, err := DialContext(context.Background(), target, tt.args.opts...)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
				assert.Nil(t, gotConn)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotConn)
			}
		})
	}
}

func TestSayHello(t *testing.T) {
	mtr := mocktracer.New()
	defer mtr.Reset()
	opentracing.SetGlobalTracer(mtr)
	ctx := context.Background()
	conn, err := DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, conn.Close())
	}()

	client := NewGreeterClient(conn)
	resp, err := client.SayHello(ctx, &HelloRequest{Name: "test"})
	require.NoError(t, err)
	assert.Equal(t, "Hello test", resp.GetMessage())
	expected := map[string]interface{}{
		"component": "grpc-client",
		"error":     false,
		"span.kind": ext.SpanKindEnum("producer"),
		"version":   "dev",
	}
	assert.Equal(t, expected, mtr.FinishedSpans()[0].Tags())
}
