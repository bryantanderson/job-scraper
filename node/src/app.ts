import express from "express";
import dotenv from "dotenv";
import { router as jobRouter } from "./routes/job";
import authMiddleware from "./middleware/auth";
import { connectToMongoDB } from "./repository/mongo";
import swaggerJsDoc from "swagger-jsdoc";
import swaggerUi from "swagger-ui-express";

dotenv.config();

const port: number = Number(process.env.PORT) || 5000;
const server: express.Express = express();

const swaggerOptions = {
    definition: {
        openapi: "3.1.0",
        info: {
            title: "NodeJS + ExpressJS",
            version: "0.1.0",
            description: "In progress. To be new main backend",
            contact: {
                name: "bciputra",
                email: "@",
            },
        },
        servers: [
            {
                url: `http://localhost:${port}`,
            },
        ],
    },
    apis: ["./src/routes/*.ts"],
};

const swaggerSpecs = swaggerJsDoc(swaggerOptions);

connectToMongoDB()
    .then(() => {
        // Register middleware
        server.use(authMiddleware);

        // Register routes
        server.use(
            "/docs",
            swaggerUi.serve,
            swaggerUi.setup(swaggerSpecs, { explorer: true })
        );
        server.use("/jobs", jobRouter);

        server.listen(port, () => {
            console.log(`Node listening on port ${port}`);
        });
    })
    .catch((error: Error) => {
        console.error("Database connection failed", error);
        process.exit();
    });
