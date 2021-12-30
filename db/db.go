package db

import ("context"
         "go.mongodb.org/mongo-driver/mongo"
         "go.mongodb.org/mongo-driver/mongo/options"
         "fmt"
         "log"
         "os"
	 "github.com/rs/xid"
       )


 var clientOptions = options.Client().ApplyURI("mongodb+srv://"+os.Getenv("DBUSER")+":"+os.Getenv("DBPASS")+"@cluster0.ma5ba.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
 var Client, Err = mongo.Connect(context.TODO(), clientOptions)
 var Database = Client.Database(os.Getenv("DBNAME"))

func Connect()  {

 if Err != nil {
     log.Fatal(Err)
 }
 Err = Client.Ping(context.TODO(), nil)

 if Err != nil {
     log.Fatal(Err)
 }

 fmt.Println("Connected to MongoDB!")
}

func GenerateId() string{

	guid := xid.New()
	id := guid.String()

	return (id)
}
