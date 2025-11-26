package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
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
var once sync.Once

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
func Evaluate(service C.struct_String_t, law C.struct_String_t, parameters C.struct_Machine_law_Params_t, referenceDate C.struct_String_t, effectiveDate C.struct_String_t, overwriteInput C.struct_Machine_law_Params_t) C.struct_Machine_law_Result_t {
	logger := logger.New("main", os.Stdout, logrus.DebugLevel)
	caseManager := manager.New(logger)
	claimManager := inmemory.New(logger, caseManager)
	ruleResolver, err := ruleresolver.New()

	services, err := serviceprovider.New(logger, time.Now(), caseManager, claimManager, ruleResolver, serviceprovider.WithRuleServiceInMemory(), serviceprovider.WithOrganizationName("vorijkapp"))

	serviceString := convertStringStruct(service)
	//serviceString := "TOESLAGEN"
	lawString := convertStringStruct(law)
	//lawString := "zorgtoeslagwet"
	evalParams := convertParameters(parameters)
	//referenceDateString := convertStringStruct(referenceDate)

	// Set up context
	ctx := context.Background()

	resultCode := 0
	result, err := services.Evaluate(ctx, serviceString, lawString, evalParams, nil, nil, nil, "", false)
	if err != nil {
		resultCode = 3
	}
	if result == nil {
		resultCode = 4
	} else {
		if result.RequirementsMet == true {
			resultCode = 1
		} else {
			resultCode = 2
		}
	}
	testString := []byte("test")
	var cTestString unsafe.Pointer
	if len(testString) > 0 {
		cTestString = C.malloc(C.size_t(len(testString)))
		C.memcpy(cTestString, unsafe.Pointer(&testString[0]), C.size_t(len(testString)))
	}

	return C.struct_Machine_law_Result_t{
		resultCode: C.uint8_t(resultCode),
		resultMessage: C.struct_String_t{
			length: C.uint16_t(len(testString)),
			string: (*C.uint8_t)(cTestString),
		},
	}
}

//export EvaluateBetalingsRegelingRijk
func EvaluateBetalingsRegelingRijk(bsn C.struct_String_t, sociaalminimum C.uint32_t, inkomen C.uint32_t, totaleSchuld C.uint32_t, eerdereRegelingNietNagekomen C.int) C.struct_Machine_law_Result_t {
	logger := logger.New("main", os.Stdout, logrus.DebugLevel)
	evalParams := map[string]any{}
	bsnString := convertStringStruct(bsn)
	logger.Warning("bsn:")
	logger.Warning(bsnString)
	logger.Warning("end bsn")
	evalParams["BSN"] = bsnString

	evalParams["SOCIAAL_MINIMUM"] = C.int(sociaalminimum)
	evalParams["INKOMEN"] = C.int(inkomen)
	evalParams["TOTALE_SCHULD"] = C.int(totaleSchuld)
	nietNagekomen := C.int(eerdereRegelingNietNagekomen)
	if nietNagekomen > 0 {
		evalParams["EERDERE_REGELING_NIET_NAGEKOMEN"] = true
	} else {
		evalParams["EERDERE_REGELING_NIET_NAGEKOMEN"] = false
	}

	serviceString := "CJIB"
	lawString := "beleidsregels_betalingsregelingen_rijk"

	// Set up context
	ctx := context.Background()

	resultCode := 0
	resultString := ""
	result, err := services.Evaluate(ctx, serviceString, lawString, evalParams, nil, nil, nil, "", false)
	if err != nil {
		resultCode = 255
	}
	if result == nil {
		resultCode = 254
	} else {
		if result.RequirementsMet {
			resultCode = 1
			isGerechtigd := result.Output["is_gerechtigd"]
			aflosCapaciteit := result.Output["afloscapaciteit"]
			typeRegeling := result.Output["type_regeling"]
			termijnBedrag := result.Output["termijnbedrag"]
			aantalTermijnen := result.Output["aantal_termijnen"]
			resultString = fmt.Sprintf("is_gerechtigd:%t;afloscapaciteit:%d;type_regeling:%s;termijnbedrag:%d;aantal_termijnen:%d", isGerechtigd, aflosCapaciteit, typeRegeling, termijnBedrag, aantalTermijnen)
		} else {
			resultCode = 2
		}

	}
	testString := []byte(resultString)
	var cTestString unsafe.Pointer
	if len(testString) > 0 {
		cTestString = C.malloc(C.size_t(len(testString)))
		C.memcpy(cTestString, unsafe.Pointer(&testString[0]), C.size_t(len(testString)))
	}

	return C.struct_Machine_law_Result_t{
		resultCode: C.uint8_t(resultCode),
		resultMessage: C.struct_String_t{
			length: C.uint16_t(len(testString)),
			string: (*C.uint8_t)(cTestString),
		},
	}
}

