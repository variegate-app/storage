package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/variegate-app/storage/pkg/discovery"
)

type testService struct {
	Name string
	Link string
}

type testServiceBag [3]testService

const testServiceName = "test service"

func TestRegistery(t *testing.T) {
	ctx := context.Background()
	Idle := 1 * time.Second
	bag := testServiceBag{
		testService{testServiceName, "8080"},
		testService{testServiceName, "8081"},
		testService{testServiceName, "8082"},
	}
	r := NewRegistry(Config{Idle})
	names := []string{}

	t.Run("test Register", func(t *testing.T) {
		for _, s := range bag {
			name := discovery.GenerateInstanceID(s.Name)
			names = append(names, name)
			err := r.Register(ctx, name, s.Name, s.Link)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
		}
	})

	t.Run("test Discover", func(t *testing.T) {
		s1, err := r.Discover(ctx, "")
		assert.Equal(t, err, discovery.ErrNotFound)
		assert.Equal(t, 0, len(s1))

		s2, err := r.Discover(ctx, testServiceName)
		assert.Nil(t, err)
		assert.Equal(t, len(bag), len(s2))
	})

	t.Run("test Deregister", func(t *testing.T) {
		err := r.Deregister(ctx, "1", "WRONG NAME")
		assert.Equal(t, err, discovery.ErrNotFound)

		srcCnt := len(names)

		for _, n := range names {
			err := r.Deregister(ctx, n, testServiceName)
			assert.Nil(t, err)

			srcCnt--

			s1, err := r.Discover(ctx, testServiceName)
			if srcCnt > 0 {
				assert.Nil(t, err)
				assert.Equal(t, srcCnt, len(s1))
			} else {
				assert.Equal(t, err, discovery.ErrNotFound)
			}

		}
	})

	t.Run("test HealthCheck", func(t *testing.T) {
		r1 := NewRegistry(Config{Idle})
		names1 := []string{}
		for _, s := range bag {
			name := discovery.GenerateInstanceID(s.Name)
			names1 = append(names1, name)
			err := r1.Register(ctx, name, s.Name, s.Link)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
		}

		time.Sleep(Idle)

		err := r1.HealthCheck("1", "WRONG NAME")
		assert.Equal(t, err, discovery.ErrNotFoundService)

		s1, err := r1.Discover(ctx, testServiceName)
		assert.Nil(t, err, nil)
		assert.Equal(t, 0, len(s1))

		for _, n := range names1 {
			err = r1.HealthCheck(n, testServiceName)
			assert.Nil(t, err)
		}

		s1, err = r1.Discover(ctx, testServiceName)
		assert.Nil(t, err)
		assert.Equal(t, len(names1), len(s1))

		err = r.HealthCheck("0", testServiceName)
		assert.Equal(t, err, discovery.ErrNotFoundInstance)

		assert.Equal(t, Idle, r.GetIdleInterval())
	})
}
