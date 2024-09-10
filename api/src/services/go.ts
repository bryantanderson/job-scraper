import http from "http";
import Candidate from "../models/candidate";
import Job from "../models/job";

const hostname = process.env.GO_SERVICE_HOSTNAME || "localhost";
const port = Number(process.env.GO_SERVICE_PORT) || 8090;

export default class GoService {
    public async createAssessment(candidate: Candidate, job: Job) {
        const postData = {
            candidate,
            job,
        };
        const options = {
            hostname: hostname,
            port: port,
            path: "/assessments",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Content-Length": Buffer.byteLength(JSON.stringify(postData)),
            },
        };

        const req = http.request(options, (res) => {
            console.log(res.statusCode);
        });

        req.write(postData);
        req.end();
    }
}