//export EvaluateToeslagenWetBestaansMinimum
func EvaluateToeslagenWetBestaansMinimum(bsn C.struct_String_t, partner C.int, woningdeler C.int, bsnPartner C.struct_String_t, bsnWoningDeler C.struct_String_t, leeftijd C.int, leeftijdPartner C.int, leeftijdWoningdeler C.int) C.struct_Machine_law_Result_t {

	logger := logger.New("main", os.Stdout, logrus.DebugLevel)

	evalParams := map[string]any{}
	bsnString := convertStringStruct(bsn)
	logger.Warning("eval params:")
	logger.Warning("bsn")
	logger.Warning(bsnString)
	evalParams["BSN"] = bsnString
	logger.Warning("leeftijd")
	evalParams["LEEFTIJD"] = leeftijd
	if partner > 0 {
		evalParams["HEEFT_PARTNER"] = true
	} else {
		evalParams["HEEFT_PARTNER"] = false
	}

	evalParams["PARTNER_BSN"] = convertStringStruct(bsnPartner)
	evalParams["LEEFTIJD_PARTNER"] = leeftijdPartner
	if woningdeler > 0 {
		evalParams["HEEFT_WONINGDELER"] = true
	} else {
		evalParams["HEEFT_WONINGDELER"] = false
	}
	evalParams["WONINGDELER_BSN"] = convertStringStruct(bsnWoningDeler)
	evalParams["LEEFTIJD_WONINGDELER"] = leeftijdWoningdeler

	caseManager := manager.New(logger)
	claimManager := inmemory.New(logger, caseManager)
	ruleResolver, _ := ruleresolver.New()

	services, _ := serviceprovider.New(logger, time.Now(), caseManager, claimManager, ruleResolver, serviceprovider.WithRuleServiceInMemory(), serviceprovider.WithOrganizationName("vorijkapp"))

	serviceString := "UWV"
	lawString := "toeslagenwet"

	// Set up context
	ctx := context.Background()

	resultCode := 0
	resultString := ""
	result, err := services.Evaluate(ctx, serviceString, lawString, evalParams, nil, nil, nil, "", false)
	if err != nil {
		resultCode = 255
	}
	if result == nil {
		resultCode = 254
	} else {
		if result.RequirementsMet {
			resultCode = 1
			sociaalMinimum := result.Output["sociaal_minimum"]
			resultString = fmt.Sprintf("is_gerechtigd:%d;", sociaalMinimum)

		} else {
			resultCode = 2
		}
	}
	testString := []byte(resultString)
	var cTestString unsafe.Pointer
	if len(testString) > 0 {
		cTestString = C.malloc(C.size_t(len(testString)))
		C.memcpy(cTestString, unsafe.Pointer(&testString[0]), C.size_t(len(testString)))
	}

	return C.struct_Machine_law_Result_t{
		resultCode: C.uint8_t(resultCode),
		resultMessage: C.struct_String_t{
			length: C.uint16_t(len(testString)),
			string: (*C.uint8_t)(cTestString),
		},
	}
}

func convertParameters(parametersStruct C.struct_Machine_law_Params_t) map[string]any {
	//evalParams := map[string]any{}
	/*parameters := strings.Split(convertStringStruct(parametersStruct.params), ";")
	for _, parameter := range parameters {
		paramFields := strings.Split(parameter, ":")
		evalParams[paramFields[0]] = paramFields[1]
	}*/
	evalParams := map[string]any{
		"BSN": "100000001",
	}
	return evalParams
}

func convertStringStruct(stringStruct C.struct_String_t) string {
	stringBytes := C.GoBytes(unsafe.Pointer(stringStruct.string), C.int(stringStruct.length))
	return string(stringBytes)
}

//export enforce_binding
func enforce_binding() {}

func main() {

}
