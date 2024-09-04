import Qualification from "./qualification";
import Responsibility from "./responsibility";
import { ObjectId } from "mongodb";

export default class Job {
    constructor(
        public id: ObjectId,
        public title: string,
        public company: string,
        public description: string,
        public responsibilities: Responsibility[],
        public qualifications: Qualification[],
        public location: string,
        public locationType: string,
        public yearsOfExperience?: number,
        public elasticId?: string
    ) {}
}
