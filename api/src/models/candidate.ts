import { ObjectId } from "mongodb";
import Education from "./education";
import Experience from "./experience";
import Skill from "./skill";

type Candidate = {
    _id: ObjectId;
    email: string;
    firstName: string;
    lastName: string;
    contactNumber: string;
    education: Education[];
    experiences: Experience[];
    skills: Skill[];
    summary?: string;
    location?: string;
};

export default Candidate;
