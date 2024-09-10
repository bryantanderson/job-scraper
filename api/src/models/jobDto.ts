type JobDto = {
    title: string;
    company: string;
    description: string;
    responsibilities: string[];
    qualifications: string[];
    location: string;
    locationType: string;
    yearsOfExperience?: number;
};

export default JobDto;
