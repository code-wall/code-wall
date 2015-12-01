package mongo

import (
	"github.com/alechewitt/code-wall/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Default(connectionString string) *MongoDatabase {
	m := new(MongoDatabase)
	m.connectString = connectionString
	return m
}

type MongoDatabase struct {
	session       *mgo.Session
	connectString string
}

func (md *MongoDatabase) connect() *mgo.Session {
	if md.session == nil {
		var err error
		md.session, err = mgo.Dial(md.connectString)
		if err != nil {
			panic(err)
		}
		md.session.SetMode(mgo.Monotonic, true)
	}
	return md.session
}

func (md *MongoDatabase) FindById(id string) (sc database.Snippet, err error) {
	s := md.connect().Copy()
	defer s.Close()
	c := s.DB("").C("snippets")
	var result SnippetResult = SnippetResult{}
	err = c.FindId(id).One(&result)
	sc = &result
	return
}

func (md *MongoDatabase) Insert(data database.Snippet) (id string, err error) {
	s := md.connect().Copy()
	defer s.Close()
	c := s.DB("").C("snippets")
	var result SnippetResult = newResult(data)
	err = c.Insert(result)
	if err == nil {
		id = result.Id
	}
	return
}

func newResult(data database.Snippet) SnippetResult {
	return SnippetResult{
		Id:       bson.NewObjectId().Hex(),
		Snippet:  data.GetSnippet(),
		Language: data.GetLanguage(),
		Created:  data.GetDateCreated(),
	}
}
