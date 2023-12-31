package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"companies/security"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const colNames = "companies"

const (
	corpo = "Corporations"
	nProf = "NonProfit"
	coop  = "Cooperative"
	solPr = "Sole Proprietorship"
)

type Company struct {
	ID                string `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	AmountOfEmployers int    `json:"amountOfEmployers"`
	Registered        bool   `json:"registered"`
	Type              string `json:"type"`
}

type Handler struct {
	dbHandler    *MongoCRUD
	secValidator *security.Security
}

func NewHandler(url, user, pass string, sec *security.Security) *Handler {
	return &Handler{dbHandler: NewMongoCrud(url, user, pass), secValidator: sec}
}

func (han *Handler) Init() {
	han.dbHandler.Connect()
}

func (han *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("token")
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, ok := han.secValidator.VerifyToken(tokenStr)
	if !ok {
		http.Error(w, fmt.Sprint(msg), http.StatusUnauthorized)
		return
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	company := Company{}
	errUnm := json.Unmarshal(reqBody, &company)
	if errUnm != nil {
		http.Error(w, fmt.Sprintf("%v", errUnm), http.StatusBadRequest)
		return
	}
	valTypeComp := validateType(company.Type)
	if !valTypeComp {
		http.Error(w, fmt.Sprintf("type of company %s not valid", company.Type), http.StatusBadRequest)
		return
	}
	result, errCreate := han.dbHandler.Create(company)
	if errCreate != nil {
		http.Error(w, fmt.Sprintf("%v", errCreate), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(result)
	if errEncode != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", errEncode), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (han *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Token")
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, ok := han.secValidator.VerifyToken(tokenStr)
	if !ok {
		http.Error(w, fmt.Sprint(msg), http.StatusUnauthorized)
		return
	}
	name := chi.URLParam(r, "name")
	result, errDel := han.dbHandler.Delete(name)
	if errDel != nil {
		http.Error(w, fmt.Sprintf("%v", errDel), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "number of deleted companies:%d", result)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func (han *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Token")
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, ok := han.secValidator.VerifyToken(tokenStr)
	if !ok {
		http.Error(w, fmt.Sprint(msg), http.StatusUnauthorized)
		return
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	company := Company{}
	errUnm := json.Unmarshal(reqBody, &company)
	if errUnm != nil {
		http.Error(w, fmt.Sprintf("%v", errUnm), http.StatusBadRequest)
		return
	}
	valTypeComp := validateType(company.Type)
	if !valTypeComp {
		http.Error(w, fmt.Sprintf("type of company %s not valid", company.Type), http.StatusBadRequest)
		return
	}
	result, errCreate := han.dbHandler.Update(company)
	if errCreate != nil {
		http.Error(w, fmt.Sprintf("%v", errCreate), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(result)
	if errEncode != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", errEncode), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, errMsg := fmt.Fprintf(w, "number of updated companies:%d", result)
	if errMsg != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}
}

func (han *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Token")
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	msg, ok := han.secValidator.VerifyToken(tokenStr)
	if !ok {
		http.Error(w, fmt.Sprint(msg), http.StatusUnauthorized)
		return
	}
	name := chi.URLParam(r, "name")
	result, errRead := han.dbHandler.Read(name)
	if errRead != nil {
		http.Error(w, fmt.Sprintf("%v", errRead), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(result)
	if errEncode != nil {
		http.Error(w, fmt.Sprintf("error building the response, %v", errEncode), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type MongoCRUD struct {
	url        string
	user       string
	passwd     string
	compDB     *mongo.Database
	collection *mongo.Collection
}

func NewMongoCrud(url, user, passwd string) *MongoCRUD {
	return &MongoCRUD{
		url:    url,
		user:   user,
		passwd: passwd,
	}
}

func (moDB *MongoCRUD) Connect() {
	cred := options.Credential{
		Username: moDB.user,
		Password: moDB.passwd,
	}
	cliOpts := options.Client().ApplyURI(moDB.url).SetAuth(cred)
	client, err := mongo.Connect(context.TODO(), cliOpts)
	if err != nil {
		panic(err)
	}
	moDB.compDB = client.Database("root")
	col := moDB.compDB.Collection(colNames)
	if col == nil {
		errCol := moDB.compDB.CreateCollection(context.TODO(), colNames)
		if errCol != nil {
			panic(errCol)
		}
		col = moDB.compDB.Collection(colNames)
	}
	moDB.collection = col
	log.Info().Msg("connected to mongoDB")
}

func (moDB *MongoCRUD) Create(comp Company) (*Company, error) {
	var doc Company
	moDB.collection.FindOne(context.TODO(), bson.M{"name": comp.Name}).Decode(&doc)
	if doc.ID != "" {
		log.Error().Msgf(fmt.Sprintf("company %s exist", comp.Name))
		return nil, fmt.Errorf("company %s exist", comp.Name)
	}
	res, err := moDB.collection.InsertOne(context.TODO(), comp)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	oidStr := oid.String()

	oidStr = strings.ReplaceAll(oidStr, "ObjectID(\"", "")
	oidStr = strings.ReplaceAll(oidStr, "\")", "")
	if ok {
		comp.ID = oidStr
	}
	return &comp, nil
}

func (moDB *MongoCRUD) Read(name string) (*Company, error) {
	res := moDB.collection.FindOne(context.TODO(), bson.M{"name": name})
	out := Company{}
	err := res.Decode(&out)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	return &out, nil
}

func (moDB *MongoCRUD) Update(comp Company) (int64, error) {
	comp.ID = ""
	res, errUpd := moDB.collection.ReplaceOne(context.TODO(), bson.M{"name": comp.Name}, comp)
	if errUpd != nil {
		log.Error().Msg(errUpd.Error())
		return 0, errUpd
	}
	return res.ModifiedCount, nil
}

func (moDB *MongoCRUD) Delete(name string) (int64, error) {
	res, err := moDB.collection.DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		log.Error().Msg(err.Error())
		return 0, err
	}
	return res.DeletedCount, nil
}

func validateType(typeCom string) bool {
	switch typeCom {
	case coop:
		return true
	case nProf:
		return true
	case corpo:
		return true
	case solPr:
		return true
	default:
		return false
	}
}
