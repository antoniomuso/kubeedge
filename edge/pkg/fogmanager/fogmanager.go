package fogmanager

import (
	"time"

	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/core"
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"

	// fogmanagerconfig "github.com/kubeedge/kubeedge/edge/pkg/fogmanager/config"
	// "github.com/kubeedge/kubeedge/edge/pkg/fogmanager/dao"
	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

//constant fogmanager module name
const (
	FogManagerModuleName = "fogManager"
)

type fogManager struct {
	enable bool
}

func newFogManager(enable bool) *fogManager {
	return &fogManager{enable: enable}
}

// Register register fogmanager
func Register(fogManager *v1alpha1.FogManager) {
	// fogmanagerconfig.InitConfigure(fogManager)
	fog := newFogManager(fogManager.Enable)
	// initDBTable(fog)
	core.Register(fog)
}

// initDBTable create table
func initDBTable(module core.Module) {
	klog.Infof("Begin to register %v db model", module.Name())
	if !module.Enable() {
		klog.Infof("Module %s is disabled, DB fog for it will not be registered", module.Name())
		return
	}
	// orm.RegisterModel(new(dao.Fog))
}

func (*fogManager) Name() string {
	return FogManagerModuleName
}

func (*fogManager) Group() string {
	return modules.FogGroup
}

func (m *fogManager) Enable() bool {
	return m.enable
}

func (m *fogManager) Start() {
	klog.Infof("Begin to register %v db model", module.Name())
	go func() {
		period := getSyncInterval()
		timer := time.NewTimer(period)
		for {
			select {
			case <-beehiveContext.Done():
				klog.Warning("FogManager stop")
				return
			case <-timer.C:
				timer.Reset(period)
				msg := model.NewMessage("").BuildRouter(FogManagerModuleName, GroupResource, model.ResourceTypePodStatus, OperationFogSync)
				beehiveContext.Send(FogManagerModuleName, *msg)
			}
		}
	}()

	m.runFogManager()
}

func getSyncInterval() time.Duration {
	// TODO: Change time interval
	return time.Duration(time.Second * 10)
}
