// Topic: Concurrent Collections

import java.util.concurrent.*;
import java.util.concurrent.atomic.*;

ConcurrentHashMap<String, Integer> cache =
    new ConcurrentHashMap<>();
cache.computeIfAbsent("key", k -> expensiveLookup(k));

AtomicLong counter = new AtomicLong(0);
long value = counter.incrementAndGet();
