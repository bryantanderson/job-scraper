import express, { Request, Response } from "express";
import CandidateDto from "../models/candidateDto";
import CandidateService from "../services/candidate";
import {
    HTTP_CREATED,
    HTTP_INTERNAL_SERVER_ERROR,
    HTTP_NO_CONTENT,
    HTTP_NOT_FOUND,
    HTTP_OK,
    INTERNAL_SERVER_ERROR,
} from "../util/constants";
import { wrapText } from "../util/util";

export const router = express.Router();
router.use(express.json());

router.post("/", async (req: Request, res: Response) => {
    try {
        const service = new CandidateService();
        const dto = req.body as CandidateDto;
        const candidate = await service.createCandidate(dto);
        res.status(HTTP_CREATED).send(candidate);
    } catch (error) {
        console.log(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(INTERNAL_SERVER_ERROR);
    }
});

router.get("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;
    try {
        const service = new CandidateService();
        const candidate = await service.getCandidate(id);

        if (candidate) {
            res.status(HTTP_OK).send(candidate);
            return;
        }
        res.status(HTTP_NOT_FOUND).send(
            wrapText(`Candidate with id ${id} does not exist`)
        );
    } catch (error) {
        console.log(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(INTERNAL_SERVER_ERROR);
    }
});

router.put("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;
    try {
        const service = new CandidateService();
        const dto = req.body as CandidateDto;
        const updatedCount = await service.updateCandidate(id, dto);

        if (updatedCount) {
            res.status(HTTP_OK).send(
                `Successfully updated ${updatedCount} candidates`
            );
            return;
        }
        res.status(HTTP_NOT_FOUND).send(
            wrapText(`Candidate with id ${id} does not exist`)
        );
    } catch (error) {
        console.log(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(INTERNAL_SERVER_ERROR);
    }
});

router.delete("/:id", async (req: Request, res: Response) => {
    const id = req?.params?.id;
    try {
        const service = new CandidateService();
        const deletedCount = await service.deleteCandidate(id);

        if (deletedCount) {
            res.status(HTTP_NO_CONTENT);
            return;
        }
        res.status(HTTP_NOT_FOUND).send(
            wrapText(`Candidate with id ${id} does not exist`)
        );
    } catch (error) {
        console.log(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(INTERNAL_SERVER_ERROR);
    }
});
