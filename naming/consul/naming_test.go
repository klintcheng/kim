package consul

import (
	"sync"
	"testing"
	"time"

	"github.com/klintcheng/kim"
	"github.com/klintcheng/kim/naming"
	"github.com/stretchr/testify/assert"
)

func Test_Naming(t *testing.T) {
	ns, err := NewNaming("localhost:8500")
	assert.Nil(t, err)

	// 准备工作
	_ = ns.Deregister("test_1")
	_ = ns.Deregister("test_2")

	serviceName := "for_test"
	// 1. 注册 test_1
	err = ns.Register(&naming.DefaultService{
		Id:        "test_1",
		Name:      serviceName,
		Namespace: "",
		Address:   "localhost",
		Port:      8000,
		Protocol:  "ws",
		Tags:      []string{"tab1", "gate"},
	})
	assert.Nil(t, err)

	// 2. 服务发现
	servs, err := ns.Find(serviceName)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(servs))
	t.Log(servs)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// 3. 监听服务实时变化（新增）
	_ = ns.Subscribe(serviceName, func(services []kim.ServiceRegistration) {
		t.Log(len(services))

		assert.Equal(t, 2, len(services))
		assert.Equal(t, "test_2", services[1].ServiceID())
		wg.Done()
	})
	time.Sleep(time.Second)

	// 4. 注册 test_2 用于验证第3步
	err = ns.Register(&naming.DefaultService{
		Id:        "test_2",
		Name:      serviceName,
		Namespace: "",
		Address:   "localhost",
		Port:      8001,
		Protocol:  "ws",
		Tags:      []string{"tab2", "gate"},
	})
	assert.Nil(t, err)

	// 等 Watch 回调中的方法执行完成
	wg.Wait()

	_ = ns.Unsubscribe(serviceName)

	// 5. 服务发现
	servs, _ = ns.Find(serviceName, "gate")
	assert.Equal(t, 2, len(servs)) // <-- 必须有两个

	// 6. 服务发现, 验证tag查询
	servs, _ = ns.Find(serviceName, "tab2")
	assert.Equal(t, 1, len(servs)) // <-- 必须有1个
	assert.Equal(t, "test_2", servs[0].ServiceID())

	// 7. 注销test_2
	err = ns.Deregister("test_2")
	assert.Nil(t, err)

	// 8. 服务发现
	servs, err = ns.Find(serviceName)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(servs))
	assert.Equal(t, "test_1", servs[0].ServiceID())

	// 9. 注销test_1
	err = ns.Deregister("test_1")
	assert.Nil(t, err)

}
