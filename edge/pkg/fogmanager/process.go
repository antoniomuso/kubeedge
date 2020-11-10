package fogmanager

import (
	"github.com/kubeedge/beehive/pkg/core/model"
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
)

func (m *fogManager) process(msg model.Message) {
	// Do something with the message
}

func (m *fogManager) runFogManager() {

}
