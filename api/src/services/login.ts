import { UUID } from "mongodb";
import jwt from "jsonwebtoken";
import LoginDto from "../models/loginDto";
import User from "../models/user";
import { collections } from "./database";
import { createHash } from "crypto";

export const jwtSecret = process.env.JWT_SECRET_KEY || "secret";

export default class LoginService {
    constructor() {}

    public async Login(dto: LoginDto): Promise<string | null> {
        const query = { username: dto.username };
        const user = (await collections.users?.findOne(
            query
        )) as unknown as User;

        if (!user || user.password !== dto.password) {
            return null;
        }
        return this.createToken(user);
    }

    public async SignUp(dto: LoginDto): Promise<string | null> {
        const user: User = {
            id: UUID.toString(),
            username: dto.username,
            password: dto.password,
        };
        const res = await collections.users?.insertOne(user);
        const insertedId = res?.insertedId.toString();

        if (!insertedId) {
            return null;
        }
        return this.createToken(user);
    }

    private createToken(user: User): string {
        const jwtPayload = {
            userId: user.id,
            username: user.username,
        };
        return jwt.sign(jwtPayload, jwtSecret, {
            expiresIn: "1h",
        });
    }

    private hashPassword(password: string): string {
        const hash = createHash("sha256");
        const res = hash.update(password);
        return res.digest("hex");
    }
}
