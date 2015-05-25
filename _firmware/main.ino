/**
  Testing firmware for the Go implemention of the Particle Cloud API.
  See https://github.com/getniwa/spark for more details
*/

char *version = "test version 1.0";

void setup() {

    // Create the 'version' variable used by the automated tests
    Spark.variable("version", version, STRING);

    // And register a function for the same purpose
    Spark.function("toggleLamp", toggleLamp);

    // Another function that causes a publish event
    Spark.function("testPublish", testPublish);
}

void loop() {
}

// A callback for the sample function
int toggleLamp(String command) {
    return 1;
}

// A callback for the sample function
int testPublish(String command) {

    // Publish a known string to end the unit tests
    Spark.publish("test-successful", "The test worked!", 60, PRIVATE);

    return 1;
}

