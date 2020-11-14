package fogmanager

import (
	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"k8s.io/klog"
)

const (
	OK = "OK"

	GroupResource    = "resource"
	OperationFogSync = "fog-internal-sync"

	OperationFunctionAction = "action"

	OperationFunctionActionResult = "action_result"

	EdgeFunctionModel   = "edgefunction"
	CloudFunctionModel  = "funcmgr"
	CloudControlerModel = "edgecontroller"

	// EdgeHubModuleName Where is located enum?
	// EdgeHubModuleName = "websocket"
)

func (m *fogManager) process(msg model.Message) {
	// Do something with the message
	println(msg.GetSource())
}

func (m *fogManager) assignLabel() {
	// msg := model.NewMessage("").BuildRouter(m.Name(), modules.HubGroup, model.ResourceTypeNode, model.UpdateOperation)
	// msg.FillBody("")
	// beehiveContext.SendToGroup(modules.HubGroup, *msg)
}

func (m *fogManager) runFogManager() {
	go func() {
		for {
			select {
			case <-beehiveContext.Done():
				klog.Warning("FogManager mainloop stop")
				return
			default:
			}
			if msg, err := beehiveContext.Receive(m.Name()); err == nil {
				klog.V(2).Infof("get a message %+v", msg)
				// println(msg.GetSource())
				m.process(msg)
			} else {
				klog.Errorf("get a message %+v: %v", msg, err)
			}
			m.assignLabel()
		}
	}()
}
