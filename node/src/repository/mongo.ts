import * as mongoDB from "mongodb";
import * as dotenv from "dotenv";

export const collections: { jobs?: mongoDB.Collection } = {};

export async function connectToMongoDB() {
    dotenv.config();

    const client: mongoDB.MongoClient = new mongoDB.MongoClient(
        process.env.MONGO_CONNECTION_STRING || ""
    );

    await client.connect();

    const db: mongoDB.Db = client.db(
        process.env.MONGO_DATABASE_NAME || "linkd"
    );

    const jobsCollection: mongoDB.Collection = db.collection(
        process.env.MONGO_JOBS_COLLECTION_NAME || "jobs"
    );

    collections.jobs = jobsCollection;
}
