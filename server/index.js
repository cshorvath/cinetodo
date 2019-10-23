const Koa = require('koa');
const logger = require("./src/util/logger");
const configMiddlewares = require("./src/util/configMiddlewares");
const router = require("./src/router");

const port = process.env.PORT || 3000;

logger.info(`Server listening on port ${port}`);
const app = new Koa();
configMiddlewares(app);
app.use(router.routes());
app.listen(port);