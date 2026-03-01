import 'dart:io';

import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';


class PathProvide extends StatelessWidget {
  const PathProvide({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Path Provider',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const PathProvidePage(title: 'Path Provider'),
    );
  }
}

class PathProvidePage extends StatefulWidget {
  const PathProvidePage({super.key, required this.title});
  final String title;

  @override
  _PathProvidePageState createState() => _PathProvidePageState();
}

class _PathProvidePageState extends State<PathProvidePage> {
  Future<Directory?>? _tempDirectory;
  Future<Directory?>? _appSupportDirectory;
  Future<Directory?>? _appLibraryDirectory;
  Future<Directory?>? _appDocumentsDirectory;
  Future<Directory?>? _externalDocumentsDirectory;
  Future<List<Directory>?>? _externalStorageDirectories;
  Future<List<Directory>?>? _externalCacheDirectories;

  void _requestTempDirectory() {
    setState(() {
      _tempDirectory = getTemporaryDirectory();
    });
  }

  Widget _buildDirectory(
      BuildContext context, AsyncSnapshot<Directory?> snapshot) {
    Text text = const Text('');
    if (snapshot.connectionState == ConnectionState.done) {
      if (snapshot.hasError) {
        text = Text('Error: ${snapshot.error}');
      } else if (snapshot.hasData) {
        text = Text('path: ${snapshot.data!.path}');
      } else {
        text = const Text('path unavailable');
      }
    }
    return Padding(padding: const EdgeInsets.all(16.0), child: text);
  }

  Widget _buildDirectories(
      BuildContext context, AsyncSnapshot<List<Directory>?> snapshot) {
    Text text = const Text('');
    if (snapshot.connectionState == ConnectionState.done) {
      if (snapshot.hasError) {
        text = Text('Error: ${snapshot.error}');
      } else if (snapshot.hasData) {
        final String combined =
        snapshot.data!.map((Directory d) => d.path).join(', ');
        text = Text('paths: $combined');
      } else {
        text = const Text('path unavailable');
      }
    }
    return Padding(padding: const EdgeInsets.all(16.0), child: text);
  }

  void _requestAppDocumentsDirectory() {
    setState(() {
      _appDocumentsDirectory = getApplicationDocumentsDirectory();
    });
  }

  void _requestAppSupportDirectory() {
    setState(() {
      _appSupportDirectory = getApplicationSupportDirectory();
    });
  }

  void _requestAppLibraryDirectory() {
    setState(() {
      _appLibraryDirectory = getLibraryDirectory();
    });
  }

  void _requestExternalStorageDirectory() {
    setState(() {
      _externalDocumentsDirectory = getExternalStorageDirectory();
    });
  }

  void _requestExternalStorageDirectories(StorageDirectory type) {
    setState(() {
      _externalStorageDirectories = getExternalStorageDirectories(type: type);
    });
  }

  void _requestExternalCacheDirectories() {
    setState(() {
      _externalCacheDirectories = getExternalCacheDirectories();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: ListView(
          children: <Widget>[
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton(
                onPressed: _requestTempDirectory,
                child: const Text('Get Temporary Directory'),
              ),
            ),
            FutureBuilder<Directory?>(
                future: _tempDirectory, builder: _buildDirectory),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton(
                onPressed: _requestAppDocumentsDirectory,
                child: const Text('Get Application Documents Directory'),
              ),
            ),
            FutureBuilder<Directory?>(
                future: _appDocumentsDirectory, builder: _buildDirectory),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton(
                onPressed: _requestAppSupportDirectory,
                child: const Text('Get Application Support Directory'),
              ),
            ),
            FutureBuilder<Directory?>(
                future: _appSupportDirectory, builder: _buildDirectory),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton(
                onPressed: _requestAppLibraryDirectory,
                child: const Text('Get Application Library Directory'),
              ),
            ),
            FutureBuilder<Directory?>(
                future: _appLibraryDirectory, builder: _buildDirectory),
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton(
                onPressed:
                Platform.isIOS ? null : _requestExternalStorageDirectory,
                child: Text(Platform.isIOS
                    ? 'External directories are unavailable on iOS'
                    : 'Get External Storage Directory'),
              ),
            ),
            FutureBuilder<Directory?>(
                future: _externalDocumentsDirectory, builder: _buildDirectory),
            Column(children: <Widget>[
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: ElevatedButton(
                  onPressed: Platform.isIOS
                      ? null
                      : () {
                    _requestExternalStorageDirectories(
                      StorageDirectory.music,
                    );
                  },
                  child: Text(Platform.isIOS
                      ? 'External directories are unavailable on iOS'
                      : 'Get External Storage Directories'),
                ),
              ),
            ]),
            FutureBuilder<List<Directory>?>(
                future: _externalStorageDirectories,
                builder: _buildDirectories),
            Column(children: <Widget>[
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: ElevatedButton(
                  onPressed:
                  Platform.isIOS ? null : _requestExternalCacheDirectories,
                  child: Text(Platform.isIOS
                      ? 'External directories are unavailable on iOS'
                      : 'Get External Cache Directories'),
                ),
              ),
            ]),
            FutureBuilder<List<Directory>?>(
                future: _externalCacheDirectories, builder: _buildDirectories),
          ],
        ),
      ),
    );
  }
}