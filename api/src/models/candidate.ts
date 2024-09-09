import Education from "./education";
import Experience from "./experience";
import Skill from "./skill";

export default interface Candidate {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    contactNumber: string;
    education: Education[];
    experiences: Experience[];
    skills: Skill[];
    summary?: string;
    location?: string;
}
