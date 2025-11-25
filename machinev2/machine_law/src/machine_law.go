package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/casemanager/manager"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/claimmanager/inmemory"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/ruleresolver"
	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/service/serviceprovider"
	"github.com/sirupsen/logrus"

	"github.com/Ealanrian/poc-machine-law-vorijk-fieldlab/machinev2/machine/logger"
)

//#include "machine_law.h"
//#include <stdlib.h>
//#include <string.h>
import "C"

var services *serviceprovider.Services
var started = false

//export machineLawStandalone
func machineLawStandalone() C.int {
	logger := logger.New("main", os.Stdout, logrus.DebugLevel)
	caseManager := manager.New(logger)
	claimManager := inmemory.New(logger, caseManager)
	ruleResolver, err := ruleresolver.New()

	if err != nil {
		return -2
	}
	services, err := serviceprovider.New(logger, time.Now(), caseManager, claimManager, ruleResolver, serviceprovider.WithRuleServiceInMemory(), serviceprovider.WithOrganizationName("vorijkapp"))
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "new service resolver") {
			return -3
		}
		return -1
	}
	if services.InStandAloneMode() {
		return 1
	} else {
		return 0
	}
}

//export Evaluate
func Evaluate(service C.struct_String_t, law C.struct_String_t, parameters C.struct_Machine_law_Params_t, referenceDate C.struct_String_t, effectiveDate C.struct_String_t, overwriteInput C.struct_Machine_law_Params_t, requestedOutput C.struct_String_t, approved C.int) C.struct_Machine_law_Result_t {
	testString := []byte("test")
	var cTestString unsafe.Pointer
	if len(testString) > 0 {
		cTestString = C.malloc(C.size_t(len(testString)))
		C.memcpy(cTestString, unsafe.Pointer(&testString[0]), C.size_t(len(testString)))
	}
	return C.struct_Machine_law_Result_t{
		resultCode: 1,
		resultMessage: C.struct_String_t{
			length: C.uint16_t(len(testString)),
			string: (*C.uint8_t)(cTestString),
		},
	}
}

//export enforce_binding
func enforce_binding() {}

func main() {

}
