export default class JobDto {
    constructor(
        public title: string,
        public company: string,
        public description: string,
        public responsibilities: string[],
        public qualifications: string[],
        public location: string,
        public locationType: string,
        public yearsOfExperience?: number
    ) {}
}
