// Topic: CompletableFuture

import java.util.concurrent.CompletableFuture;

CompletableFuture<String> userFuture = CompletableFuture
    .supplyAsync(() -> fetchUser(userId))
    .thenApply(user -> user.getName())
    .exceptionally(ex -> "anonymous");

CompletableFuture<Void> combined = CompletableFuture.allOf(
    fetchProfile(id),
    fetchOrders(id),
    fetchPreferences(id)
);

combined.thenRun(() -> System.out.println("All loaded"));