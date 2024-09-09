import express, { Request, Response } from "express";
import {
    HTTP_INTERNAL_SERVER_ERROR,
    HTTP_NOT_AUTHORIZED,
    HTTP_OK,
} from "../util/constants";
import LoginDto from "../models/loginDto";
import LoginService from "../services/login";

/**
 * @swagger
 * components:
 *   schemas:
 *     LoginDto:
 *       type: object
 *       properties:
 *         username:
 *           type: string
 *           example: "user123"
 *         password:
 *           type: string
 *           example: "password123"
 *       required:
 *         - username
 *         - password
 */
export const router = express.Router();
router.use(express.json());

/**
 * @swagger
 * tags:
 *   name: Auth
 * /auth/login:
 *   post:
 *     tags: [Auth]
 *     summary: Login
 *     description: Authenticates a user and returns an access token.
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/LoginDto'
 *     responses:
 *       '200':
 *         description: Access token successfully generated
 *         content:
 *           application/json:
 *             schema:
 *               type: string
 *       '401':
 *         description: Unauthorized
 *       '500':
 *         description: Internal server error
 */
router.post("/login", async (req: Request, res: Response) => {
    try {
        const service = new LoginService();
        const dto = req.body as LoginDto;
        const accessToken = await service.Login(dto);
        if (!accessToken) {
            res.status(HTTP_NOT_AUTHORIZED).send(
                "Account not found. Please sign-up to create a new account."
            );
            return;
        }
        res.status(HTTP_OK).send(accessToken);
    } catch (error) {
        console.error(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(error);
    }
});

/**
 * @swagger
 * tags:
 *   name: Auth
 * /auth/sign-in:
 *   post:
 *     tags: [Auth]
 *     summary: Sign-Up
 *     description: Registers a new user and returns an access token.
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/LoginDto'
 *     responses:
 *       '200':
 *         description: Access token successfully generated
 *         content:
 *           application/json:
 *             schema:
 *               type: string
 *       '500':
 *         description: Internal server error
 */
router.post("/sign-in", async (req: Request, res: Response) => {
    try {
        const service = new LoginService();
        const dto = req.body as LoginDto;
        const accessToken = await service.SignUp(dto);

        if (!accessToken) {
            res.status(HTTP_INTERNAL_SERVER_ERROR).send(
                "We are unable to process your request at this moment."
            );
        }
        res.status(HTTP_OK).send(accessToken);
    } catch (error) {
        console.error(error);
        res.status(HTTP_INTERNAL_SERVER_ERROR).send(error);
    }
});
