package main

import (
	"os"
	"time"

	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/casemanager/manager"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/claimmanager/inmemory"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/ruleresolver"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/service/serviceprovider"
	"github.com/sirupsen/logrus"

	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/logger"
)

//#include "machine_law.h"

import "C"

var services *serviceprovider.Services
var started = false

//export startMachineLawEngine
func startMachineLawEngine() {
	logger := logger.New("main", os.Stdout, logrus.DebugLevel)
	caseManager := manager.New(logger)
	claimManager := inmemory.New(logger, caseManager)
	ruleResolver, _ := ruleresolver.New()
	_, _ = serviceprovider.New(logger, time.Now(), caseManager, claimManager, ruleResolver, serviceprovider.WithRuleServiceInMemory())
	started = true
}

//export machineLawStandalone
func machineLawStandalone() int {
	logger := logger.New("main", os.Stdout, logrus.DebugLevel)
	caseManager := manager.New(logger)
	claimManager := inmemory.New(logger, caseManager)
	ruleResolver, err := ruleresolver.New()

	if err != nil {
		return -2
	}
	services, err := serviceprovider.New(logger, time.Now(), caseManager, claimManager, ruleResolver, serviceprovider.WithRuleServiceInMemory(), serviceprovider.WithOrganizationName("vorijkapp"))
	if err != nil {
		return -1
	}
	if services.InStandAloneMode() {
		return 1
	} else {
		return 0
	}
}

//export enforce_binding
func enforce_binding() {}

func main() {

}
