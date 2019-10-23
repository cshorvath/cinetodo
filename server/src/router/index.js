const Router = require("koa-router");

const controller = require("../controller/list");

const router = new Router({prefix: "/v1/list"});

router.get("/:listId", controller.get);
router.post("/", controller.create);
router.post("/:listId/movie", controller.addMovie);
router.patch("/:listId/movie/:movieId", controller.updateMovie);
router.del("/:listId/movie/:movieId", controller.updateMovie);

module.exports = router;