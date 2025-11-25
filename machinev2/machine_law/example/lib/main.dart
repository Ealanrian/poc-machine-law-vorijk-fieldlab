import 'package:flutter/material.dart';
import 'dart:async';

import 'package:machine_law/machine_law.dart' as machine_law;

void main() {
  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  late int sumResult;
  late Future<int> sumAsyncResult;
  late int machineLawAlive;
  String? result = '';
  String _apiResult = '';

  @override
  void initState() {
    super.initState();
    machineLawAlive = machine_law.machineLawStandalone();
  }

  Future<void> setResultState(String? result) async {
    setState(() {
      _apiResult = result ?? 'Returned null';
    });
  }


  @override
  Widget build(BuildContext context) {
    const textStyle = TextStyle(fontSize: 25);
    const spacerSmall = SizedBox(height: 10);
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Native Packages'),
        ),
        body: SingleChildScrollView(
          child: Container(
            padding: const EdgeInsets.all(10),
            child: Column(
              children: [
                const Text(
                  'This calls a native function through FFI that is shipped as source in the package. '
                  'The native code is built as part of the Flutter Runner build.',
                  style: textStyle,
                  textAlign: TextAlign.center,
                ),
                spacerSmall,
                Text(
                  'machineLaw alive = $machineLawAlive',
                  style: textStyle,
                  textAlign: TextAlign.center,
                ),
                spacerSmall,
                ElevatedButton(
                  onPressed: () async {
                    bool? success;
                    try {
                      machine_law.startMachineLawEngine();
                      result = "machine started";
                    } catch (error) {
                      result = "It failed "+error.toString();
                    }

                    await setResultState(result);
                  },
                  child: const Text('start engine'),
                ),
                ElevatedButton(
                  onPressed: () async {
                    bool? success;
                    try {
                      machineLawAlive = machine_law.machineLawStandalone();
                      if (machineLawAlive == 1 ) {
                        result = "its alive!";
                      } else if (machineLawAlive == -1) {
                        result = "starting failed";
                      } else {
                        result ="push it";
                      }
                    } catch (error) {
                      result = "It failed "+error.toString();
                    }

                    await setResultState(result);
                  },
                  child: const Text('check state'),
                ),
                const Divider(),
                const Text(
                  'Result',
                  style: TextStyle(fontSize: 22),
                ),
                Text(_apiResult),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
