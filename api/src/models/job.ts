import { ObjectId } from "mongodb";
import Qualification from "./qualification";
import Responsibility from "./responsibility";

type Job = {
    _id: ObjectId;
    title: string;
    company: string;
    description: string;
    responsibilities: Responsibility[];
    qualifications: Qualification[];
    location: string;
    locationType: string;
    yearsOfExperience?: number;
};

export default Job;
