import { ObjectId, UUID } from "mongodb";
import Job from "../models/job";
import JobDto from "../models/jobDto";
import Responsibility from "../models/responsibility";
import Qualification from "../models/qualification";
import { collections } from "../repository/mongo";

export default class JobService {
    constructor() {}

    public async createJob(dto: JobDto): Promise<Job | null> {
        const id = new ObjectId(UUID.toString());
        const job: Job = new Job(
            id,
            dto.title,
            dto.company,
            dto.description,
            dto.responsibilities.map((resp) => new Responsibility(resp)),
            dto.qualifications.map((qual) => new Qualification(qual)),
            dto.location,
            dto.locationType,
            dto.yearsOfExperience
        );
        const res = await collections.jobs?.insertOne(job);
        const insertedId = res?.insertedId.toString();

        if (insertedId) {
            return job;
        }
        return null;
    }

    public async getJob(id: string): Promise<Job | null> {
        const query = { _id: new ObjectId(id) };
        const job = (await collections.jobs?.findOne(query)) as unknown as Job;

        if (job) {
            return job;
        }
        return null;
    }

    public async putJob(id: string, dto: JobDto): Promise<string | null> {
        const existingJob = await this.getJob(id);

        if (!existingJob) {
            return null;
        }

        this.updateJob(existingJob, dto);

        const query = { _id: new ObjectId(id) };
        const result = await collections?.jobs?.updateOne(query, {
            $set: existingJob,
        });

        if (result?.upsertedId) {
            return result.upsertedId.toString();
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

    private updateJob(job: Job, dto: JobDto): void {
        job.title = dto.title;
        job.company = dto.company;
        job.description = dto.description;
        job.location = dto.location;
        job.locationType = dto.locationType;
        job.yearsOfExperience = dto.yearsOfExperience;
    }
}
