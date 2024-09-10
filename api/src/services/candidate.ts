import { ObjectId } from "mongodb";
import Candidate from "../models/candidate";
import CandidateDto from "../models/candidateDto";
import { collections } from "./database";
import { copyFields } from "../util/util";
import { getCache } from "./cache";

export default class CandidateService {
    private cacheExpiry: number = 3600;

    public async createCandidate(dto: CandidateDto): Promise<Candidate | null> {
        const candidate: Candidate = {
            _id: new ObjectId(),
            email: dto.email,
            firstName: dto.firstName,
            lastName: dto.lastName,
            contactNumber: dto.contactNumber,
            education: dto.education,
            experiences: dto.experiences,
            skills: dto.skills.map((skill) => ({ description: skill })),
            location: dto.location,
            summary: dto.summary,
        };
        const res = await collections.candidates?.insertOne(candidate);

        if (res?.acknowledged) {
            await this.cacheCandidate(candidate);
            return candidate;
        }
        return null;
    }

    public async getCandidate(id: string): Promise<Candidate | null> {
        const cachedCandidate = await this.getCachedCandidate(id);

        if (cachedCandidate) {
            return cachedCandidate;
        }

        const query = { _id: new ObjectId(id) };
        const res = await collections.candidates?.findOne(query);

        if (res) {
            const candidate = res as Candidate;
            await this.cacheCandidate(candidate);
            return candidate;
        }
        return null;
    }

    public async updateCandidate(
        id: string,
        dto: CandidateDto
    ): Promise<number | null> {
        const existingCandidate = await this.getCandidate(id);

        if (!existingCandidate) {
            throw new Error("Candidate does not exist");
        }

        copyFields(existingCandidate, dto);

        const query = { _id: new ObjectId(id) };
        const res = await collections.candidates?.updateOne(query, {
            $set: existingCandidate,
        });

        if (res) {
            await this.invalidateCachedCandidate(id);
            return res.upsertedCount;
        }
        return null;
    }

    public async deleteCandidate(id: string): Promise<number | null> {
        const query = { _id: new ObjectId(id) };
        const res = await collections.candidates?.deleteOne(query);

        if (res) {
            return res.deletedCount;
        }
        return null;
    }

    private async cacheCandidate(candidate: Candidate): Promise<void> {
        const cache = await getCache();
        await cache.setEx(
            this.createCacheKey(candidate._id.toString()),
            this.cacheExpiry,
            JSON.stringify(candidate)
        );
    }

    private async getCachedCandidate(id: string): Promise<Candidate | null> {
        const cache = await getCache();
        const stringCandidate = await cache.get(this.createCacheKey(id));

        if (stringCandidate) {
            return JSON.parse(stringCandidate) as Candidate;
        }
        return null;
    }

    private async invalidateCachedCandidate(id: string): Promise<void> {
        const cache = await getCache();
        await cache.del(this.createCacheKey(id));
    }

    private createCacheKey(id: string): string {
        return `candidates-${id}`;
    }
}
