const List = require("../model/List");

const create = async (ctx) => {
    const {title} = ctx.request.body;
    const list = new List({title});
    return await list.save();
};

const get = async (ctx) => {
    return List.findById(ctx.params.listId);
};

const addMovie = async (ctx) => {
    const {listId} = ctx.params;
    const {title, originalTitle, year, director, id} = ctx.request.body;
    List.findByIdAndUpdate(listId, {$addToSet: {movies: {title, originalTitle, year, director, id}}})
};

const updateMovie = async (ctx) => {
    const {listId, movieId} = ctx.params.listId;
    const {seen} = ctx.request.body;
    await List.update(
        {'_id': listId, 'movies._id': movieId},
        {
            $set: {
                'movies.$.seen': seen
            }
        }
    );
};

const deleteMovie = async (ctx) => {
    const {listId, movieId} = ctx.params.listId;
    await List.findByIdAndUpdate(listId,
        {
            '$pull': {
                'movies': {'_id': movieId}
            }
        }
    );
};


module.exports = {create, get, addMovie, updateMovie, deleteMovie};