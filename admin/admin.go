package admin

import("github.com/gin-gonic/gin"
       "daksha-leaderboard/db"
       "daksha-leaderboard/models"
       "go.mongodb.org/mongo-driver/mongo/options"
       "go.mongodb.org/mongo-driver/bson"
       "context"
       "net/http"
       "time"
       "os"
       "fmt"
       "encoding/json"
       "strconv"
       )

func Login(c *gin.Context){
	session,err := c.Request.Cookie("logged_in")
	if err!=nil || session.Value == ""{
		c.HTML(401,"login.tmpl",gin.H{})
		return
	}
	c.Redirect(303,"/update")
}

func AdminLogin(c *gin.Context){
	username := c.PostForm("username")
	pass := c.PostForm("password")

	if username == os.Getenv("USERNAME") && pass == os.Getenv("PASS"){
		http.SetCookie(c.Writer,&http.Cookie{
					Name:"logged_in",
					Value:"true",
					Path:"/",
					Expires:time.Now().Add(time.Minute * time.Duration(30)),
		})
		c.Redirect(303,"/update")
		return
	}else{
		c.HTML(401,"login.tmpl",gin.H{})
	}
}

func Add(c *gin.Context){
	session,err := c.Request.Cookie("logged_in")
	if err!=nil || session.Value == ""{
		c.Redirect(303,"/admin")
		return
	}
	c.HTML(200,"add-clg.tmpl",gin.H{})

}

func AddCollege(c *gin.Context){
	var college models.College

	college.Name = c.PostForm("name")
	college.Id = db.GenerateId()
	insert,err := db.Database.Collection("colleges").InsertOne(context.TODO(),college)
	if err!=nil{
		c.JSON(400,gin.H{"msg":"failed to add college"})
		return
	}

	c.HTML(200,"add-clg.tmpl",gin.H{"msg":"college added successfully","insertionID":insert.InsertedID})
}

func UpdatePage(c * gin.Context){
	var colleges []*models.College

	cur,err := db.Database.Collection("colleges").Find(context.TODO(),bson.D{{}})
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"unknown error. failed to fetch data"})
	}

	for cur.Next(context.TODO()){
                var c models.College
                err := cur.Decode(&c)
                if err!=nil{
                        fmt.Println(err)
                }

                colleges = append(colleges,&c)
        }
	c.HTML(200,"admin.tmpl",gin.H{"colleges":colleges})
}

func Update(c *gin.Context){

	var points models.Points

	a := c.PostForm("college")
        college := models.College{}
        json.Unmarshal([]byte(a),&college)

	points.College = college
	points.Score,_ = strconv.ParseFloat(c.PostForm("score"),64)

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
	fmt.Println(updateScore)

	c.Redirect(303,"/update")
}
