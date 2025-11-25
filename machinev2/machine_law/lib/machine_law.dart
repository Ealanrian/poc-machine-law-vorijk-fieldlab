
import 'dart:async';
import 'dart:ffi';
import 'dart:io';
import 'dart:isolate';


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
  final Pointer<Machine_law_Params_t> parameters = createParametersPointer();


  Machine_law_Result_t result =_bindings.Evaluate(service.ref, law.ref, parameters, referenceDate, effectiveDate, overwriteInput, requestedOutput, approved);


  freeString(service);
  freeString(law);
  return 0;
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

Pointer<Machine_law_Params_t> createParametersPointer() {
  final Pointer<Machine_law_Params_t> cParams = calloc<Machine_law_Params_t>();
  cParams.ref.numParams = 0;
  cParams.ref.params = calloc<Machine_law_param_t>(10);
  final Pointer<Machine_law_param_t> cParam = calloc<Machine_law_param_t>(10);
  cParam.asTypedList()

  return cParams;
}

/// A longer lived native function, which occupies the thread calling it.
///
/// Do not call these kind of native functions in the main isolate. They will
/// block Dart execution. This will cause dropped frames in Flutter applications.
/// Instead, call these native functions on a separate isolate.
///
/// Modify this to suit your own use case. Example use cases:
///
/// 1. Reuse a single isolate for various different kinds of requests.
/// 2. Use multiple helper isolates for parallel execution.
/*Future<int> sumAsync(int a, int b) async {
  final SendPort helperIsolateSendPort = await _helperIsolateSendPort;
  final int requestId = _nextSumRequestId++;
  final _SumRequest request = _SumRequest(requestId, a, b);
  final Completer<int> completer = Completer<int>();
  _sumRequests[requestId] = completer;
  helperIsolateSendPort.send(request);
  return completer.future;
}*/

const String _libName = 'machine_law';

/// The dynamic library in which the symbols for [MachineLawBindings] can be found.
final DynamicLibrary _dylib = () {
  if (Platform.isMacOS || Platform.isIOS) {
    return DynamicLibrary.open('$_libName.framework/$_libName');
  }
  if (Platform.isAndroid || Platform.isLinux) {
    return DynamicLibrary.open('$_libName.so');
  }
  if (Platform.isWindows) {
    return DynamicLibrary.open('$_libName.dll');
  }
  throw UnsupportedError('Unknown platform: ${Platform.operatingSystem}');
}();

/// The bindings to the native functions in [_dylib].
final MachineLawBindings _bindings = MachineLawBindings(_dylib);


/// A request to compute `sum`.
///
/// Typically sent from one isolate to another.
class _SumRequest {
  final int id;
  final int a;
  final int b;

  const _SumRequest(this.id, this.a, this.b);
}

/// A response with the result of `sum`.
///
/// Typically sent from one isolate to another.
class _SumResponse {
  final int id;
  final int result;

  const _SumResponse(this.id, this.result);
}

/// Counter to identify [_SumRequest]s and [_SumResponse]s.
int _nextSumRequestId = 0;

/// Mapping from [_SumRequest] `id`s to the completers corresponding to the correct future of the pending request.
final Map<int, Completer<int>> _sumRequests = <int, Completer<int>>{};

/// The SendPort belonging to the helper isolate.
///
/*
Future<SendPort> _helperIsolateSendPort = () async {
  // The helper isolate is going to send us back a SendPort, which we want to
  // wait for.
  final Completer<SendPort> completer = Completer<SendPort>();

  // Receive port on the main isolate to receive messages from the helper.
  // We receive two types of messages:
  // 1. A port to send messages on.
  // 2. Responses to requests we sent.
  final ReceivePort receivePort = ReceivePort()
    ..listen((dynamic data) {
      if (data is SendPort) {
        // The helper isolate sent us the port on which we can sent it requests.
        completer.complete(data);
        return;
      }
      if (data is _SumResponse) {
        // The helper isolate sent us a response to a request we sent.
        final Completer<int> completer = _sumRequests[data.id]!;
        _sumRequests.remove(data.id);
        completer.complete(data.result);
        return;
      }
      throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
    });

  // Start the helper isolate.
  await Isolate.spawn((SendPort sendPort) async {
    final ReceivePort helperReceivePort = ReceivePort()
      ..listen((dynamic data) {
        // On the helper isolate listen to requests and respond to them.
        if (data is _SumRequest) {
          final int result = _bindings.sum_long_running(data.a, data.b);
          final _SumResponse response = _SumResponse(data.id, result);
          sendPort.send(response);
          return;
        }
        throw UnsupportedError('Unsupported message type: ${data.runtimeType}');
      });

    // Send the port to the main isolate on which we can receive requests.
    sendPort.send(helperReceivePort.sendPort);
  }, receivePort.sendPort);

  // Wait until the helper isolate has sent us back the SendPort on which we
  // can start sending requests.
  return completer.future;
}();

 */
