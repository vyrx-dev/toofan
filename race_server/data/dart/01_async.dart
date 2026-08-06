// Topic: Async and Futures
import 'dart:async';

Future<String> fetchData(String url) async {
  await Future.delayed(Duration(seconds: 1));
  if (url.isEmpty) {
    throw Exception('Empty URL');
  }
  return '{"status": "ok"}';
}

Stream<int> countDown(int from) async* {
  for (var i = from; i >= 0; i--) {
    await Future.delayed(Duration(milliseconds: 100));
    yield i;
  }
}

void main() async {
  try {
    final data = await fetchData('/api');
    print(data);
  } catch (e) {
    print('Error: $e');
  }

  await for (final n in countDown(3)) {
    print(n);
  }
}
