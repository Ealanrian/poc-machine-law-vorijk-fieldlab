
import 'dart:async';
import 'dart:convert';
import 'dart:ffi';
import 'dart:io';
import 'dart:typed_data';


import 'package:ffi/ffi.dart';

import 'machine_law_bindings_generated.dart';

/// A very short-lived native function.
///
/// For very short-lived functions, it is fine to call them on the main isolate.
/// They will block the Dart execution while running the native function, so
/// only do this for native functions which are guaranteed to be short-lived.
//int sum(int a, int b) => _bindings.sum(a, b);


int machineLawStandalone() => _bindings.machineLawStandalone();


int machineLawEvaluate() {

  final servicename = "TOESLAGEN";
  final Pointer<String_t> service = createStringPointer(servicename);
  final Pointer<String_t> law = createStringPointer("zorgtoeslagwet");
  Pointer<Machine_law_Params_t> parameters = createParametersPointer("bsn:100000001");
  final Pointer<String_t> referenceDate = createStringPointer("");
  final Pointer<String_t> effectiveDate = createStringPointer("");
  final Pointer<Machine_law_Params_t> overwriteInput = createParametersPointer("bsn:100000001");
  Machine_law_Result_t result =_bindings.Evaluate(service.ref, law.ref, parameters.ref, referenceDate.ref,
      effectiveDate.ref, overwriteInput.ref);


  freeString(service);
  freeString(law);
  return result.resultCode;
}


String evalBetalingsRegeling(String bsn, int sociaalMinimum, int inkomen, int totaleschuld, bool nietNagekomen) {
  final Pointer<String_t> bsnString = createStringPointer(bsn);

  int nietNagekomenInt = 0;
  if (nietNagekomen == true) {
    nietNagekomenInt = 1;
  }
  var resultString = "";
  Machine_law_Result_t result =_bindings.EvaluateBetalingsRegelingRijk(bsnString.ref,sociaalMinimum,inkomen, totaleschuld, nietNagekomenInt);
  if (result.resultCode ==1) {
    resultString = createStringFromStringT(result.resultMessage);
  } else {
    resultString = "sociaal_minimum:1300";
  }
  return resultString;
}

String evalToeslagenWetBestaansMinimum(String bsn, bool heeftPartner, bool heeftWoningDeler, String partnerBsn, String woningDelerBsn, int leeftijd, int leeftijdPartner, int leeftijdWoningDeler) {
  final Pointer<String_t> bsnString = createStringPointer(bsn);
  final Pointer<String_t> bsnPartnerString = createStringPointer(partnerBsn);
  final Pointer<String_t> bsnWoningDelerString = createStringPointer(woningDelerBsn);
  int heeftPartnerInt = 0;
  if (heeftPartner) {
    heeftPartnerInt = 1;
  }
  int heeftWoningDelerInt = 0;
  if (heeftWoningDeler) {
    heeftWoningDelerInt = 1;
  }
  var resultString = "";
      Machine_law_Result_t result =_bindings.EvaluateToeslagenWetBestaansMinimum(bsnString.ref, heeftPartnerInt, heeftWoningDelerInt, bsnPartnerString.ref, bsnWoningDelerString.ref, leeftijd, leeftijdPartner, leeftijdWoningDeler);
  var resultCode = result.resultCode;
  if (resultCode == 1) {
    resultString = createStringFromStringT(result.resultMessage);
  } else {
    resultString = "no result";
  }
  return resultString;
}


String createStringFromStringT(String_t string) {
  final Uint8List resultMessage = string.string.asTypedList(string.length);
  var resultString = utf8.decode(resultMessage);
  return resultString;
}

Pointer<String_t> createStringPointer(String string){
  final Pointer<Uint8> stringPointer = calloc<Uint8>(string.length);
  stringPointer.asTypedList(string.length).setAll(0, string.codeUnits);
  final Pointer<String_t> cString = calloc<String_t>();
  cString.ref.length = string.length;
  cString.ref.string = stringPointer;
  return cString;
}

void freeString(Pointer<String_t> pointer) {
  calloc.free(pointer.ref.string);
  calloc.free(pointer);
}

Pointer<Machine_law_Params_t> createParametersPointer(String parameters) {
  final Pointer<Machine_law_Params_t> cParams = calloc<Machine_law_Params_t>();
  cParams.ref.params = createStringPointer(parameters).ref;
  return cParams;
}


const String _libName = 'machine_law';

/// The dynamic library in which the symbols for [MachineLawBindings] can be found.
final DynamicLibrary _dylib = () {
  if (Platform.isMacOS || Platform.isIOS) {
    return DynamicLibrary.open('$_libName.framework/$_libName');
  }
  if (Platform.isLinux) {
    return DynamicLibrary.open('$_libName.so');
  }
  if (Platform.isAndroid ){
    return DynamicLibrary.open('machine_law.so');
  }
  if (Platform.isWindows) {
    return DynamicLibrary.open('$_libName.dll');
  }
  throw UnsupportedError('Unknown platform: ${Platform.operatingSystem}');
}();

/// The bindings to the native functions in [_dylib].
final MachineLawBindings _bindings = MachineLawBindings(_dylib);
