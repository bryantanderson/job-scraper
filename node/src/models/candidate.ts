import Education from "./education";
import Experience from "./experience";
import Skill from "./skill";

export default class Candidate {
    constructor(
        public id: string,
        public education: Education[],
        public experiences: Experience[],
        public skills: Skill[],
        public summary?: string,
        public location?: string
    ) {}
}
