// Topic: Stream Pipeline

import java.util.*;
import java.util.stream.*;

List<String> result = users.stream()
    .filter(u -> u.age() > 18)
    .map(User::name)
    .collect(Collectors.toList());
