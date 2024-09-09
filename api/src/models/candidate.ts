import Education from "./education";
import Experience from "./experience";
import Skill from "./skill";

export default interface Candidate {
    id: string;
    education: Education[];
    experiences: Experience[];
    skills: Skill[];
    summary?: string;
    location?: string;
}
