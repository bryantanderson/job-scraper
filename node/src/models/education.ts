export default class Education {
    constructor(
        public title: string,
        public institute: string,
        public description: string,
        public startYear: number,
        public endYear?: number
    ) {}
}
