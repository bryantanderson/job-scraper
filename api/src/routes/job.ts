import { HTTP_STATUS } from "../util/constants";
import express, { Request, Response } from "express";
import JobDto from "../models/jobDto";
import JobService from "../services/job";

/**
 * @swagger
 * components:
 *   securitySchemes:
 *     BearerAuth:
 *       type: http
 *       scheme: bearer
 *   schemas:
 *     Responsibility:
 *       type: object
 *       properties:
 *         description:
 *           type: string
 *           description: The responsibility description

 *     Qualification:
 *       type: object
 *       properties:
 *         description:
 *           type: string
 *           description: The qualification description
 * 
 *     JobDto:
 *       type: object
 *       required:
 *         - title
 *         - company
 *         - description
 *         - responsibilities
 *         - qualifications
 *         - location
 *         - locationType
 *       properties:
 *         title:
 *           type: string
 *           description: The title of the job
 *         company:
 *           type: string
 *           description: The company offering the job
 *         description:
 *           type: string
 *           description: A detailed description of the job
 *         responsibilities:
 *           type: array
 *           items:
 *             type: string
 *           description: List of responsibilities for the job
 *         qualifications:
 *           type: array
 *           items:
 *             type: string
 *           description: List of qualifications for the job
 *         location:
 *           type: string
 *           description: The location of the job
 *         locationType:
 *           type: string
 *           description: The type of location (e.g., remote, onsite)
 *         yearsOfExperience:
 *           type: number
 *           description: Optional years of experience required for the job
 *
 *     Job:
 *       type: object
 *       required:
 *         - id
 *         - title
 *         - company
 *         - description
 *         - responsibilities
 *         - qualifications
 *         - location
 *         - locationType
 *       properties:
 *         id:
 *           type: string
 *           description: The auto-generated id of the job
 *         title:
 *           type: string
 *           description: The title of the job
 *         company:
 *           type: string
 *           description: The company offering the job
 *         description:
 *           type: string
 *           description: A detailed description of the job
 *         responsibilities:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/Responsibility'
 *           description: List of responsibilities for the job
 *         qualifications:
 *           type: array
 *           items:
 *             $ref: '#/components/schemas/Qualification'
 *           description: List of qualifications for the job
 *         location:
 *           type: string
 *           description: The location of the job
 *         locationType:
 *           type: string
 *           description: The type of location (e.g., remote, onsite)
 *         yearsOfExperience:
 *           type: number
 *           description: Optional years of experience required for the job
 *         elasticId:
 *           type: string
 *           description: Optional Elasticsearch ID for the job
 * 
 *       example:
 *         id: 60d0fe4f5311236168a109ca
 *         title: Backend Developer
 *         company: Tech Corp
 *         description: Develop RESTful APIs using Node.js and Express.
 *         responsibilities:
 *           - description: Design and implement RESTful APIs using Node.js
 *         qualifications:
 *           - description: JavaScript
 *         location: New York, NY
 *         locationType: Remote
 *         yearsOfExperience: 3
 *         elasticId: abc123xyz
 * 
 */

export const router = express.Router();
router.use(express.json());

/**
 * @swagger
 * tags:
 *   name: Jobs
 * /jobs:
 *   post:
 *     security:
 *       - BearerAuth: []
 *     summary: Create a new job
 *     tags: [Jobs]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/JobDto'
 *     responses:
 *       201:
 *         description: Job successfully created
 *       500:
 *         description: Internal server error
 */
router.post("/", async (req: Request, res: Response) => {
    try {
        const service: JobService = new JobService();
        const jobDto: JobDto = req.body as JobDto;
        const job = await service.createJob(jobDto);

        if (job) {
            res.status(HTTP_STATUS.CREATED).send(job);
            return;
        }
        res.status(HTTP_STATUS.INTERNAL_SERVER_ERROR).send(
            "Unable to create job"
        );
    } catch (error) {
        console.error(error);
        res.status(HTTP_STATUS.INTERNAL_SERVER_ERROR).send(error);
    }
});

/**
 * @swagger
 * tags:
 *   name: Jobs
 * /jobs/{id}:
 *   get:
 *     security:
 *       - BearerAuth: []
 *     summary: Get a job by ID
 *     tags: [Jobs]
 *     parameters:
 *       - in: path
 *         name: id
 *         schema:
 *           type: string
 *         required: true
 *         description: The ID of the job
 *     responses:
 *       200:
 *         description: The job description by ID
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Job'
 *       404:
 *         description: Job not found
 *       500:
 *         description: Internal server error
 */
router.get("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;

    try {
        const service: JobService = new JobService();
        const job = await service.getJob(id);

        if (job) {
            res.status(HTTP_STATUS.OK).send(job);
            return;
        }
        res.status(HTTP_STATUS.NOT_FOUND).send("Job does not exist");
    } catch (error) {
        console.error(error);
        res.status(HTTP_STATUS.INTERNAL_SERVER_ERROR).send(error);
    }
});

/**
 * @swagger
 * tags:
 *   name: Jobs
 * /jobs/{id}:
 *   put:
 *     security:
 *       - BearerAuth: []
 *     summary: Update a job by ID
 *     tags: [Jobs]
 *     parameters:
 *       - in: path
 *         name: id
 *         schema:
 *           type: string
 *         required: true
 *         description: The ID of the job
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/JobDto'
 *     responses:
 *       200:
 *         description: Successfully updated job
 *       304:
 *         description: Job not modified
 *       500:
 *         description: Internal server error
 */
router.put("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;

    try {
        const service: JobService = new JobService();
        const dto: JobDto = req.body as JobDto;
        const updatedCount = await service.putJob(id, dto);

        if (updatedCount) {
            res.status(HTTP_STATUS.OK).send(
                `Successfully updated ${updatedCount} jobs`
            );
            return;
        }
        res.status(HTTP_STATUS.NOT_MODIFIED).send(
            `Job with id: ${id} was not updated`
        );
    } catch (error) {
        console.error(error);
        res.status(HTTP_STATUS.INTERNAL_SERVER_ERROR).send(error);
    }
});

/**
 * @swagger
 * tags:
 *   name: Jobs
 * /jobs/{id}:
 *   delete:
 *     security:
 *       - BearerAuth: []
 *     summary: Delete a job by ID
 *     tags: [Jobs]
 *     parameters:
 *       - in: path
 *         name: id
 *         schema:
 *           type: string
 *         required: true
 *         description: The ID of the job
 *     responses:
 *       204:
 *         description: Successfully deleted the job
 *       404:
 *         description: Job not found
 *       500:
 *         description: Internal server error
 */
router.delete("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;

    try {
        const service: JobService = new JobService();
        const deletedCount = await service.deleteJob(id);

        if (deletedCount === null || deletedCount === 0) {
            res.status(HTTP_STATUS.NOT_FOUND).send(
                `Job with id: ${id} does not exist`
            );
            return;
        }
        res.status(HTTP_STATUS.NO_CONTENT).send();
    } catch (error) {
        console.error(error);
        res.status(HTTP_STATUS.INTERNAL_SERVER_ERROR).send(error);
    }
});
