import {
    HTTP_BAD_REQUEST,
    HTTP_CREATED,
    HTTP_INTERNAL_SERVER_ERROR,
    HTTP_NOT_FOUND,
    HTTP_OK,
} from "../constants";
import http from "http";
import {
    wrapResponse,
    wrapText,
    getRequestBody,
    getRequestPathname,
} from "../util/util";
import { ObjectId } from "mongodb";
import { collections } from "../repository/mongo";
import Job from "../models/job";

const postJobHandler = async (
    req: http.IncomingMessage,
    res: http.ServerResponse
) => {
    try {
        const job: Job = JSON.parse(getRequestBody(req)) as Job;
        const insertRes = await collections.jobs?.insertOne(job);
        if (insertRes) {
            wrapResponse(res, HTTP_CREATED);
        } else {
            wrapResponse(res, HTTP_INTERNAL_SERVER_ERROR);
        }
    } catch (error) {
        console.error(error);
        wrapResponse(res, HTTP_INTERNAL_SERVER_ERROR);
    }
};

const putJobHandler = async (
    req: http.IncomingMessage,
    res: http.ServerResponse
) => {
    const id = getJobId(req);

    if (id === null) {
        wrapResponse(res, HTTP_BAD_REQUEST);
        return;
    }

    try {
        const updatedJob: Job = JSON.parse(getRequestBody(req)) as Job;
        const query = { _id: new ObjectId(id) };

        const result = collections?.jobs?.updateOne(query, {
            $set: updatedJob,
        });

        if (result) {
            wrapResponse(res, HTTP_OK);
        } else {
            wrapResponse(res, HTTP_INTERNAL_SERVER_ERROR);
        }
    } catch (error) {
        console.error(error);
        wrapResponse(res, HTTP_INTERNAL_SERVER_ERROR);
    }
};

const getJobHandler = async (
    req: http.IncomingMessage,
    res: http.ServerResponse
) => {
    const id = getJobId(req);

    if (id === null) {
        wrapResponse(res, HTTP_BAD_REQUEST);
        return;
    }

    try {
        const query = { _id: new ObjectId(id) };
        const job = (await collections.jobs?.findOne(query)) as unknown as Job;
        if (job) {
            wrapResponse(res, HTTP_OK);
        } else {
            wrapResponse(res, HTTP_NOT_FOUND, wrapText("Job does not exist"));
        }
    } catch (error) {
        console.error(error);
        wrapResponse(res, HTTP_INTERNAL_SERVER_ERROR);
    }
};

const getJobId = (req: http.IncomingMessage) => {
    const requestPathname: string | null = getRequestPathname(req);

    if (requestPathname === null) {
        return null;
    }

    const pathParts: string[] = requestPathname.split("/");

    if (pathParts.length === 0) {
        return null;
    }

    return pathParts[-1];
};

const router: Map<
    string,
    (req: http.IncomingMessage, res: http.ServerResponse) => Promise<void>
> = new Map([
    ["GET", getJobHandler],
    ["POST", postJobHandler],
    ["PUT", putJobHandler],
]);

export default router;
