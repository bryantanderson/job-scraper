export default class Experience {
    constructor(
        public title: string,
        public company: string,
        public description: string,
        public startDate: Date,
        public endDate?: Date
    ) {}
}
