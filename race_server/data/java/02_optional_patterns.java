// Topic: Optional Chaining

import java.util.Optional;

String city = Optional.ofNullable(getUser())
    .map(User::getAddress)
    .map(Address::getCity)
    .filter(c -> !c.isEmpty())
    .orElse("Unknown");

Optional<Integer> parsed = Optional.ofNullable(input)
    .filter(s -> s.matches("\\d+"))
    .map(Integer::parseInt);

parsed.ifPresentOrElse(
    val -> System.out.println("Got: " + val),
    () -> System.out.println("No valid number")
);