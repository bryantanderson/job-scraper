import { ObjectId } from "mongodb";
import Candidate from "../models/candidate";
import CandidateDto from "../models/candidateDto";
import { collections } from "./mongo";
import { copyFields } from "../util/util";

export default class CandidateService {
    constructor() {}

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
            return candidate;
        }
        return null;
    }

    public async getCandidate(id: string): Promise<Candidate | null> {
        const query = { _id: new ObjectId(id) };
        const res = await collections.candidates?.findOne(query);

        if (res) {
            return res as Candidate;
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
}
