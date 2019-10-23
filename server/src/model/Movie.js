const {Schema, model} = require("mongoose");

module.exports = new Schema({
    _id: Number,
    title: String,
    originalTitle: String,
    year: Number,
    director: String
});
