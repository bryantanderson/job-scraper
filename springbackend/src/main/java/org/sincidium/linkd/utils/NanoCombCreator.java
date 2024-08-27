package org.sincidium.linkd.utils;

import java.security.SecureRandom;
import java.util.UUID;

// https://stackoverflow.com/questions/63015186/make-ulid-lexicographic-ordering-more-sensitive-to-time
public final class NanoCombCreator {

    private long prevTime = 0;
    private long prevNano = 0;

    private static final long ONE_MILLION_NANOSECONDS = 1_000_000L;

    private static final SecureRandom SECURE_RANDOM = new SecureRandom();

    /**
     * Returns a time component in nanoseconds.
     *
     * It uses `System.currentTimeMillis()` to get the system time in milliseconds
     * accuracy. The nanoseconds resolution is simulated by calling
     * `System.nanoTime()` between subsequent calls within the same millisecond.
     * It's not precise, but it provides some monotonicity to the values generates.
     *
     * @return the current time in nanoseconds
     */
    private synchronized long getTimeComponent() {

        final long time = System.currentTimeMillis();
        final long nano = System.nanoTime();
        final long elapsed; // nanoseconds since last call

        if (time == prevTime) {
            elapsed = (nano - prevNano);
            if (elapsed > ONE_MILLION_NANOSECONDS) {
                try {
                    // make the clock to catch up
                    Thread.sleep(1);
                } catch (InterruptedException ex) {
                    System.err.println("something went wrong...");
                }
            }
        } else {
            prevTime = time;
            prevNano = nano;
            elapsed = 0;
        }

        return (time * ONE_MILLION_NANOSECONDS) + elapsed;
    }

    /**
     * Returns the random component using a secure random generator.
     *
     * @return a random value.
     */
    private synchronized long getRandomComponent() {
        return SECURE_RANDOM.nextLong();
    }

    /**
     * Returns a Nano COMB.
     *
     * A Nano COMB is inspired on ULID and COMB generators.
     *
     * It is composed of 64 bits for time and 64 for random bits.
     *
     * @return a UUID
     */
    public synchronized UUID create() {

        final long timeBits = getTimeComponent();
        final long randomBits = getRandomComponent();

        return new UUID(timeBits, randomBits);
    }

}
