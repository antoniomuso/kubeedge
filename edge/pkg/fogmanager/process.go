package fogmanager

import (
	"encoding/json"

	beehiveContext "github.com/kubeedge/beehive/pkg/core/context"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/edge/pkg/common/modules"
	v1 "k8s.io/api/core/v1"
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
	println(msg.GetOperation())

	klog.V(2).Infof(" Arrived message from %+v to %+v", msg.GetResource(), m.Name())
	m.assignLabel()
}

func (m *fogManager) assignLabel() {
	node := &v1.Node{}
	node.APIVersion = "v1"
	node.Kind = "Node"
	node.ObjectMeta.Name = "edge-node"
	node.ObjectMeta.Labels = make(map[string]string)
	node.ObjectMeta.Labels["fog-colony-name"] = "home"

	msg := model.NewMessage("").BuildRouter(m.Name(), modules.HubGroup, model.ResourceTypeNode, model.UpdateOperation)
	content, ok := json.Marshal(node)
	if ok != nil {
		klog.Errorf("Marhing error")
	}
	msg.Content = content
	beehiveContext.SendToGroup(modules.HubGroup, *msg)
}

func (m *fogManager) runFogManager() {
	klog.Infof("Fog manager Start")
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

		}
	}()
}
