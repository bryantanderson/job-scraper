import { ObjectId } from "mongodb";
import Job from "../models/job";
import JobDto from "../models/jobDto";
import { collections } from "./database";
import { copyFields } from "../util/util";

export default class JobService {
    constructor() {}

    public async createJob(dto: JobDto): Promise<Job | null> {
        const job: Job = {
            _id: new ObjectId(),
            title: dto.title,
            company: dto.company,
            description: dto.description,
            responsibilities: dto.responsibilities.map((resp) => ({
                description: resp,
            })),
            qualifications: dto.qualifications.map((qual) => ({
                description: qual,
            })),
            location: dto.location,
            locationType: dto.locationType,
            yearsOfExperience: dto.yearsOfExperience,
        };
        const res = await collections.jobs?.insertOne(job);
        const insertedId = res?.insertedId.toString();

        if (insertedId) {
            return job;
        }
        return null;
    }

    public async getJob(id: string): Promise<Job | null> {
        const query = { _id: new ObjectId(id) };
        const res = await collections.jobs?.findOne(query);

        if (res) {
            return res as Job;
        }
        return null;
    }

    public async putJob(id: string, dto: JobDto): Promise<number | null> {
        const existingJob = await this.getJob(id);

        if (!existingJob) {
            return null;
        }

        copyFields(existingJob, dto);

        const query = { _id: new ObjectId(id) };
        const res = await collections?.jobs?.updateOne(query, {
            $set: existingJob,
        });

        if (res) {
            return res.upsertedCount;
        }
        return null;
    }

    public async deleteJob(id: string): Promise<number | null> {
        const query = { _id: new ObjectId(id) };
        const result = await collections?.jobs?.deleteOne(query);

        if (result) {
            return result.deletedCount;
        }
        return null;
    }
}
