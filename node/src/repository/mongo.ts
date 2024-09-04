import * as mongoDB from "mongodb";

export const collections: {
    jobs?: mongoDB.Collection;
    users?: mongoDB.Collection;
    candidates?: mongoDB.Collection;
} = {};

export async function connectToMongoDB() {
    const databaseName: string = process.env.MONGO_DATABASE_NAME || "linkd";

    const client: mongoDB.MongoClient = new mongoDB.MongoClient(
        process.env.MONGO_CONNECTION_STRING || ""
    );

    await client.connect();

    const jobsCollectionName: string =
        process.env.MONGO_JOBS_COLLECTION_NAME || "jobs";
    const candidatesCollectionName: string =
        process.env.MONGO_CANDIDATES_COLLECTION_NAME || "candidates";
    const usersCollectionName: string =
        process.env.MONGO_USERS_COLLECTION_NAME || "users";

    const db: mongoDB.Db = client.db(databaseName);

    collections.jobs = db.collection(jobsCollectionName);
    collections.users = db.collection(usersCollectionName);
    collections.candidates = db.collection(candidatesCollectionName);
}
