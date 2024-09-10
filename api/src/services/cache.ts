import type { RedisClientType } from "redis";
import { createClient } from "redis";

let isReady: boolean;
let redisClient: RedisClientType;

const url = process.env.REDIS_CONNECTION_STRING || "redis://localhost:6379";

async function getCache(): Promise<RedisClientType> {
    if (isReady) {
        return redisClient;
    }
    redisClient = createClient({
        url,
    });
    redisClient.on("error", (error) => {
        console.error(`Redis client error:`, error);
    });
    redisClient.on("connect", () => {
        console.log("Redis connected");
    });
    redisClient.on("reconnecting", () => {
        console.log("Redis reconnecting");
    });
    redisClient.on("ready", () => {
        console.log("Redis ready");
    });
    await redisClient.connect();
    return redisClient;
}

getCache()
    .then((connection) => {
        redisClient = connection;
    })
    .catch((error) => {
        console.log(`Failed to connect to Redis with error: ${error}`);
    });

export { getCache };
