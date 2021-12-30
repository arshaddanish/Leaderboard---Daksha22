package admin

import("github.com/gin-gonic/gin"
       "daksha-leaderboard/db"
       "daksha-leaderboard/models"
       "go.mongodb.org/mongo-driver/mongo/options"
       "go.mongodb.org/mongo-driver/bson"
       "context"
       )

func AddCollege(c *gin.Context){
	var college models.College
	c.BindJSON(&college)

	college.Id = db.GenerateId()
	insert,err := db.Database.Collection("colleges").InsertOne(context.TODO(),college)
	if err!=nil{
		c.JSON(400,gin.H{"msg":"failed to add college"})
		return
	}

	c.JSON(200,gin.H{"msg":"added successfully","insertionID":insert.InsertedID})
}

func Update(c *gin.Context){
	var points models.Points
	c.BindJSON(&points)

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"college.id",points.College.Id},{"college.name",points.College.Name}}
	update := bson.D{
		{"$inc",bson.D{
			{"score",points.Score},
		}},
	}

	updateScore,err := db.Database.Collection("scores").UpdateOne(context.TODO(),filter,update,opts)
	if err!=nil{
		c.JSON(400,gin.H{"msg":"failed to update"})
		return
	}

	c.JSON(200,gin.H{"msg":"updated","results":updateScore})
}
