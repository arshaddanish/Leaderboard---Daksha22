package main
import("github.com/gin-gonic/gin"
       _"net/http"
       "daksha-leaderboard/db"
       "daksha-leaderboard/models"
       "daksha-leaderboard/admin"
       "go.mongodb.org/mongo-driver/mongo/options"
       "go.mongodb.org/mongo-driver/bson"
       "context"
       "fmt"
       )


func LeaderBoard(c *gin.Context){
	collection := db.Database.Collection("scores")

//	Cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{
//		bson.D{{
//			"setWindowFields",bson.M{
//				"$sortBy":bson.M{"score":-1},
//				"$output":bson.M{
//					"rank":bson.M{"$denseRank":bson.D{},},
//				}},
//			}},
//	})
//	if err != nil {
//	    panic(err)
//	}
//	var response []bson.M
//	if err = Cursor.All(context.TODO(), &response); err != nil {
//	    panic(err)
//	}
//	fmt.Println(response)

	var data []*models.Points

	findOptions := options.Find().SetSort(bson.D{{"score",-1}})
	cur,err := collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		fmt.Println(err)
	}

	for cur.Next(context.TODO()){
		var p models.Points
		err := cur.Decode(&p)
		if err!=nil{
			fmt.Println(err)
		}

		data = append(data,&p)
	}

	data[0].Rank = 1;
	for i:=1;i<len(data);i++{
		if data[i].Score == data[i-1].Score{
			data[i].Rank = data[i-1].Rank;
		}else{
			data[i].Rank = data[i-1].Rank + 1;
		}
	}
	c.HTML(200,"leaderboard.tmpl",gin.H{"data":data})
}



func main(){
	r := gin.Default()

	r.Static("/assets","./assets")
	r.LoadHTMLGlob("*.tmpl")

	db.Connect()

	r.GET("/",LeaderBoard)
	r.POST("/add",admin.AddCollege)
	r.POST("/update",admin.Update)
	r.Run()
}
