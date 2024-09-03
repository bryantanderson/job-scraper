import http from "http";
import dotenv from "dotenv";
import { wrapResponse, wrapText } from "./util/util";
import { HTTP_BAD_REQUEST, HTTP_NOT_FOUND } from "./constants";
import jobRouter from "./routes/job";

dotenv.config();

const port = Number(process.env.PORT) || 3000;
const scheme = process.env.SCHEME || "http";
const hostname = process.env.HOSTNAME || "0.0.0.0";
const server = http.createServer(
    (req: http.IncomingMessage, res: http.ServerResponse) => {
        // Router
        if (req.url === undefined || req.method === undefined) {
            wrapResponse(
                res,
                HTTP_BAD_REQUEST,
                wrapText("Undefined URL or method")
            );
            return;
        }
        const urlParts: string[] = req.url.split("/");
        const baseRoute = urlParts.length > 1 ? urlParts[1] : urlParts[0];

        switch (baseRoute) {
            case "jobs":
                jobRouter.get(req.method)?.apply(this, [req, res]);
                break;
            default:
                wrapResponse(
                    res,
                    HTTP_NOT_FOUND,
                    wrapText("Route does not exist")
                );
                return;
        }
    }
);

server.listen(port, hostname, () => {
    console.log(`Node listening on address ${scheme}://${hostname}:${port}`);
});
