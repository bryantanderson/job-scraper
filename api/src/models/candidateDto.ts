import Education from "./education";
import Experience from "./experience";

export default interface CandidateDto {
    email: string;
    firstName: string;
    lastName: string;
    contactNumber: string;
    education: Education[];
    experiences: Experience[];
    skills: string[];
    summary?: string;
    location?: string;
}
