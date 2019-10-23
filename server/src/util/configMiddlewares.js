const cors = require('@koa/cors');
const helmet = require('koa-helmet');
const bodyparser = require('koa-bodyparser');
const morgan = require('koa-morgan');

const {NODE_ENV} = process.env;
const MORGAN_LOG_TYPE = NODE_ENV === 'development' ? 'dev' : 'combined';

const responseTime = () => async (ctx, next) => {
    const start = Date.now();
    await next();
    const duration = Math.ceil(Date.now() - start);
    ctx.set('X-Response-Time', `${duration}ms`)
};

module.exports = (app) => {
    app.use(responseTime());
    // HTTP CORS configuration
    app.use(cors());
    // HTTP/Headers protection
    app.use(helmet());
    // Request body parser
    app.use(bodyparser());
    // Request logger
    app.use(morgan(MORGAN_LOG_TYPE))
};