const mongoose = require('mongoose');
const logger = require('./util/logger');

const {DATABASE_URI} = process.env;


mongoose.connection
    .once('open', _ => logger.info('Connected to database with success.'))
    .on('error', _ => logger.error('Error connecting to database', _));

module.exports = () => mongoose.connect(DATABASE_URI, {useMongoClient: true})