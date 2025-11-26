#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#if _WIN32
#include <windows.h>
#else
#include <pthread.h>
#include <unistd.h>
#endif

#if _WIN32
#define FFI_PLUGIN_EXPORT __declspec(dllexport)
#else
#define FFI_PLUGIN_EXPORT
#endif


 struct String_t{
    uint16_t length;
    uint8_t* string;
} ;


 struct Machine_law_Params_t{
    struct String_t params;
} ;

 struct Machine_law_Result_t{
    uint8_t resultCode;
    struct String_t resultMessage;
} ;


FFI_PLUGIN_EXPORT struct Machine_law_Result_t Evaluate(struct String_t service, struct String_t law, struct Machine_law_Params_t parameters, struct String_t referenceDate, struct String_t effectiveDate, struct Machine_law_Params_t overwriteInput);
FFI_PLUGIN_EXPORT struct Machine_law_Result_t EvaluateBetalingsRegelingRijk(struct String_t bsn, uint32_t sociaalMinimum, uint32_t inkomen, uint32_t totaleSchuld, int eerdereRegelingNietNagekomen );
FFI_PLUGIN_EXPORT struct Machine_law_Result_t EvaluateToeslagenWetBestaansMinimum(struct String_t bsn, int partner, int woningdeler, struct String_t bsnPartner, struct String_t woningdelerBsn, int leeftijd, int leeftijdPartner, int leeftijdWoningDeler);
FFI_PLUGIN_EXPORT int machineLawStandalone();

FFI_PLUGIN_EXPORT void freeParams(struct Machine_law_Params_t params);
FFI_PLUGIN_EXPORT void freeRuleResult(struct Machine_law_Result_t result);
