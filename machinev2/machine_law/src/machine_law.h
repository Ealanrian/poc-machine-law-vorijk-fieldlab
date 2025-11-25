#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

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

#ifndef MACHINE_LAW_H
#define MACHINE_LAW_H
#include <stddef.h>

typedef struct {
    // Define the structure for parameters
    int numParams;
    char** paramNames;
    void** paramValues;
} Params_t;

typedef struct {
    int resultCode;
    char* resultMessage;
} Result_t;

FFI_PLUGIN_EXPORT Result_t Evaluate(const char* service, const char* law, Params_t parameters, const char* referenceDate, const char* effectiveDate, Params_t overwriteInput, const char* requestedOutput, int approved);
FFI_PLUGIN_EXPORT int machineLawStandalone();
FFI_PLUGIN_EXPORT void startMachineLawEngine();
void freeParams(Params_t params);
void freeRuleResult(Result_t result);

#endif