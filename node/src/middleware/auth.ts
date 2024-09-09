import { Request, Response, NextFunction } from "express";
import jwt from "jsonwebtoken";
import { jwtSecret } from "../services/login";
import { HTTP_NOT_AUTHORIZED } from "../util/constants";

const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
    try {
        const nonProtectedRoutes = ["/auth/login", "auth/signup"];

        if (process.env.MODE === "dev") {
            console.log(
                "Running in development mode, authentication is disabled."
            );
            return next();
        }

        if (nonProtectedRoutes.includes(req.path)) {
            return next();
        }

        const token = req.headers?.authorization?.split(" ")[1];

        if (!token) {
            res.status(HTTP_NOT_AUTHORIZED).send("Token not provided");
            return;
        }

        jwt.verify(token, jwtSecret);
        next();
    } catch (error) {
        console.error(error);
        res.status(HTTP_NOT_AUTHORIZED).send("Invalid access token");
    }
};

export default authMiddleware;
