const winston = require('winston');

const logger = new winston.createLogger({
    transports: [
        new winston.transports.Console()
    ],
});

module.exports = logger;