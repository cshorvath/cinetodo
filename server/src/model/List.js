const {Schema} = require("mongoose");
const db = require("../db");
const shortid = require("shortid");

const Movie = require("./Movie");

const schema = new Schema({
    _id: {type: String, default: shortid.generate},
    title: String,
    date: {type: Date, default: Date.now},
    movies: [Movie]
});

module.exports = db.model('List', schema);