import express, { Request, Response } from "express";
import {
    HTTP_INTERNAL_SERVER_ERROR,
    HTTP_NOT_AUTHORIZED,
    HTTP_OK,
} from "../util/constants";
import LoginDto from "../models/loginDto";
import LoginService from "../services/login";

export const router = express.Router();
router.use(express.json());

router.post("/login", async (req: Request, res: Response) => {
    try {
        const service = new LoginService();
        const dto = req.body as LoginDto;
        const accessToken = await service.Login(dto);

        if (!accessToken) {
            res.status(HTTP_NOT_AUTHORIZED);
        }
        res.status(HTTP_OK).send(accessToken);
    } catch (error) {
        console.error(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(error);
    }
});

router.post("/sign-in", async (req: Request, res: Response) => {
    try {
        const service = new LoginService();
        const dto = req.body as LoginDto;
        const accessToken = await service.SignUp(dto);

        if (!accessToken) {
            res.status(HTTP_INTERNAL_SERVER_ERROR);
        }
        res.status(HTTP_OK).send(accessToken);
    } catch (error) {
        console.error(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(error);
    }
});
