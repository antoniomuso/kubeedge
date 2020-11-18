package config

import (
	"sync"

	"github.com/kubeedge/kubeedge/pkg/apis/componentconfig/edgecore/v1alpha1"
)

var Config Configure
var once sync.Once

type Configure struct {
	v1alpha1.FogManager
}

func InitConfigure(fogmanager *v1alpha1.FogManager) {
	once.Do(func() {
		Config = Configure{
			FogManager: *fogmanager,
		}
	})
}
